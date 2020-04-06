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
	"time"

	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/config-parser/v2/types"
	"github.com/haproxytech/dataplaneapi/haproxy"
	log "github.com/sirupsen/logrus"
)

//Node is structure required for connection to cluster
type Node struct {
	Address     string `json:"address"`
	APIBasePath string `json:"api_base_path"`
	APIPassword string `json:"api_password"`
	APIUser     string `json:"api_user"`
	Certificate string `json:"certificate,omitempty"`
	Description string `json:"description,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Port        int64  `json:"port,omitempty"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

//ClusterSync fetches certificates for joining cluster
type ClusterSync struct {
	cfg         *Configuration
	certFetch   chan struct{}
	cli         *client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c *ClusterSync) Monitor(cfg *Configuration, cli *client_native.HAProxyClient) {
	c.cfg = cfg
	c.cli = cli

	go c.monitorBootstrapKey()

	c.certFetch = make(chan struct{}, 2)
	go c.fetchCert()

	key := c.cfg.BootstrapKey.Load()
	certFetched := cfg.Cluster.CertFetched.Load()

	if key != "" && !certFetched {
		c.cfg.Notify.BootstrapKeyChanged.Notify()
	}
}

func (c *ClusterSync) monitorBootstrapKey() {
	for range c.cfg.Notify.BootstrapKeyChanged.Subscribe("monitorBootstrapKey") {
		key := c.cfg.BootstrapKey.Load()
		c.cfg.Cluster.CertFetched.Store(false)
		if key == "" {
			//do we need to delete cert here maybe?
			c.cfg.Cluster.ActiveBootstrapKey.Store("")
			err := c.cfg.Save()
			if err != nil {
				log.Panic(err)
			}
			continue
		}
		if key == c.cfg.Cluster.ActiveBootstrapKey.Load() {
			fetched := c.cfg.Cluster.CertFetched.Load()
			if !fetched {
				c.certFetch <- struct{}{}
			}
			continue
		}
		data, err := decodeBootstrapKey(key)
		if err != nil {
			log.Warning(err)
		}
		if len(data) != 8 {
			log.Warning("bottstrap key in unrecognized format")
			continue
		}
		url := fmt.Sprintf("%s://%s", data[0], data[1])
		c.cfg.Cluster.URL.Store(url)
		c.cfg.Cluster.Port.Store(data[2])
		c.cfg.Cluster.APIBasePath.Store(data[3])
		c.cfg.Cluster.APINodesPath.Store(data[4])
		c.cfg.Cluster.Name.Store(data[5])
		c.cfg.Cluster.Description.Store(data[6])
		c.cfg.Mode.Store("cluster")
		err = c.cfg.Save()
		if err != nil {
			log.Panic(err)
		}
		csr, key, err := generateCSR()
		if err != nil {
			log.Warning(err)
			continue
		}
		err = ioutil.WriteFile(c.cfg.Cluster.CertificateKeyPath.Load(), []byte(key), 0644)
		if err != nil {
			log.Warning(err)
			continue
		}
		err = ioutil.WriteFile(c.cfg.Cluster.CertificateCSR.Load(), []byte(csr), 0644)
		if err != nil {
			log.Warning(err)
			continue
		}
		err = c.cfg.Save()
		if err != nil {
			log.Panic(err)
		}
		err = c.issueJoinRequest(url, data[2], data[3], data[4], csr, key)
		if err != nil {
			log.Warning(err)
			continue
		}
		c.certFetch <- struct{}{}
	}
}

func (c *ClusterSync) issueJoinRequest(url, port, basePath string, nodesPath string, csr, key string) error {
	url = fmt.Sprintf("%s:%s%s/%s", url, port, basePath, nodesPath)
	serverCfg := c.cfg.Server
	users := GetUsersStore().GetUsers()
	if len(users) == 0 {
		return fmt.Errorf("no users configured in %v userlist in conf", c.cfg.HAProxy.Userlist)
	}
	var user *types.User
	for _, u := range users {
		if u.IsInsecure {
			user = &u
			break
		}
	}
	if user == nil {
		return fmt.Errorf("no available user for cluster comunication")
	}

	nodeData := Node{
		//ID:          "",
		Address:     serverCfg.Host,
		APIBasePath: serverCfg.APIBasePath,
		APIPassword: user.Password,
		APIUser:     user.Name,
		Certificate: csr,
		Description: "",
		Name:        c.cfg.Name.Load(),
		Port:        int64(serverCfg.Port),
		Status:      "waiting_approval",
		Type:        "community",
	}
	bytesRepresentation, _ := json.Marshal(nodeData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return fmt.Errorf("error creating new POST request for cluster comunication")
	}
	req.Header.Add("X-Bootstrap-Key", c.cfg.BootstrapKey.Load())
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
		//write id to file
		errFID := ioutil.WriteFile(c.cfg.HAProxy.NodeIDFile, []byte(responseData.ID), 0644)
		if errFID != nil {
			return errFID
		}
		version, errVersion := c.cli.Configuration.GetVersion("")
		if errVersion != nil || version < 1 {
			//silently fallback to 1
			version = 1
		}
		t, err1 := c.cli.Configuration.StartTransaction(version)
		if err1 != nil {
			return err1
		}
		//write id to peers
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
		//restart HAProxy
		errRestart := c.ReloadAgent.Restart()
		if errRestart != nil {
			return errRestart
		}
	}
	c.cfg.Cluster.ID.Store(responseData.ID)
	c.cfg.Cluster.Token.Store(resp.Header.Get("X-Node-Key"))
	c.cfg.Cluster.ActiveBootstrapKey.Store(c.cfg.BootstrapKey.Load())
	c.cfg.Status.Store(responseData.Status)
	log.Infof("Cluster joined, status: %s", responseData.Status)
	if responseData.Status == "active" {
		err = ioutil.WriteFile(c.cfg.Cluster.CertificatePath.Load(), []byte(responseData.Certificate), 0644)
		if err != nil {
			return err
		}
		c.cfg.Cluster.CertFetched.Store(true)
		c.cfg.Notify.Reload.Notify()
	}
	err = c.cfg.Save()
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterSync) activateFetchCert(err error) {
	go func(err error) {
		log.Warning(err)
		time.Sleep(1 * time.Minute)
		c.certFetch <- struct{}{}
	}(err)
}

func (c *ClusterSync) fetchCert() {
	for range c.certFetch {
		key := c.cfg.BootstrapKey.Load()
		if key == "" || c.cfg.Cluster.Token.Load() == "" {
			continue
		}
		//if not, sleep and start all over again
		certFetched := c.cfg.Cluster.CertFetched.Load()
		if !certFetched {
			url := c.cfg.Cluster.URL.Load()
			port := c.cfg.Cluster.Port.Load()
			apiBasePath := c.cfg.Cluster.APIBasePath.Load()
			apiNodesPath := c.cfg.Cluster.APINodesPath.Load()
			id := c.cfg.Cluster.ID.Load()
			url = fmt.Sprintf("%s:%s/%s/%s/%s", url, port, apiBasePath, apiNodesPath, id)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				c.activateFetchCert(err)
				continue
			}
			req.Header.Add("X-Node-Key", c.cfg.Cluster.Token.Load())
			req.Header.Add("Content-Type", "application/json")
			httpClient := createHTTPClient()
			resp, err := httpClient.Do(req)
			if err != nil {
				c.activateFetchCert(err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				c.activateFetchCert(err)
				continue
			}
			if resp.StatusCode != 200 {
				c.activateFetchCert(fmt.Errorf("status code not proper [%d] %s", resp.StatusCode, string(body)))
				continue
			}
			var responseData Node
			err = json.Unmarshal(body, &responseData)
			if err != nil {
				c.activateFetchCert(err)
				continue
			}
			c.cfg.Status.Store(responseData.Status)
			log.Warningf("Fetching certificate, status: %s", responseData.Status)

			if responseData.Status == "active" {
				err = ioutil.WriteFile(c.cfg.Cluster.CertificatePath.Load(), []byte(responseData.Certificate), 0644)
				if err != nil {
					log.Warning(err.Error())
					continue
				}
				c.cfg.Cluster.CertFetched.Store(true)
				c.cfg.Notify.Reload.Notify()
			}
			err = c.cfg.Save()
			if err != nil {
				log.Warning(err)
			}
		}
		if !certFetched {
			go func() {
				time.Sleep(1 * time.Minute)
				c.certFetch <- struct{}{}
			}()
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
				InsecureSkipVerify: true, //this is deliberate, might only have self signed certificate
			},
		},
		Timeout: time.Duration(10) * time.Second,
	}
	return client
}
