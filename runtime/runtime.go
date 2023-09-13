// Copyright 2023 HAProxy Technologies
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

package runtime

import (
	"context"
	"net"
	"strings"
	"sync"
	"syscall"

	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/log"
	commands "github.com/haproxytech/dataplaneapi/runtime/commands"
)

var (
	debugServerOnce sync.Once
	debugServer     *DebugServer
)

type DebugServer struct {
	Client      client_native.HAProxyClient
	DAPIVersion string
	CnChannel   chan client_native.HAProxyClient
	comm        Commands
}

func GetServer() *DebugServer {
	debugServerOnce.Do(func() {
		debugServer = &DebugServer{}
		debugServer.CnChannel = make(chan client_native.HAProxyClient)
	})
	return debugServer
}

func (s *DebugServer) Start(ctx context.Context, cancelFunc context.CancelFunc) {
	go s.handleClientNativeClientUpdate(ctx)

	cfg := configuration.Get()
	log.Debug("-- command socket waiting for DAPI conf readiness...")
	<-cfg.Notify.ServerStarted.Subscribe("commandSocket")
	log.Debug("-- command socket Server. conf ready...")

	cfg.Load()
	if cfg.HAProxy.DebugSocketPath == "" {
		log.Debug("-- command socket not set (--debug_socket_path). Not running the command socket server")
		cancelFunc()
		return
	}
	_ = syscall.Unlink(cfg.HAProxy.DebugSocketPath) // we don't care if it can't be removed

	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "unix", cfg.HAProxy.DebugSocketPath)
	if err != nil { // but we care if we can't listen on it
		log.Error(err)
		return
	}

	ver := strings.TrimSpace(s.DAPIVersion)

	s.comm.Register(commands.Goroutines{})
	s.comm.Register(commands.Stack{})
	s.comm.Register(commands.DataplaneapiVersion{Version: ver})
	s.comm.Register(commands.PProf{})
	s.comm.Register(commands.DataplaneapiConfiguration{})

	stopped := false

	go func() {
		<-ctx.Done()
		stopped = true
		l.Close()
		log.Info("-- command socket Shutting down...")
	}()

	log.Infof("-- command socket Starting on %s", cfg.HAProxy.DebugSocketPath)
	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				if stopped {
					return
				}
				log.Info("runtime accept error", err.Error())
				return
			}
			go serve(&s.comm, fd)
		}
	}()
}

func (s *DebugServer) handleClientNativeClientUpdate(ctx context.Context) {
	for {
		select {
		case cn, ok := <-s.CnChannel:
			if !ok {
				return
			}
			log.Debug("-- command socket updating client_native client")
			s.comm.UnRegister(commands.ConfCmdKey)
			s.comm.Register(commands.HAProxyConfiguration{Client: cn})
			s.Client = cn
		case <-ctx.Done():
			log.Debug("-- command socket handleClientNativeClientUpdate stopped")
			return
		}
	}
}
