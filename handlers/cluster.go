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

package handlers

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/operations/cluster"
)

// CreateClusterHandlerImpl implementation of the CreateClusterHandler interface
type CreateClusterHandlerImpl struct {
	Client      client_native.HAProxyClient
	Config      *configuration.Configuration
	ReloadAgent haproxy.IReloadAgent
}

// GetClusterHandlerImpl implementation of the GetClusterHandler interface
type GetClusterHandlerImpl struct {
	Config *configuration.Configuration
}

// ClusterInitiateCertificateRefreshHandlerImpl implementation of the ClusterInitiateCertificateRefreshHandler interface
type ClusterInitiateCertificateRefreshHandlerImpl struct {
	Config *configuration.Configuration
}

type DeleteClusterHandlerImpl struct {
	Client      client_native.HAProxyClient
	Config      *configuration.Configuration
	Users       *configuration.Users
	ReloadAgent haproxy.IReloadAgent
}

type EditClusterHandlerImpl struct {
	Config *configuration.Configuration
}

// Handle executing the request and returning a response
func (h *ClusterInitiateCertificateRefreshHandlerImpl) Handle(params cluster.InitiateCertificateRefreshParams, principal interface{}) middleware.Responder {
	if h.Config.Mode.Load() != configuration.ModeCluster {
		return cluster.NewInitiateCertificateRefreshForbidden()
	}
	h.Config.Notify.CertificateRefresh.Notify()
	return cluster.NewInitiateCertificateRefreshOK()
}

func (h *CreateClusterHandlerImpl) err500(err error, transaction *models.Transaction) middleware.Responder {
	if transaction != nil {
		configuration, confErr := h.Client.Configuration()
		if confErr == nil {
			_ = configuration.DeleteTransaction(transaction.ID)
		}
	}
	msg := err.Error()
	code := int64(500)
	return cluster.NewPostClusterDefault(500).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func (h *CreateClusterHandlerImpl) err406(err error, transaction *models.Transaction) middleware.Responder {
	// 406 Not Acceptable
	if transaction != nil {
		configuration, confErr := h.Client.Configuration()
		if confErr == nil {
			_ = configuration.DeleteTransaction(transaction.ID)
		}
	}
	msg := err.Error()
	code := int64(406)
	return cluster.NewPostClusterDefault(406).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func (h *CreateClusterHandlerImpl) err409(err error, transaction *models.Transaction) middleware.Responder {
	// 409 Conflict
	if transaction != nil {
		configuration, confErr := h.Client.Configuration()
		if confErr == nil {
			_ = configuration.DeleteTransaction(transaction.ID)
		}
	}
	msg := err.Error()
	code := int64(409)
	return cluster.NewPostClusterDefault(409).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func (h *CreateClusterHandlerImpl) Handle(params cluster.PostClusterParams, principal interface{}) middleware.Responder {
	key := h.Config.Cluster.BootstrapKey.Load()
	if params.Data.BootstrapKey != "" && key != params.Data.BootstrapKey {
		// before we switch to cluster mode, check if folder for storage is compatible with dataplane
		key, err := configuration.DecodeBootstrapKey(params.Data.BootstrapKey)
		if err != nil {
			return h.err406(err, nil)
		}
		log.Warningf("received instructions from %s to join cluster %s at %s", params.HTTPRequest.RemoteAddr, key["name"], key["address"])
		errStorageDir := configuration.CheckIfStorageDirIsOK(key["storage-dir"], h.Config)
		if errStorageDir != nil {
			log.Warningf("configured storage dir incompatible with cluster configuration: %s", errStorageDir)
			return h.err409(errStorageDir, nil)
		}
		// Init NOTICE file to inform user that the cluster storage folder is programmatically managed by Fusion API
		if errStorageInit := configuration.InitStorageNoticeFile(key["storage-dir"]); errStorageInit != nil {
			log.Warningf("unable to create notice file, %s: skipping it", errStorageInit.Error())
		}
		// enforcing API advertising options
		if a := params.AdvertisedAddress; a != nil {
			h.Config.APIOptions.APIAddress = *a
		}
		if p := params.AdvertisedPort; p != nil {
			h.Config.APIOptions.APIPort = *p
		}
		h.Config.Mode.Store(configuration.ModeCluster)
		h.Config.Cluster.BootstrapKey.Store(params.Data.BootstrapKey)
		h.Config.Cluster.Clear()
		// ensuring configuration file saving occurs before notifying the monitor about the bootstrap key change
		defer func() {
			h.Config.Notify.BootstrapKeyChanged.Notify()
		}()
	}
	err := h.Config.Save()
	if err != nil {
		return h.err500(err, nil)
	}
	return cluster.NewPostClusterOK().WithPayload(getClusterSettings(h.Config))
}

// Handle executing the request and returning a response
func (h *GetClusterHandlerImpl) Handle(params cluster.GetClusterParams, principal interface{}) middleware.Responder {
	return cluster.NewGetClusterOK().WithPayload(getClusterSettings(h.Config))
}

func (h *DeleteClusterHandlerImpl) Handle(params cluster.DeleteClusterParams, principal interface{}) middleware.Responder {
	log.Warningf("received instructions from %s to switch to standalone mode", params.HTTPRequest.RemoteAddr)
	// Only do when dataplane is in cluster mode, if not, do nothing and return 204
	if h.Config.Mode.Load() == configuration.ModeCluster {
		log.Warning("clearing cluster users")
		for _, u := range h.Users.GetUsers() {
			// remove all users for cluster communication
			if strings.HasPrefix(u.Name, "dpapi-c-") {
				errRU := h.Users.RemoveUser(u)
				if errRU != nil {
					log.Error(errRU.Error())
				}
			}
		}
		// If we don't want to keep the haproxy configuration, set it to dummy config
		if params.Configuration == nil || *params.Configuration != "keep" {
			log.Warning("clearing configuration as requested")

			conf, err := h.Client.Configuration()
			if err != nil {
				return h.err500(err, nil)
			}
			version, errVersion := conf.GetVersion("")
			if errVersion != nil || version < 1 {
				// silently fallback to 1
				version = 1
			}

			config := fmt.Sprintf(DummyConfig, time.Now().Format("01-02-2006 15:04:05 MST"), h.Config.Name.Load())
			if err = conf.PostRawConfiguration(&config, version, true); err != nil {
				return h.err500(err, nil)
			}

			// we need to restart haproxy
			err = h.ReloadAgent.Restart()
			if err != nil {
				return h.err500(err, nil)
			}
			// Deleting the storage directory used by Fusion:
			// avoiding at all entering any nil pointer dereference.
			if storageData := h.Config.GetStorageData(); storageData != nil && storageData.Cluster != nil && storageData.Cluster.StorageDir != nil {
				if storageErr := configuration.RemoveStorageFolder(*storageData.Cluster.StorageDir); storageErr != nil {
					log.Warningf("failed to clean-up the cluster storage directory: %s", storageErr.Error())
				}
			}
		}
		h.Config.Cluster.BootstrapKey.Store("")
		h.Config.Mode.Store(configuration.ModeSingle)
		h.Config.Status.Store("active")
		h.Config.Cluster.Clear()
		defer func() {
			log.Warning("reloading to apply configuration changes")
			h.Config.Notify.Reload.Notify()
		}()
	}
	err := h.Config.Save()
	if err != nil {
		return h.err500(err, nil)
	}
	return cluster.NewDeleteClusterNoContent()
}

func (h *DeleteClusterHandlerImpl) err500(err error, transaction *models.Transaction) middleware.Responder {
	if transaction != nil {
		configuration, confErr := h.Client.Configuration()
		if confErr == nil {
			_ = configuration.DeleteTransaction(transaction.ID)
		}
	}
	msg := err.Error()
	code := int64(500)
	return cluster.NewDeleteClusterDefault(500).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func (h *EditClusterHandlerImpl) Handle(params cluster.EditClusterParams, principal interface{}) middleware.Responder {
	// Only do when dataplane is in cluster mode, if not, do nothing and return 204
	if h.Config.Mode.Load() == configuration.ModeCluster {
		// for now change only cluster log targets in PUT method
		if params.Data != nil && params.Data.Cluster != nil {
			if clusterLogTargetsChanged(h.Config.Cluster.ClusterLogTargets, params.Data.Cluster.ClusterLogTargets) {
				h.Config.Cluster.ClusterLogTargets = params.Data.Cluster.ClusterLogTargets
				h.Config.Cluster.ClusterID.Store(params.Data.Cluster.ClusterID)
				err := h.Config.Save()
				if err != nil {
					return h.err500(err)
				}
				defer h.Config.Notify.Reload.Notify()
			}
		}
		return cluster.NewEditClusterOK().WithPayload(getClusterSettings(h.Config))
	}
	return h.err406(fmt.Errorf("dataplaneapi in single mode"))
}

func (h *EditClusterHandlerImpl) err406(err error) middleware.Responder {
	// 406 Not Acceptable
	msg := err.Error()
	code := int64(406)
	return cluster.NewEditClusterDefault(406).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func (h *EditClusterHandlerImpl) err500(err error) middleware.Responder {
	msg := err.Error()
	code := int64(500)
	return cluster.NewEditClusterDefault(500).WithPayload(&models.Error{
		Code:    &code,
		Message: &msg,
	})
}

func getClusterSettings(cfg *configuration.Configuration) *models.ClusterSettings {
	portStr := cfg.Cluster.Port.Load()
	port := int64(portStr)
	var clusterSettings *models.ClusterSettingsCluster
	if cfg.Mode.Load() == configuration.ModeCluster {
		clusterSettings = &models.ClusterSettingsCluster{
			Address:           cfg.Cluster.URL.Load(),
			Port:              &port,
			APIBasePath:       cfg.Cluster.APIBasePath.Load(),
			Name:              cfg.Cluster.Name.Load(),
			Description:       cfg.Cluster.Description.Load(),
			ClusterLogTargets: cfg.Cluster.ClusterLogTargets,
		}
	}
	settings := &models.ClusterSettings{
		BootstrapKey: cfg.Cluster.BootstrapKey.Load(),
		Cluster:      clusterSettings,
		Mode:         cfg.Mode.Load(),
		Status:       cfg.Status.Load(),
	}
	return settings
}

func clusterLogTargetsChanged(oldCLT []*models.ClusterLogTarget, newCLT []*models.ClusterLogTarget) bool {
	if len(oldCLT) == len(newCLT) {
		eqCtr := 0
		for _, oldT := range oldCLT {
			for _, newT := range newCLT {
				if reflect.DeepEqual(oldT, newT) {
					eqCtr++
				}
			}
		}
		return !(eqCtr == len(oldCLT))
	}
	return true
}

const DummyConfig = `# NOTE: This configuration file was managed by the Fusion Control Plane.
# Fusion released the control at %s

defaults
  mode http
  timeout connect 5000
  timeout client 30000
  timeout server 10000

frontend disabled
  bind /tmp/dataplaneapi-%s.sock name tmp

`
