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

package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

// RegisterRouter registers all cluster routes onto r.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, cfg *dataplaneapi_config.Configuration, users *dataplaneapi_config.Users, ra haproxy.IReloadAgent) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client, Config: cfg, Users: users, ReloadAgent: ra}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for cluster management.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	Config      *dataplaneapi_config.Configuration
	Users       *dataplaneapi_config.Users
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetCluster(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, getClusterSettings(h.Config))
}

func (h *HandlerImpl) PostCluster(w http.ResponseWriter, r *http.Request, params PostClusterParams) {
	var data ClusterSettings
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	key := h.Config.Cluster.BootstrapKey.Load()
	if data.BootstrapKey != "" && key != data.BootstrapKey {
		bootstrapKey, err := dataplaneapi_config.DecodeBootstrapKey(data.BootstrapKey)
		if err != nil {
			h.errStatus(w, http.StatusNotAcceptable, err, nil)
			return
		}
		log.Warningf("received instructions from %s to join cluster %s at %s", r.RemoteAddr, bootstrapKey["name"], bootstrapKey["address"])
		if errStorageDir := dataplaneapi_config.CheckIfStorageDirIsOK(bootstrapKey["storage-dir"], h.Config); errStorageDir != nil {
			log.Warningf("configured storage dir incompatible with cluster configuration: %s", errStorageDir)
			h.errStatus(w, http.StatusConflict, errStorageDir, nil)
			return
		}
		if errStorageInit := dataplaneapi_config.InitStorageNoticeFile(bootstrapKey["storage-dir"]); errStorageInit != nil {
			log.Warningf("unable to create notice file, %s: skipping it", errStorageInit.Error())
		}
		if params.AdvertisedAddress != "" {
			h.Config.APIOptions.APIAddress = params.AdvertisedAddress
		}
		if params.AdvertisedPort != 0 {
			h.Config.APIOptions.APIPort = int64(params.AdvertisedPort)
		}
		h.Config.Mode.Store(dataplaneapi_config.ModeCluster)
		h.Config.Cluster.BootstrapKey.Store(data.BootstrapKey)
		h.Config.Cluster.Clear()
		defer func() {
			h.Config.Notify.BootstrapKeyChanged.Notify()
		}()
	}
	if err := h.Config.SaveClusterModeData(); err != nil {
		h.errStatus(w, http.StatusInternalServerError, err, nil)
		return
	}
	respond.JSON(w, http.StatusOK, getClusterSettings(h.Config))
}

func (h *HandlerImpl) EditCluster(w http.ResponseWriter, r *http.Request, params EditClusterParams) {
	if h.Config.Mode.Load() != dataplaneapi_config.ModeCluster {
		h.errStatus(w, http.StatusNotAcceptable, errors.New("dataplaneapi in single mode"), nil)
		return
	}
	var data ClusterSettings
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	if data.Cluster != nil {
		if clusterLogTargetsChanged(h.Config.Cluster.ClusterLogTargets, data.Cluster.ClusterLogTargets) {
			h.Config.Cluster.ClusterLogTargets = data.Cluster.ClusterLogTargets
			h.Config.Cluster.ClusterID.Store(data.Cluster.ClusterID)
			if err := h.Config.Save(); err != nil {
				h.errStatus(w, http.StatusInternalServerError, err, nil)
				return
			}
			defer h.Config.Notify.Reload.Notify()
		}
	}
	respond.JSON(w, http.StatusOK, getClusterSettings(h.Config))
}

func (h *HandlerImpl) DeleteCluster(w http.ResponseWriter, r *http.Request, params DeleteClusterParams) {
	log.Warningf("received instructions from %s to switch to standalone mode", r.RemoteAddr)
	if h.Config.Mode.Load() == dataplaneapi_config.ModeCluster {
		log.Warning("clearing cluster users")
		for _, u := range h.Users.GetUsers() {
			if strings.HasPrefix(u.Name, storagetype.DapiClusterUserPrefix) {
				if errRU := h.Users.RemoveUser(u); errRU != nil {
					log.Error(errRU.Error())
				}
			}
		}
		if params.Configuration != "keep" {
			log.Warning("clearing configuration as requested")
			conf, err := h.Client.Configuration()
			if err != nil {
				h.errStatus(w, http.StatusInternalServerError, err, nil)
				return
			}
			version, errVersion := conf.GetVersion("")
			if errVersion != nil || version < 1 {
				version = 1
			}
			config := fmt.Sprintf(dummyConfig, time.Now().Format("01-02-2006 15:04:05 MST"), h.Config.Name.Load())
			if err = conf.PostRawConfiguration(&config, version, true); err != nil {
				h.errStatus(w, http.StatusInternalServerError, err, nil)
				return
			}
			if err = h.ReloadAgent.Restart(); err != nil {
				h.errStatus(w, http.StatusInternalServerError, err, nil)
				return
			}
			if storageData := h.Config.GetStorageData(); storageData != nil && storageData.DeprecatedCluster != nil && storageData.DeprecatedCluster.StorageDir != nil {
				if storageErr := dataplaneapi_config.RemoveStorageFolder(*storageData.DeprecatedCluster.StorageDir); storageErr != nil {
					log.Warningf("failed to clean-up the cluster storage directory: %s", storageErr.Error())
				}
			}
		}
		h.Config.Cluster.BootstrapKey.Store("")
		h.Config.Status.Store("active")
		h.Config.HAProxy.ClusterTLSCertDir = ""
		h.Config.Cluster.Clear()
		defer func() {
			log.Warning("reloading to apply configuration changes")
			h.Config.Notify.Reload.Notify()
		}()
	}
	if err := h.Config.SaveClusterModeData(); err != nil {
		h.errStatus(w, http.StatusInternalServerError, err, nil)
		return
	}
	respond.NoContent(w)
}

func (h *HandlerImpl) InitiateCertificateRefresh(w http.ResponseWriter, r *http.Request) {
	if h.Config.Mode.Load() != dataplaneapi_config.ModeCluster {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	h.Config.Notify.CertificateRefresh.Notify()
	respond.JSON(w, http.StatusOK, nil)
}

func (h *HandlerImpl) errStatus(w http.ResponseWriter, status int, err error, _ any) {
	code := int64(status)
	msg := err.Error()
	respond.JSON(w, status, &models.Error{Code: &code, Message: &msg})
}

func getClusterSettings(cfg *dataplaneapi_config.Configuration) *models.ClusterSettings {
	portStr := cfg.Cluster.Port.Load()
	port := int64(portStr)
	var clusterSettings *models.ClusterSettingsCluster
	if cfg.Mode.Load() == dataplaneapi_config.ModeCluster {
		clusterSettings = &models.ClusterSettingsCluster{
			Address:           cfg.Cluster.URL.Load(),
			Port:              &port,
			APIBasePath:       cfg.Cluster.APIBasePath.Load(),
			Name:              cfg.Cluster.Name.Load(),
			Description:       cfg.Cluster.Description.Load(),
			ClusterLogTargets: cfg.Cluster.ClusterLogTargets,
		}
	}
	return &models.ClusterSettings{
		BootstrapKey: cfg.Cluster.BootstrapKey.Load(),
		Cluster:      clusterSettings,
		Mode:         cfg.Mode.Load(),
		Status:       cfg.Status.Load(),
	}
}

func clusterLogTargetsChanged(oldCLT, newCLT []*models.ClusterLogTarget) bool {
	if len(oldCLT) != len(newCLT) {
		return true
	}
	eqCtr := 0
	for _, oldT := range oldCLT {
		for _, newT := range newCLT {
			if reflect.DeepEqual(oldT, newT) {
				eqCtr++
			}
		}
	}
	return eqCtr != len(oldCLT)
}

const dummyConfig = `# NOTE: This configuration file was managed by the Fusion Control Plane.
# Fusion released the control at %s

defaults
  mode http
  timeout connect 5000
  timeout client 30000
  timeout server 10000

frontend disabled
  bind /tmp/dataplaneapi-%s.sock name tmp

`
