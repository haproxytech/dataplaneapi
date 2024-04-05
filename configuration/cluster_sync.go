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
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/renameio"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/config-parser/v5/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/haproxytech/dataplaneapi/log"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

const DataplaneAPIType = "community"

// Node is structure required for connection to cluster
type Node struct {
	Facts       map[string]string `json:"facts"`
	Address     string            `json:"address"`
	APIBasePath string            `json:"api_base_path"`
	APIPassword string            `json:"api_password"`
	APIUser     string            `json:"api_user"`
	Certificate string            `json:"certificate,omitempty"`
	Description string            `json:"description,omitempty"`
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Type        string            `json:"type"`
	Port        int64             `json:"port,omitempty"`
}

// ClusterSync fetches certificates for joining cluster
type ClusterSync struct {
	cfg         *Configuration
	certFetch   chan struct{}
	cli         client_native.HAProxyClient
	Context     context.Context
	ReloadAgent haproxy.IReloadAgent
}

var expectedResponseCodes = map[string]int{
	"POST": 201,
	"PUT":  200,
}

func (c *ClusterSync) Monitor(cfg *Configuration, cli client_native.HAProxyClient) {
	c.cfg = cfg
	c.cli = cli

	go c.monitorBootstrapKey()
	go c.monitorCertificateRefresh()

	c.certFetch = make(chan struct{}, 2)
	go c.fetchCert()

	<-c.cfg.Notify.ServerStarted.Subscribe("clusterMonitor")

	key := c.cfg.Cluster.BootstrapKey.Load()
	certFetched := cfg.Cluster.CertificateFetched.Load()

	if key != "" && !certFetched {
		c.cfg.Notify.BootstrapKeyChanged.Notify()
	}
}

func (c *ClusterSync) monitorCertificateRefresh() {
	for {
		select {
		case <-c.cfg.Notify.CertificateRefresh.Subscribe("monitorCertificateRefresh"):
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
			err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s-csr.crt", c.cfg.Name.Load())), []byte(csr), 0o644)
			if err != nil {
				log.Warning(err)
				continue
			}
			err = c.issueRefreshRequest(url, data["port"], data["api-base-path"], data["path"], csr, key)
			if err != nil {
				log.Warning(err)
				continue
			}
		case <-c.Context.Done():
			return
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
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	bytesRepresentation, _ := json.Marshal(nodeData)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return errors.New("error creating new POST request for cluster comunication")
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body))
	}
	var responseData Node
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	log.Infof("Cluster re joined, status: %s", responseData.Status)
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", c.cfg.Name.Load())), []byte(csr), 0o644)
	if err != nil {
		log.Warning(err)
		return err
	}
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", c.cfg.Name.Load())), []byte(key), 0o644)
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
			log.Warningf("detected change in configured bootstrap key")
			key := c.cfg.Cluster.BootstrapKey.Load()
			c.cfg.Cluster.CertificateFetched.Store(false)
			if key == "" {
				// do we need to delete cert here maybe?
				log.Warningf("setting bootstrap key to empty")
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
			// check if storage key is ok
			errStorageDir := CheckIfStorageDirIsOK(data["storage-dir"], c.cfg)
			if errStorageDir != err {
				log.Error(errStorageDir)
				continue
			}
			// Init NOTICE file to inform user that the cluster storage folder is programmatically managed by Fusion API
			if errStorageInit := InitStorageNoticeFile(data["storage-dir"]); errStorageInit != nil {
				log.Warningf("unable to create notice file, %s: skipping it", errStorageInit.Error())
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
			c.cfg.Cluster.StorageDir.Store(data["storage-dir"])
			c.cfg.Cluster.ClusterID.Store(data["cluster-id"])
			c.cfg.HAProxy.ClusterTLSCertDir = path.Join(data["storage-dir"], "certs-cluster")
			c.cfg.Cluster.CertificateDir.Store(path.Join(data["storage-dir"], "certs-cluster"))
			c.cfg.Mode.Store(ModeCluster)
			err = c.cfg.Save()
			if err != nil {
				log.Panic(err)
			}
			csr, key, err := generateCSR()
			if err != nil {
				log.Warning(err)
				break
			}
			err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", c.cfg.Name.Load())), []byte(key), 0o644)
			if err != nil {
				log.Warning(err)
				break
			}
			err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s-csr.crt", c.cfg.Name.Load())), []byte(csr), 0o644)
			if err != nil {
				log.Warning(err)
				break
			}
			err = c.cfg.Save()
			if err != nil {
				log.Panic(err)
			}
			registerMerhod := "POST"
			if method, ok := data["register-method"]; ok {
				registerMerhod = method
			}
			log.Warningf("issuing cluster join request to cluster %s at %s", data["name"], data["address"])
			userStore := GetUsersStore()
			user, pwd, err := misc.CreateClusterUser()
			if err != nil {
				log.Error(err)
				break
			}
			err = userStore.AddUser(user)
			if err != nil {
				log.Error(err)
				break
			}
			backOff := 1
			numTries := 0
			maxTries := 10
			for {
				err = c.issueJoinRequest(url, data["port"], data["api-base-path"], c.cfg.Cluster.APIRegisterPath.Load(), registerMerhod, csr, key, user, pwd)
				if err == nil {
					break
				}
				log.Error(err)
				if !misc.IsNetworkErr(err) {
					break
				}
				numTries++
				backOff *= 2
				if backOff > 60 {
					backOff = 60
				}
				if numTries > maxTries {
					log.Error("Joining cluster failed")
					break
				}
				log.Warningf("Joining cluster will be retried after %d seconds [%d/%d]", backOff, numTries, maxTries)
				time.Sleep(time.Second * time.Duration(backOff))
			}
			if err != nil {
				break
			}

			if !c.cfg.Cluster.CertificateFetched.Load() {
				log.Warningf("starting certificate fetch")
				c.certFetch <- struct{}{}
			}
		case <-c.Context.Done():
			return
		}
	}
}

func (c *ClusterSync) issueJoinRequest(url, port, basePath string, registerPath string, registerMethod string, csr, key string, user types.User, userPWD string) error {
	url = fmt.Sprintf("%s:%s/%s", url, port, strings.TrimLeft(path.Join(basePath, registerPath), "/"))
	apiCfg := c.cfg.APIOptions

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
		APIPassword: userPWD,
		APIUser:     user.Name,
		Certificate: csr,
		Description: "",
		Name:        c.cfg.Name.Load(),
		Port:        apiPort,
		Status:      "waiting_approval",
		Type:        DataplaneAPIType,
	}
	nodeData.Facts = c.getNodeFacts()

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	bytesRepresentation, _ := json.Marshal(nodeData)

	req, err := http.NewRequest(registerMethod, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf("error creating new %s request for cluster comunication", registerMethod)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != expectedResponseCodes[registerMethod] {
		return fmt.Errorf("invalid status code [%d] %s", resp.StatusCode, string(body))
	}
	log.Warningf("success sending local info, joining in progress")
	var responseData Node
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}
	if c.cfg.HAProxy.NodeIDFile != "" {
		configuration, errCfg := c.cli.Configuration()
		if errCfg != nil {
			return errCfg
		}
		// write id to file
		errFID := renameio.WriteFile(c.cfg.HAProxy.NodeIDFile, []byte(responseData.ID), 0o644)
		if errFID != nil {
			return errFID
		}
		version, errVersion := configuration.GetVersion("")
		if errVersion != nil || version < 1 {
			// silently fallback to 1
			version = 1
		}
		t, err1 := configuration.StartTransaction(version)
		if err1 != nil {
			return err1
		}
		// write id to peers
		_, peerSections, errorGet := configuration.GetPeerSections(t.ID)
		if errorGet != nil {
			return errorGet
		}
		peerFound := false
		dataplaneID := c.cfg.Cluster.ID.Load()
		if dataplaneID == "" {
			dataplaneID = "localhost"
		}
		for _, section := range peerSections {
			_, peerEntries, err1 := configuration.GetPeerEntries(section.Name, t.ID)
			if err1 != nil {
				return err1
			}
			for _, peer := range peerEntries {
				if peer.Name == dataplaneID {
					peerFound = true
					peer.Name = responseData.ID
					errEdit := configuration.EditPeerEntry(dataplaneID, section.Name, peer, t.ID, 0)
					if errEdit != nil {
						_ = configuration.DeleteTransaction(t.ID)
						return err
					}
				}
			}
		}
		if !peerFound {
			_ = configuration.DeleteTransaction(t.ID)
			return fmt.Errorf("peer [%s] not found in HAProxy config", dataplaneID)
		}
		_, err = configuration.CommitTransaction(t.ID)
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
	log.Warning("cluster joined")
	_, err = c.checkCertificate(responseData)
	if err != nil {
		return err
	}
	return c.cfg.Save()
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
	err = renameio.WriteFile(path.Join(c.cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", c.cfg.Name.Load())), []byte(node.Certificate), 0o644)
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
				url = fmt.Sprintf("%s:%d/%s", url, port, strings.TrimLeft(path.Join(apiBasePath, apiNodesPath, id), "/"))
				req, err := http.NewRequest(http.MethodGet, url, nil)
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
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					c.activateFetchCert(err)
					break
				}
				if resp.StatusCode != http.StatusOK {
					c.activateFetchCert(fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body)))
					break
				}
				var responseData Node
				json := jsoniter.ConfigCompatibleWithStandardLibrary
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
						log.Warningf("retrying certificate fetch")
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
				InsecureSkipVerify: true, // this is deliberate, might only have self signed certificate
			},
		},
		Timeout: time.Duration(30) * time.Second,
	}
	return client
}
