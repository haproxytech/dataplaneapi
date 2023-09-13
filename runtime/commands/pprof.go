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

package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	np "net/http/pprof"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/haproxytech/dataplaneapi/log"
)

type PProf struct{}

func (g PProf) Definition() definition {
	return definition{
		Key:  "pprof",
		Info: "pprof dumps",
		Commands: []allCommands{
			{"pprof examples", "pprof usage examples"},
			{"pprof cpu start filename [1h42m]", "start cpu pprof, time is optional default 10 mins"},
			{"pprof cpu stop", "stop cpu pprof"},
			{"pprof mem filename", "dump mem pprof"},
			{"pprof web start port [10m]", "pprof web interface, time is optional default 10 mins"},
			{"pprof web stop", "pprof web interface"},
		},
	}
}

var (
	pprofServer      *http.Server
	mu               sync.Mutex
	stopWebServer    context.CancelFunc
	stopCPUProfiling context.CancelFunc
)

func (g PProf) Command(cmd []string) (response []byte, err error) {
	if len(cmd) < 3 {
		if len(cmd) != 2 || cmd[1] != "examples" {
			return []byte(strings.Trim("", " ")), errors.New("unrecognized command")
		}
	}
	mu.Lock()
	defer mu.Unlock()
	switch cmd[1] {
	case "web":
		switch cmd[2] {
		case "start":
			if stopWebServer != nil {
				return []byte{}, errors.New("pprof web is already running")
			}
			if len(cmd) < 4 {
				return []byte{}, errors.New("port required: pprof web start port")
			}
			sleepTime := time.Minute * 10
			if len(cmd) > 4 {
				sleepTime, err = time.ParseDuration(cmd[4])
				if err != nil {
					return []byte{}, err
				}
			}
			mux := http.NewServeMux()

			mux.HandleFunc("/debug/pprof/", np.Index)
			mux.HandleFunc("/debug/pprof/cmdline", np.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", np.Profile)
			mux.HandleFunc("/debug/pprof/symbol", np.Symbol)
			mux.HandleFunc("/debug/pprof/trace", np.Trace)

			pprofServer = &http.Server{
				Addr:    ":" + cmd[3],
				Handler: mux,
			}

			go func() {
				err = pprofServer.ListenAndServe()
				if err != nil {
					log.Error(err)
				}
			}()

			var ctx context.Context
			ctx, stopWebServer = context.WithCancel(context.Background())

			go func(context.Context) {
				<-ctx.Done()
				mu.Lock()
				defer mu.Unlock()
				ctxServer := context.Background()
				err = pprofServer.Shutdown(ctxServer) //nolint:contextcheck
				if err != nil {
					return
				}
				pprofServer = nil
			}(ctx)

			go func() {
				time.Sleep(sleepTime)
				mu.Lock()
				defer mu.Unlock()
				select {
				case <-ctx.Done():
					return
				default:
				}
				if stopWebServer != nil {
					stopWebServer()
				}
			}()

			return []byte(fmt.Sprintf("pprof web started on port %s with duration of %s", cmd[3], sleepTime)), nil
		case "stop":
			if stopWebServer != nil {
				stopWebServer()
				return []byte("pprof web server stopping"), nil
			}
			return []byte(""), errors.New("pprof web server not running")
		default:
			return []byte{}, errors.New("unrecognized command")
		}
	case "cpu":
		switch cmd[2] {
		case "start":
			if len(cmd) < 4 {
				return []byte{}, errors.New("filename required: pprof cpu start filename")
			}
			sleepTime := time.Minute * 10
			if len(cmd) > 4 {
				sleepTime, err = time.ParseDuration(cmd[4])
				if err != nil {
					return []byte{}, err
				}
			}
			f, err := os.Create(cmd[3])
			if err != nil {
				return []byte{}, err
			}
			err = pprof.StartCPUProfile(f)
			if err != nil {
				return []byte{}, err
			}
			var ctx context.Context
			ctx, stopCPUProfiling = context.WithCancel(context.Background())

			go func(context.Context) {
				<-ctx.Done()
				mu.Lock()
				defer mu.Unlock()
				pprof.StopCPUProfile()
				stopCPUProfiling = nil
			}(ctx)

			go func() {
				time.Sleep(sleepTime)
				mu.Lock()
				defer mu.Unlock()
				select {
				case <-ctx.Done():
					return
				default:
				}
				if stopCPUProfiling != nil {
					stopCPUProfiling()
				}
			}()

			return []byte(fmt.Sprintf("CPU Profile with duration of %s", sleepTime)), nil
		case "stop":
			if stopCPUProfiling == nil {
				return []byte(""), errors.New("CPU Profile not running")
			}
			stopCPUProfiling()
			return []byte("CPU Profile stopped"), nil
		default:
			return []byte{}, errors.New("unrecognized command")
		}
	case "mem":
		if len(cmd) < 3 {
			return []byte{}, errors.New("filename required: pprof mem filename")
		}
		f, err := os.Create(cmd[2])
		if err != nil {
			return []byte{}, err
		}
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			return []byte{}, err
		}
		f.Close()
		return []byte("pprof.WriteHeapProfile executed"), nil
	case "examples":
		return []byte(pprofExamples), nil
	default:
		return []byte{}, errors.New("unrecognized command")
	}
}

const pprofExamples = `CPU profiling:
  echo "pprof cpu start /tmp/dataplane.cpu.pprof 1h42m" | socat -t 30 UNIX-CONNECT:dataplaneapi.sock -
  ...
  # make api requests
  ...
  echo "pprof cpu stop" | socat -t 30 UNIX-CONNECT:dataplaneapi.sock -
  # to view in browser, run
  go tool pprof -web dataplane.cpu.pprof

MEMORY profiling:
  echo "pprof mem /tmp/dataplane.mem.pprof" | socat -t 30 UNIX-CONNECT:dataplaneapi.sock -
  # to view in browser, run
  go tool pprof -web dataplane.mem.pprof

WEB:
  echo "pprof web start 8888 1h" | socat -t 30 UNIX-CONNECT:dataplaneapi.sock -
  # heap
  go tool pprof  http://localhost:6060/debug/pprof/heap
  # profile
  go tool pprof -web http://localhost:8888/debug/pprof/profile
  # trace
  curl -s http://127.0.0.1:8888/debug/pprof/trace?seconds=60 > cpu-trace.out
  go tool trace cpu-trace.out
  # pprof
  curl -s http://127.0.0.1:8888/debug/pprof/profile?seconds=60 > ./cpu.out
  ...
  # various API requests
  ...
  go tool pprof -web cpu.out
  echo "pprof web stop" | socat -t 30 UNIX-CONNECT:dataplaneapi.sock -
`
