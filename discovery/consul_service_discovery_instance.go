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

package discovery

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/log"
	jsoniter "github.com/json-iterator/go"
)

type consulService struct {
	params  *models.Consul
	name    string
	servers []configuration.ServiceServer
	changed bool
}

func (c *consulService) GetName() string {
	return c.name
}

func (c *consulService) GetFrom() string {
	return c.params.Defaults
}

func (c *consulService) GetBackendName() string {
	return fmt.Sprintf("consul-backend-%s-%d-%s", *c.params.Address, *c.params.Port, c.name)
}

func (c *consulService) Changed() bool {
	return c.changed
}

func (c *consulService) GetServers() []configuration.ServiceServer {
	return c.servers
}

type consulInstance struct {
	ctx             context.Context
	params          *models.Consul
	httpClient      *http.Client
	discoveryConfig *ServiceDiscoveryInstance
	prevIndexes     map[string]uint64
	update          chan struct{}
	logFields       map[string]interface{}
	timeout         time.Duration
	prevEnabled     bool
}

func (c *consulInstance) start() error {
	c.logDebug("discovery job starting")
	if err := c.setAPIClient(); err != nil {
		return err
	}
	c.update = make(chan struct{}, 1)
	go c.watch()
	return nil
}

func (c *consulInstance) setAPIClient() error {
	c.httpClient = &http.Client{
		Timeout: time.Minute,
	}
	return nil
}

func (c *consulInstance) watch() {
	watchTimer := time.NewTimer(c.timeout)
	defer watchTimer.Stop()

	for {
		select {
		case _, ok := <-c.update:
			if !ok {
				return
			}
			c.logDebug("discovery job update triggered")
			if err := c.setAPIClient(); err != nil {
				c.logErrorf("error while setting up the API client: %s", err.Error())
				c.stop()
				continue
			}
			err := c.discoveryConfig.UpdateParams(discoveryInstanceParams{
				Allowlist:       c.params.ServiceAllowlist,
				Denylist:        c.params.ServiceDenylist,
				LogFields:       c.logFields,
				ServerSlotsBase: int(*c.params.ServerSlotsBase),
				SlotsGrowthType: *c.params.ServerSlotsGrowthType,
				SlotsIncrement:  int(c.params.ServerSlotsGrowthIncrement),
			})
			if err != nil {
				c.stop()
				c.logErrorf("error while updating the instance: %s", err.Error())
			}
		case <-c.ctx.Done():
			c.stop()
		case <-watchTimer.C:
			c.logDebug("discovery job reconciliation started")
			if err := c.updateServices(); err != nil {
				// c.log.Errorf("error while updating service: %w", err)
				c.stop()
			}
			c.logDebug("discovery job reconciliation completed")
			watchTimer.Reset(c.timeout)
		}
	}
}

func (c *consulInstance) stop() {
	c.logDebug("discovery job stopping")
	c.httpClient = nil
	c.prevEnabled = false
	close(c.update)
}

func (c *consulInstance) updateServices() error {
	services := make([]ServiceInstance, 0)
	params := &queryParams{}
	if c.params.Namespace != "" {
		params.Namespace = c.params.Namespace
	}

	if c.params.ServiceNameRegexp != "" {
		params.Filter = fmt.Sprintf("ServiceName matches \"%s\"", c.params.ServiceNameRegexp)
	}
	cServices, _, err := c.queryCatalogServices(params)
	if err != nil {
		return err
	}
	newIndexes := make(map[string]uint64)
	for se := range cServices {
		if se == "consul" {
			continue
		}
		nodes, meta, err := c.queryHealthService(se, &queryParams{})
		if err != nil {
			continue
		}
		newIndexes[se] = meta.LastIndex
		services = append(services, &consulService{
			name:    se,
			params:  c.params,
			servers: c.convertToServers(nodes),
			changed: c.hasServiceChanged(se, meta.LastIndex),
		})
	}
	c.prevIndexes = newIndexes
	return c.discoveryConfig.UpdateServices(services)
}

func (c *consulInstance) convertToServers(nodes []*serviceEntry) []configuration.ServiceServer {
	servers := make([]configuration.ServiceServer, 0)
	for _, node := range nodes {
		if !c.validateHealthChecks(node) {
			continue
		}
		if node.Service.Address != "" {
			servers = append(servers, configuration.ServiceServer{
				Address: node.Service.Address,
				Port:    node.Service.Port,
			})
		} else {
			servers = append(servers, configuration.ServiceServer{
				Address: node.Node.Address,
				Port:    node.Service.Port,
			})
		}
	}
	return servers
}

func (c *consulInstance) validateHealthChecks(node *serviceEntry) bool {
	if c.params.HealthCheckPolicy == nil {
		return true
	}
	switch *c.params.HealthCheckPolicy {
	case models.ConsulHealthCheckPolicyAny:
		return c.validateHealthChecksAny(node)
	case models.ConsulHealthCheckPolicyAll:
		return c.validateHealthChecksAll(node)
	case models.ConsulHealthCheckPolicyMin:
		return c.validateHealthChecksMin(node)
	case models.ConsulHealthCheckPolicyNone:
		return true
	default:
		return true
	}
}

func (c *consulInstance) validateHealthChecksAny(node *serviceEntry) bool {
	if node.Checks == nil || len(node.Checks) == 0 {
		return false
	}

	for _, check := range node.Checks {
		if check.Status == "passing" {
			return true
		}
	}
	return false
}

func (c *consulInstance) validateHealthChecksAll(node *serviceEntry) bool {
	if node.Checks == nil || len(node.Checks) == 0 {
		return false
	}

	for _, check := range node.Checks {
		if check.Status != "passing" {
			return false
		}
	}
	return true
}

func (c *consulInstance) validateHealthChecksMin(node *serviceEntry) bool {
	if node.Checks == nil || len(node.Checks) == 0 {
		return false
	}

	passing := 0
	for _, check := range node.Checks {
		if check.Status == "passing" {
			passing++
		}
	}
	return passing >= int(c.params.HealthCheckPolicyMin)
}

func (c *consulInstance) hasServiceChanged(service string, index uint64) bool {
	prevIndex, ok := c.prevIndexes[service]
	if !ok {
		return true
	}
	return prevIndex != index
}

func (c *consulInstance) updateTimeout(timeoutSeconds int) error {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", timeoutSeconds))
	if err != nil {
		return err
	}
	c.timeout = timeout
	return nil
}

func (c *consulInstance) handleStateChange() error {
	if c.stateChangedToEnabled() {
		if err := c.start(); err != nil {
			c.prevEnabled = false
			return err
		}
		c.prevEnabled = *c.params.Enabled
		return nil
	}
	if c.stateChangedToDisabled() {
		c.stop()
	}
	if *c.params.Enabled {
		c.update <- struct{}{}
	}
	c.prevEnabled = *c.params.Enabled
	return nil
}

func (c *consulInstance) stateChangedToEnabled() bool {
	return !c.prevEnabled && *c.params.Enabled
}

func (c *consulInstance) stateChangedToDisabled() bool {
	return c.prevEnabled && !*c.params.Enabled
}

func (c *consulInstance) logDebug(message string) {
	log.WithFields(c.logFields, log.DebugLevel, message)
}

func (c *consulInstance) logErrorf(format string, args ...interface{}) {
	log.WithFieldsf(c.logFields, log.ErrorLevel, format, args...)
}

func (c *consulInstance) queryCatalogServices(params *queryParams) (map[string][]string, *queryMetadata, error) {
	nodes := map[string][]string{}
	meta, err := c.doConsulQuery(http.MethodGet, "/v1/catalog/services", params, &nodes)
	if err != nil {
		return nil, nil, err
	}
	return nodes, meta, nil
}

func (c *consulInstance) queryHealthService(se string, params *queryParams) ([]*serviceEntry, *queryMetadata, error) {
	services := []*serviceEntry{}
	path, err := url.JoinPath("/v1/health/service/", se)
	if err != nil {
		return nil, nil, err
	}
	meta, err := c.doConsulQuery(http.MethodGet, path, params, &services)
	if err != nil {
		return nil, nil, err
	}
	return services, meta, nil
}

func (c *consulInstance) doConsulQuery(method string, path string, params *queryParams, resp interface{}) (*queryMetadata, error) {
	mode := "http://"
	if c.params.Mode != nil {
		mode = *c.params.Mode + "://"
	}
	fullPath, err := url.JoinPath(
		mode,
		net.JoinHostPort(*c.params.Address, strconv.FormatInt(*c.params.Port, 10)),
		path,
	)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, fullPath, nil)
	if err != nil {
		return nil, err
	}

	// Global Consul parameters.
	if c.params.Token != "" {
		req.Header.Add("X-Consul-Token", c.params.Token)
	}

	q := url.Values{}

	// Request's parameters.
	if params.Namespace != "" {
		q.Add("ns", c.params.Namespace)
	}

	if params.Filter != "" {
		q.Add("filter", params.Filter)
	}

	req.URL.RawQuery = q.Encode()

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(raw, resp)
	if err != nil {
		return nil, err
	}

	meta, err := extractMetadata(httpResp)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

type serviceEntry struct {
	Node *struct {
		Address string
	}
	Service *struct {
		Address string
		Port    int
	}
	Checks []*struct {
		Status string
	}
}

type queryParams struct {
	Namespace string
	Filter    string
}

type queryMetadata struct {
	LastIndex uint64
}

func extractMetadata(resp *http.Response) (*queryMetadata, error) {
	meta := queryMetadata{}
	indexStr := resp.Header.Get("X-Consul-Index")
	if indexStr != "" {
		lastIndex, err := strconv.ParseUint(indexStr, 10, 64)
		if err != nil {
			return nil, err
		}
		meta.LastIndex = lastIndex
	}
	return &meta, nil
}
