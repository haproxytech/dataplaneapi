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
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
)

const (
	HAProxyServiceNameTag  = "HAProxy:Service:Name"
	HAProxyServicePortTag  = "HAProxy:Service:Port"
	HAProxyInstancePortTag = "HAProxy:Instance:Port"
)

type awsInstance struct {
	ctx             context.Context
	params          *models.AwsRegion
	update          chan struct{}
	state           map[string]map[string]time.Time
	discoveryConfig *ServiceDiscoveryInstance
	logFields       map[string]interface{}
	timeout         time.Duration
}

type awsService struct {
	instances    map[string]types.Instance
	name         string
	region       string
	instanceName string
	ipv4         string
	changed      bool
}

func (a awsService) GetName() string {
	return a.name
}

func (a awsService) GetFrom() string {
	return ""
}

func (a awsService) GetBackendName() string {
	return fmt.Sprintf("aws-%s-%s-%s", a.region, a.instanceName, a.GetName())
}

func (a awsService) Changed() bool {
	return a.changed
}

func (a awsService) GetServers() (servers []configuration.ServiceServer) {
	for _, instance := range a.instances {
		port, _ := a.instancePortFromEC2(instance)
		var address string
		switch a.ipv4 {
		case models.AwsRegionIPV4AddressPrivate:
			address = aws.ToString(instance.PrivateIpAddress)
		case models.AwsRegionIPV4AddressPublic:
			address = aws.ToString(instance.PublicIpAddress)
		default:
			continue
		}
		// In case of public IPv4 and the instance doesn't have it, ignoring.
		if len(address) == 0 {
			continue
		}
		servers = append(servers, configuration.ServiceServer{
			Address: address,
			Port:    port,
		})
	}
	return
}

func newAWSRegionInstance(ctx context.Context, params *models.AwsRegion, client configuration.Configuration, reloadAgent haproxy.IReloadAgent) (*awsInstance, error) {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", *params.RetryTimeout))
	if err != nil {
		return nil, err
	}

	logFields := map[string]interface{}{"ServiceDiscovery": "AWS", "ID": *params.ID}

	ai := &awsInstance{
		params:    params,
		timeout:   timeout,
		ctx:       ctx,
		logFields: logFields,
		state:     make(map[string]map[string]time.Time),
		discoveryConfig: NewServiceDiscoveryInstance(client, reloadAgent, discoveryInstanceParams{
			Allowlist:       []string{},
			Denylist:        []string{},
			LogFields:       logFields,
			ServerSlotsBase: int(*params.ServerSlotsBase),
			SlotsGrowthType: *params.ServerSlotsGrowthType,
			SlotsIncrement:  int(params.ServerSlotsGrowthIncrement),
		}),
	}
	if err = ai.updateTimeout(*params.RetryTimeout); err != nil {
		return nil, err
	}

	return ai, nil
}

func (a *awsInstance) filterConverter(in []*models.AwsFilters) (out []types.Filter) {
	out = make([]types.Filter, len(in))
	for i, l := range in {
		filter := l
		out[i] = types.Filter{
			Name:   filter.Key,
			Values: []string{aws.ToString(filter.Value)},
		}
	}
	return
}

func (a *awsInstance) updateTimeout(timeoutSeconds int64) error {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", timeoutSeconds))
	if err != nil {
		return err
	}
	a.timeout = timeout
	return nil
}

func (a *awsInstance) start() {
	a.update = make(chan struct{})

	go func() {
		a.logDebug("discovery job starting")

		discoveryTimer := time.NewTimer(a.timeout)
		defer discoveryTimer.Stop()

		for {
			select {
			case _, ok := <-a.update:
				if !ok {
					return
				}
				a.logDebug("discovery job update triggered")
				err := a.discoveryConfig.UpdateParams(discoveryInstanceParams{
					Allowlist:       []string{},
					Denylist:        []string{},
					LogFields:       a.logFields,
					ServerSlotsBase: int(*a.params.ServerSlotsBase),
					SlotsGrowthType: *a.params.ServerSlotsGrowthType,
					SlotsIncrement:  int(a.params.ServerSlotsGrowthIncrement),
				})
				if err != nil {
					a.stop()
				}
			case <-discoveryTimer.C:
				a.logDebug("discovery job update triggered")

				var api *ec2.Client
				var err error

				if api, err = a.setAPIClient(); err != nil {
					a.logErrorf("error while setting up the API client: %s", err.Error())
					a.stop()
				}
				if err = a.updateServices(api); err != nil {
					switch t := err.(type) {
					case *configuration.ConfError:
						switch t.Code() {
						case configuration.ErrObjectAlreadyExists:
							continue
						default:
							a.stop()
							a.logErrorf("error while updating service: %s", err.Error())
						}
					default:
						a.stop()
					}
				}

				a.logDebug("discovery job reconciliation completed")
				discoveryTimer.Reset(a.timeout)
			case <-a.ctx.Done():
				a.stop()
			}
		}
	}()
}

func (a *awsInstance) setAPIClient() (*ec2.Client, error) {
	opts := []func(options *config.LoadOptions) error{
		config.WithRegion(*a.params.Region),
	}
	if len(a.params.AccessKeyID) > 0 && len(a.params.SecretAccessKey) > 0 {
		opts = append(opts, config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     a.params.AccessKeyID,
				SecretAccessKey: a.params.SecretAccessKey,
			},
		}))
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot generate the AWS instance due to a configuration setup error: %w", err)
	}

	return ec2.NewFromConfig(cfg), nil
}

func (a *awsInstance) updateServices(api *ec2.Client) (err error) {
	var io *ec2.DescribeInstancesOutput

	io, err = api.DescribeInstances(a.ctx, &ec2.DescribeInstancesInput{
		Filters: append([]types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{HAProxyServiceNameTag, HAProxyServicePortTag},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		}, a.filterConverter(a.params.Allowlist)...),
	})
	if err != nil {
		return
	}

	mapService := make(map[string]*awsService)

	for _, r := range io.Reservations {
		for _, i := range r.Instances {
			var sn string
			sn, err = a.serviceNameFromEC2(i)
			if err != nil {
				a.logErrorf("unable to retrieve service name for the instance %s", *i.InstanceId)

				continue
			}
			// creating empty service in case it isn't there
			if _, ok := mapService[sn]; !ok {
				mapService[sn] = &awsService{
					name:         sn,
					region:       *a.params.Region,
					instanceName: *a.params.Name,
					ipv4:         *a.params.IPV4Address,
					instances:    make(map[string]types.Instance),
				}
			}
			instanceID := aws.ToString(i.InstanceId)

			if _, portErr := mapService[sn].instancePortFromEC2(i); portErr != nil {
				a.logErrorf("unable to retrieve service port for the instance %s", *i.InstanceId)

				continue
			}

			mapService[sn].instances[instanceID] = i
		}
	}

	if len(a.params.Denylist) > 0 {
		// AWS API doesn't provide negative filter search, so doing on our own
		io, err = api.DescribeInstances(a.ctx, &ec2.DescribeInstancesInput{
			Filters: a.filterConverter(a.params.Denylist),
		})
		if err == nil {
			for _, r := range io.Reservations {
				for _, i := range r.Instances {
					var sn string
					sn, err = a.serviceNameFromEC2(i)
					// definitely we can skip, there's no Service metadata tag
					if err != nil {
						continue
					}
					// neither tracked as Service, we can skip
					if _, ok := mapService[sn]; !ok {
						continue
					}
					// we have an occurrence, we have to delete
					instanceID := aws.ToString(i.InstanceId)
					delete(mapService[sn].instances, instanceID)
				}
			}
		}
	}

	var services []ServiceInstance
	for _, s := range mapService {
		// We don't have a proper way to understand if a Service has changed, or not, this can be achieved
		// iterating over the instances being part of the Service and check the last launch time:
		// if something differs, a change occurred.
		s.changed = func() bool {
			if _, ok := a.state[s.name]; !ok {
				return true
			}
			if len(a.state[s.name]) != len(s.instances) {
				return true
			}
			for _, instance := range s.instances {
				instanceID := aws.ToString(instance.InstanceId)
				v, ok := a.state[s.name][instanceID]
				if !ok {
					return true
				}
				if v != *instance.LaunchTime {
					return true
				}
			}
			return false
		}()
		services = append(services, s)

		a.state[s.name] = func(instances map[string]types.Instance) (hash map[string]time.Time) {
			hash = make(map[string]time.Time)
			for _, instance := range instances {
				id := aws.ToString(instance.InstanceId)
				hash[id] = aws.ToTime(instance.LaunchTime)
			}
			return
		}(s.instances)
	}

	return a.discoveryConfig.UpdateServices(services)
}

func (a *awsInstance) stop() {
	a.logDebug("discovery job stopping")
	close(a.update)
}

func (a *awsService) instancePortFromEC2(instance types.Instance) (port int, err error) {
	for _, t := range instance.Tags {
		switch {
		case *t.Key == HAProxyServicePortTag:
			port, err = strconv.Atoi(*t.Value)
		case *t.Key == HAProxyInstancePortTag:
			return strconv.Atoi(*t.Value)
		}
	}
	return
}

func (a *awsInstance) serviceNameFromEC2(instance types.Instance) (string, error) {
	var name, port string
L:
	for _, t := range instance.Tags {
		switch {
		case *t.Key == HAProxyServiceNameTag:
			name = aws.ToString(t.Value)
		case *t.Key == HAProxyServicePortTag:
			port = aws.ToString(t.Value)
		case len(name) > 0 && len(port) > 0:
			break L
		}
	}

	if len(name) == 0 || len(port) == 0 {
		return "", fmt.Errorf("missing metadata for instance %s", *instance.InstanceId)
	}

	return fmt.Sprintf("%s-%s", name, port), nil
}

func (a *awsInstance) logDebug(message string) {
	log.WithFields(a.logFields, log.DebugLevel, message)
}

func (a *awsInstance) logErrorf(format string, args ...interface{}) {
	log.WithFieldsf(a.logFields, log.ErrorLevel, format, args...)
}
