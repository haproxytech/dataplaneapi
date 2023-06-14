//go:build aws
// +build aws

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
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/google/uuid"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/assert"

	"github.com/haproxytech/dataplaneapi/haproxy"
)

var (
	AWSSecretAccessKey string
	AWSAccessKeyID     string
)

const (
	region        = "us-east-1"
	ami           = "ami-0742b4e673072066f"
	additionalTag = "HAProxy:Integration:Test"
	serviceName   = "my-app"
	servicePort   = 8080
)

func TestAWS(t *testing.T) {
	if len(AWSSecretAccessKey) == 0 {
		t.Fatal("missing AWS Secret Access Key, cannot test on AWS")
	}
	if len(AWSAccessKeyID) == 0 {
		t.Fatal("missing AWS Access Key ID, cannot test on AWS")
	}

	var tmp string
	var err error
	tmp, err = ioutil.TempDir("", "haproxy")
	assert.Nil(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(tmp)
		_ = os.Remove(tmp)
	})

	cfgFile := fmt.Sprintf("%s/haproxy.cfg", tmp)
	_, err = os.Create(cfgFile)
	assert.Nil(t, err)

	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile: cfgFile,
		TransactionDir:    tmp,
	}

	err = confClient.Init(confParams)
	assert.Nil(t, err)

	var ra haproxy.IReloadAgent

	ra, err = haproxy.NewReloadAgent(haproxy.ReloadAgentParams{
		Delay:      1,
		ReloadCmd:  "true",
		RestartCmd: "true",
		ConfigFile: cfgFile,
		BackupDir:  tmp,
		Retention:  0,
		Ctx:        context.Background(),
	})
	assert.Nil(t, err)

	var instance *awsInstance
	instance, err = newAWSRegionInstance(context.Background(), &models.AwsRegion{
		AccessKeyID: AWSAccessKeyID,
		Allowlist: []*models.AwsFilters{
			{
				Key:   aws.String("tag-key"),
				Value: aws.String(additionalTag),
			},
		},
		Denylist:                   []*models.AwsFilters{},
		Description:                "just an integration test on AWS",
		Enabled:                    aws.Bool(true),
		IPV4Address:                aws.String(models.AwsRegionIPV4AddressPrivate),
		ID:                         aws.String(uuid.New().String()),
		Name:                       aws.String("integration-test"),
		Region:                     aws.String(region),
		RetryTimeout:               aws.Int64(1),
		SecretAccessKey:            AWSSecretAccessKey,
		ServerSlotsBase:            aws.Int64(10),
		ServerSlotsGrowthIncrement: 10,
		ServerSlotsGrowthType:      aws.String("linear"),
	}, confClient, ra)
	assert.Nil(t, err)

	var cfg aws.Config
	cfg, err = config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(*instance.params.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     instance.params.AccessKeyID,
				SecretAccessKey: instance.params.SecretAccessKey,
			},
		}),
	)
	assert.Nil(t, err)

	asg := autoscaling.NewFromConfig(cfg)

	err = createLaunchConfiguration(instance.params.ID, asg)

	assert.Nil(t, err)
	t.Cleanup(func() {
		_, _ = asg.DeleteLaunchConfiguration(context.Background(), &autoscaling.DeleteLaunchConfigurationInput{
			LaunchConfigurationName: instance.params.ID,
		})
	})

	err = createAutoScalingGroup(instance.params.ID, asg)
	assert.Nil(t, err)
	t.Cleanup(func() {
		_, _ = asg.DeleteAutoScalingGroup(context.Background(), &autoscaling.DeleteAutoScalingGroupInput{
			AutoScalingGroupName: aws.String(*instance.params.ID),
			ForceDelete:          aws.Bool(true),
		})
	})

	instance.start()
	t.Cleanup(func() {
		instance.stop()
	})

	run := func(t *testing.T, dc int32) {
		err = scaleAutoScalingGroup(instance.params.ID, dc, asg)
		assert.Nil(t, err)

		err = checkAutoScalingGroupCapacity(instance.params.ID, dc, asg)
		assert.Nil(t, err)

		backendName := fmt.Sprintf("aws-%s-%s-%s-%d", region, *instance.params.Name, serviceName, servicePort)
		ec2Client, _ := instance.setAPIClient()
		assert.Eventually(t, func() bool {
			return checkBackendServers(instance.params.ID, backendName, asg, ec2Client, confClient)
		}, 2*time.Minute, time.Second)
	}

	for _, dc := range []int32{1, 5, 10} {
		t.Run(fmt.Sprintf("scaling capacity out of %d servers", dc), func(t *testing.T) {
			run(t, dc)
		})
	}
	for _, dc := range []int32{5, 1, 0} {
		t.Run(fmt.Sprintf("scaling capacity in of %d servers", dc), func(t *testing.T) {
			run(t, dc)
		})
	}
}

func createLaunchConfiguration(name *string, client *autoscaling.Client) (err error) {
	_, err = client.CreateLaunchConfiguration(context.Background(), &autoscaling.CreateLaunchConfigurationInput{
		LaunchConfigurationName:  name,
		AssociatePublicIpAddress: aws.Bool(false),
		ImageId:                  aws.String(ami),
		InstanceType:             aws.String("t2.micro"),
	})
	return
}

func createAutoScalingGroup(instanceId *string, client *autoscaling.Client) (err error) {
	_, err = client.CreateAutoScalingGroup(context.Background(), &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName:    instanceId,
		MaxSize:                 aws.Int32(10),
		MinSize:                 aws.Int32(0),
		DesiredCapacity:         aws.Int32(0),
		AvailabilityZones:       []string{region + "a"},
		LaunchConfigurationName: instanceId,
		Tags: []types.Tag{
			{
				PropagateAtLaunch: aws.Bool(true),
				Key:               aws.String(HAProxyServiceNameTag),
				Value:             aws.String(serviceName),
			},
			{
				PropagateAtLaunch: aws.Bool(true),
				Key:               aws.String(HAProxyServicePortTag),
				Value:             aws.String(fmt.Sprintf("%d", servicePort)),
			},
			{
				PropagateAtLaunch: aws.Bool(true),
				Key:               aws.String(additionalTag),
				Value:             aws.String("true"),
			},
		},
	})
	return
}

func checkBackendServers(asgName *string, backendName string, asg *autoscaling.Client, ec2Client *ec2.Client, confClient *configuration.Client) (ok bool) {
	var out *autoscaling.DescribeAutoScalingGroupsOutput
	var err error
	out, err = asg.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{*asgName},
	})
	if err != nil {
		return false
	}

	_, _, err = confClient.GetBackend(backendName, "")
	if err != nil {
		return false
	}

	var servers models.Servers
	_, servers, err = confClient.GetServers(backendName, "")
	if err != nil {
		return false
	}

	instanceIDs := make([]string, len(out.AutoScalingGroups[0].Instances))
	for k, i := range out.AutoScalingGroups[0].Instances {
		instanceIDs[k] = aws.ToString(i.InstanceId)
	}

	set := make(map[string]string)

	if len(instanceIDs) > 0 {
		instances, _ := ec2Client.DescribeInstances(context.Background(), &ec2.DescribeInstancesInput{
			InstanceIds: instanceIDs,
		})
		for _, r := range instances.Reservations {
			for _, i := range r.Instances {
				id := aws.ToString(i.PrivateIpAddress)
				set[id] = aws.ToString(i.InstanceId)
			}
		}
	}

	var counter int
	for ip := range set {
		for _, server := range servers {
			if ip == server.Address {
				counter++
			}
		}
	}

	return counter == len(set)
}

func scaleAutoScalingGroup(asgName *string, desiredCapacity int32, asg *autoscaling.Client) (err error) {
	_, err = asg.SetDesiredCapacity(context.Background(), &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: asgName,
		DesiredCapacity:      aws.Int32(desiredCapacity),
	})
	return
}

func checkAutoScalingGroupCapacity(asgName *string, desiredCapacity int32, asg *autoscaling.Client) (err error) {
	var out *autoscaling.DescribeAutoScalingGroupsOutput
	ctx, c := context.WithTimeout(context.Background(), 2*time.Minute)
	defer c()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("the desired capacity is not matched")
		case <-time.After(time.Second):
			out, err = asg.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
				AutoScalingGroupNames: []string{*asgName},
			})
			if err != nil {
				continue
			}
			if len(out.AutoScalingGroups[0].Instances) != int(desiredCapacity) {
				continue
			}
			return
		}
	}
}
