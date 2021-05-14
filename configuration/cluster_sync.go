// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/GehirnInc/crypt"
	"github.com/google/renameio"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/config-parser/v3/types"
	log "github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

const DataplaneAPIType = "community"

// Node is structure required for connection to cluster
type Node struct {
	Address     string            `json:"address"`
	APIBasePath string            `json:"api_base_path"`
	APIPassword string            `json:"api_password"`
	APIUser     string            `json:"api_user"`
	Certificate string            `json:"certificate,omitempty"`
	Description string            `json:"description,omitempty"`
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Port        int64             `json:"port,omitempty"`
	Status      string            `json:"status"`
	Type        string            `json:"type"`
	Variables   map[string]string `json:"variables"`
}

// ClusterSync fetches certificates for joining cluster
type ClusterSync struct {
	cfg         *Configuration
	certFetch   chan struct{}
	cli         *client_native.HAProxyClient
	Context     context.Context
	ReloadAgent haproxy.IReloadAgent
}

func (c *ClusterSync) Monitor(cfg *Configuration, cli *client_native.HAProxyClient) {
	c.cfg = cfg
	c.cli = cli

	go c.monitorBootstrapKey()
	if c.cfg.Mode.Load() == "cluster" {
		go c.monitorCertificateRefresh()
	}

	c.certFetch = make(chan struct{}, 2)
	go c.fetchCert()

	key := c.cfg.Cluster.BootstrapKey.Load()
	certFetched := cfg.Cluster.CertificateFetched.Load()

	if key != "" && !certFetched {
		runtime.Gosched()
		c.cfg.Notify.BootstrapKeyChanged.Notify()
	}
}

func (c *ClusterSync) monitorCertificateRefresh() {
	for range c.cfg.Notify.CertificateRefresh.Subscribe("monitorCertificateRefresh") {
		log.Info("refreshing certificate")

		key := c.cfg.Cluster.BootstrapKey.Load()
		data, err := DecodeBootstrapKey(key)
		if err != nil {
			log.Warning(err)
			continue
		}
		url := fmt.Sprintf("%s://%s", data["schema"], data["address"])

		csr, key, err := generateCSR()
		if err != nil {
			log.Warning(err)
			continue
		}
		err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s-csr.crt", c.cfg.Name.Load())), []byte(csr), 0644)
		if err != nil {
			log.Warning(err)
			continue
		}
		err = c.issueRefreshRequest(url, data["port"], data["api-base-path"], data["path"], csr, key)
		if err != nil {
			log.Warning(err)
			continue
		}
	}
}

func (c *ClusterSync) issueRefreshRequest(url, port, basePath string, nodesPath string, csr, key string) error {
	url = fmt.Sprintf("%s:%s/%s", url, port, strings.TrimLeft(path.Join(basePath, nodesPath, c.cfg.Cluster.ID.Load()), "/"))
	apiAddress := c.cfg.APIOptions.APIAddress
	if apiAddress == "" {
		apiAddress = c.cfg.RuntimeData.Host
	}
	nodeData := Node{
		ID:          c.cfg.Cluster.ID.Load(),
		Address:     apiAddress,
		Certificate: csr,
		Status:      cfg.Status.Load(),
		Type:        DataplaneAPIType,
	}
	bytesRepresentation, _ := json.Marshal(nodeData)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf("error creating new POST request for cluster comunication")
	}
	req.Header.Add("X-Node-Key", c.cfg.Cluster.Token.Load())
	req.Header.Add("Content-Type", "application/json")
	log.Infof("Refreshing certificate %s", url)
	httpClient := createHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body))
	}
	var responseData Node
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	log.Infof("Cluster re joined, status: %s", responseData.Status)
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", c.cfg.Name.Load())), []byte(csr), 0644)
	if err != nil {
		log.Warning(err)
		return err
	}
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", c.cfg.Name.Load())), []byte(key), 0644)
	if err != nil {
		log.Warning(err)
		return err
	}
	c.cfg.Cluster.Token.Store(resp.Header.Get("X-Node-Key"))
	err = c.cfg.Save()
	if err != nil {
		log.Warning(err)
		return err
	}
	c.cfg.Notify.Reload.Notify()
	return nil
}

func (c *ClusterSync) monitorBootstrapKey() {
	for {
		select {
		case <-c.cfg.Notify.BootstrapKeyChanged.Subscribe("monitorBootstrapKey"):
			key := c.cfg.Cluster.BootstrapKey.Load()
			c.cfg.Cluster.CertificateFetched.Store(false)
			if key == "" {
				// do we need to delete cert here maybe?
				c.cfg.Cluster.ActiveBootstrapKey.Store("")
				err := c.cfg.Save()
				if err != nil {
					log.Panic(err)
				}
				break
			}
			if key == c.cfg.Cluster.ActiveBootstrapKey.Load() {
				if !c.cfg.Cluster.CertificateFetched.Load() {
					c.certFetch <- struct{}{}
				}
				break
			}
			data, err := DecodeBootstrapKey(key)
			if err != nil {
				log.Warning(err)
			}
			url := fmt.Sprintf("%s://%s", data["schema"], data["address"])
			c.cfg.Cluster.URL.Store(url)
			c.cfg.Cluster.Port.Store(func() int {
				i, _ := strconv.Atoi(data["port"])
				return i
			}())
			c.cfg.Cluster.APIBasePath.Store(data["api-base-path"])
			registerPath, ok := data["register-path"]
			if !ok {
				c.cfg.Cluster.APIRegisterPath.Store(data["path"])
				c.cfg.Cluster.APINodesPath.Store(data["path"])
			} else {
				c.cfg.Cluster.APIRegisterPath.Store(registerPath)
				c.cfg.Cluster.APINodesPath.Store(data["nodes-path"])
			}
			c.cfg.Cluster.Name.Store(data["name"])
			c.cfg.Cluster.Description.Store(data["description"])
			c.cfg.Mode.Store("cluster")
			err = c.cfg.Save()
			if err != nil {
				log.Panic(err)
			}
			csr, key, err := generateCSR()
			if err != nil {
				log.Warning(err)
				break
			}
			err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", c.cfg.Name.Load())), []byte(key), 0644)
			if err != nil {
				log.Warning(err)
				break
			}
			err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s-csr.crt", c.cfg.Name.Load())), []byte(csr), 0644)
			if err != nil {
				log.Warning(err)
				break
			}
			err = c.cfg.Save()
			if err != nil {
				log.Panic(err)
			}
			err = c.issueJoinRequest(url, data["port"], data["api-base-path"], c.cfg.Cluster.APIRegisterPath.Load(), csr, key)
			if err != nil {
				log.Warning(err)
				break
			}
			if !c.cfg.Cluster.CertificateFetched.Load() {
				c.certFetch <- struct{}{}
			}
		case <-c.Context.Done():
			return
		}
	}
}

func (c *ClusterSync) issueJoinRequest(url, port, basePath string, registerPath string, csr, key string) error {
	url = fmt.Sprintf("%s:%s/%s", url, port, strings.TrimLeft(path.Join(basePath, registerPath), "/"))
	apiCfg := c.cfg.APIOptions
	userStore := GetUsersStore()

	// create a new user for connecting to cluster
	name, err := misc.RandomString(8)
	if err != nil {
		return err
	}
	pwd, err := misc.RandomString(24)
	if err != nil {
		return err
	}

	cryptAlg := crypt.New(crypt.SHA512)
	hash, err := cryptAlg.Generate([]byte(pwd), nil)
	if err != nil {
		return err
	}
	name = fmt.Sprintf("dpapi-c-%s", name)
	log.Infof("Creating user %s for cluster connection", name)
	user := types.User{
		Name:       name,
		IsInsecure: false,
		Password:   hash,
	}
	err = userStore.AddUser(user)
	if err != nil {
		return err
	}

	apiAddress := apiCfg.APIAddress
	if apiAddress == "" {
		apiAddress = c.cfg.RuntimeData.Host
	}
	apiPort := apiCfg.APIPort
	if apiPort == 0 {
		apiPort = int64(c.cfg.RuntimeData.Port)
	}

	nodeData := Node{
		// ID:          "",
		Address:     apiAddress,
		APIBasePath: c.cfg.RuntimeData.APIBasePath,
		APIPassword: pwd,
		APIUser:     user.Name,
		Certificate: csr,
		Description: "",
		Name:        c.cfg.Name.Load(),
		Port:        apiPort,
		Status:      "waiting_approval",
		Type:        DataplaneAPIType,
	}
	nodeData.Variables = map[string]string{}

	// report the dataplane_cmdline if started from within haproxy
	if c.cfg.HAProxy.MasterWorkerMode || os.Getenv("HAPROXY_MWORKER") == "1" {
		nodeData.Variables["dataplane_cmdline"] = c.cfg.Cmdline.String()
	}

	processInfos, err := c.cli.Runtime.GetInfo()
	if err != nil || len(processInfos) < 1 {
		log.Error("unable to fetch processInfo")
	} else {
		if processInfos[0].Info != nil {
			nodeData.Variables["haproxy_version"] = processInfos[0].Info.Version
		} else {
			log.Error("empty process info")
		}
	}

	bytesRepresentation, _ := json.Marshal(nodeData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf("error creating new POST request for cluster comunication")
	}
	req.Header.Add("X-Bootstrap-Key", c.cfg.Cluster.BootstrapKey.Load())
	req.Header.Add("Content-Type", "application/json")
	log.Infof("Joining cluster %s", url)
	httpClient := createHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body))
	}
	var responseData Node
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	if c.cfg.HAProxy.NodeIDFile != "" {
		// write id to file
		errFID := ioutil.WriteFile(c.cfg.HAProxy.NodeIDFile, []byte(responseData.ID), 0644) // nolint:gosec
		if errFID != nil {
			return errFID
		}
		version, errVersion := c.cli.Configuration.GetVersion("")
		if errVersion != nil || version < 1 {
			// silently fallback to 1
			version = 1
		}
		t, err1 := c.cli.Configuration.StartTransaction(version)
		if err1 != nil {
			return err1
		}
		// write id to peers
		_, peerSections, errorGet := c.cli.Configuration.GetPeerSections(t.ID)
		if errorGet != nil {
			return errorGet
		}
		peerFound := false
		dataplaneID := c.cfg.Cluster.ID.Load()
		if dataplaneID == "" {
			dataplaneID = "localhost"
		}
		for _, section := range peerSections {
			_, peerEntries, err1 := c.cli.Configuration.GetPeerEntries(section.Name, t.ID)
			if err1 != nil {
				return err1
			}
			for _, peer := range peerEntries {
				if peer.Name == dataplaneID {
					peerFound = true
					peer.Name = responseData.ID
					errEdit := c.cli.Configuration.EditPeerEntry(dataplaneID, section.Name, peer, t.ID, 0)
					if errEdit != nil {
						_ = c.cli.Configuration.DeleteTransaction(t.ID)
						return err
					}
				}
			}
		}
		if !peerFound {
			_ = c.cli.Configuration.DeleteTransaction(t.ID)
			return fmt.Errorf("peer [%s] not found in HAProxy config", dataplaneID)
		}
		_, err = c.cli.Configuration.CommitTransaction(t.ID)
		if err != nil {
			return err
		}
		// restart HAProxy
		errRestart := c.ReloadAgent.Restart()
		if errRestart != nil {
			return errRestart
		}
	}
	c.cfg.Cluster.ID.Store(responseData.ID)
	c.cfg.Cluster.Name.Store(responseData.Name)
	c.cfg.Cluster.Token.Store(resp.Header.Get("X-Node-Key"))
	c.cfg.Cluster.ActiveBootstrapKey.Store(c.cfg.Cluster.BootstrapKey.Load())
	log.Info("Cluster joined")
	_, err = c.checkCertificate(responseData)
	if err != nil {
		return err
	}
	err = c.cfg.Save()
	if err != nil {
		return err
	}
	return nil
}

// checkCertificate checks if we have received valid certificate or we just got CSR back
//
// two options are possible here:
// -----BEGIN CERTIFICATE----- or -----BEGIN CERTIFICATE REQUEST-----
func (c *ClusterSync) checkCertificate(node Node) (fetched bool, err error) {
	if !strings.HasPrefix(node.Certificate, "-----BEGIN CERTIFICATE-----") {
		c.cfg.Status.Store("unconfigured")
		return false, nil
	}
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", c.cfg.Name.Load())), []byte(node.Certificate), 0644)
	if err != nil {
		c.cfg.Status.Store("unconfigured")
		return false, err
	}
	c.cfg.Cluster.CertificateFetched.Store(true)
	c.cfg.Notify.Reload.Notify()
	c.cfg.Status.Store("active")
	return true, nil
}

func (c *ClusterSync) activateFetchCert(err error) {
	go func(err error) {
		log.Warning(err)
		time.Sleep(1 * time.Minute)
		if !c.cfg.Cluster.CertificateFetched.Load() {
			c.certFetch <- struct{}{}
		}
	}(err)
}

func (c *ClusterSync) fetchCert() {
	for {
		select {
		case <-c.Context.Done():
			close(c.certFetch)
			return
		case <-c.certFetch:
			key := c.cfg.Cluster.BootstrapKey.Load()
			if key == "" || c.cfg.Cluster.Token.Load() == "" {
				break
			}
			// if not, sleep and start all over again
			certFetched := c.cfg.Cluster.CertificateFetched.Load()
			if !certFetched {
				url := c.cfg.Cluster.URL.Load()
				port := c.cfg.Cluster.Port.Load()
				apiBasePath := c.cfg.Cluster.APIBasePath.Load()
				apiNodesPath := c.cfg.Cluster.APINodesPath.Load()
				id := c.cfg.Cluster.ID.Load()
				url = fmt.Sprintf("%s:%d/%s/%s/%s", url, port, apiBasePath, apiNodesPath, id)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					c.activateFetchCert(err)
					break
				}
				req.Header.Add("X-Node-Key", c.cfg.Cluster.Token.Load())
				req.Header.Add("Content-Type", "application/json")
				httpClient := createHTTPClient()
				resp, err := httpClient.Do(req)
				if err != nil {
					c.activateFetchCert(err)
					break
				}
				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					c.activateFetchCert(err)
					break
				}
				if resp.StatusCode != 200 {
					c.activateFetchCert(fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body)))
					break
				}
				var responseData Node
				err = json.Unmarshal(body, &responseData)
				if err != nil {
					c.activateFetchCert(err)
					break
				}
				log.Warningf("Fetching certificate, status: %s", responseData.Status)

				certFetched, err = c.checkCertificate(responseData)
				if err != nil {
					log.Warning(err.Error())
					break
				}
				err = c.cfg.Save()
				if err != nil {
					log.Warning(err)
				}
			}
			if !certFetched {
				time.AfterFunc(time.Minute, func() {
					if !c.cfg.Cluster.CertificateFetched.Load() {
						c.certFetch <- struct{}{}
					}
				})
			}
		}
	}
}

func generateCSR() (string, string, error) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 2048)

	oidEmailAddress := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}
	emailAddress := "test@example.com"
	subj := pkix.Name{
		CommonName:         "haproxy.com",
		Country:            []string{"US"},
		Province:           []string{""},
		Locality:           []string{"Waltham"},
		Organization:       []string{"HAProxy Technologies LLC"},
		OrganizationalUnit: []string{"IT"},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: oidEmailAddress,
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte(emailAddress),
				},
			},
		},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	var buf bytes.Buffer
	err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	if err != nil {
		return "", "", err
	}

	caPrivKeyPEMBuff := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEMBuff, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyBytes),
	})
	if err != nil {
		return "", "", err
	}
	return buf.String(), caPrivKeyPEMBuff.String(), err
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			TLSClientConfig: &tls.Config{
				//nolint
				InsecureSkipVerify: true, // this is deliberate, might only have self signed certificate
			},
		},
		Timeout: time.Duration(10) * time.Second,
	}
	return client
}
