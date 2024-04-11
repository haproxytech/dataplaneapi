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

package resilient

import (
	"context"
	"errors"

	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/runtime"
	cn "github.com/haproxytech/dataplaneapi/client-native"
	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
)

type Client struct {
	client_native.HAProxyClient
}

func NewClient(c client_native.HAProxyClient) *Client {
	return &Client{
		c,
	}
}

// Runtime is a wrapper around HAProxyClient.Runtime
// that retries once to configure the runtime client if it failed
func (c *Client) Runtime() (runtime.Runtime, error) {
	runtime, err := c.HAProxyClient.Runtime()

	// We already have a valid runtime
	// Let's return it
	if err == nil {
		return runtime, nil
	}

	// Now, for let's try to reconfigure once the runtime
	cfg, err := c.HAProxyClient.Configuration()
	if err != nil {
		return nil, err
	}

	dpapiCfg := dataplaneapi_config.Get()
	haproxyOptions := dpapiCfg.HAProxy

	// Let's disable the delayed start by putting a max value to 0
	// This is important to not block the handlers by waiting the DelayedStartMax that we wait for when we start
	haproxyOptions.DelayedStartMax = 0
	// let's retry
	rnt := cn.ConfigureRuntimeClient(context.Background(), cfg, haproxyOptions)
	if rnt == nil {
		return nil, errors.New("retry - unable to configure runtime client")
	}
	c.HAProxyClient.ReplaceRuntime(rnt)

	return c.HAProxyClient.Runtime()
}
