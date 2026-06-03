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

package info

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/sys/unix"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all info routes onto r using spec-based request validation
// and a shared error handler.
func RegisterRouter(r chi.Router, version, buildTime string, systemInfo bool) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{
		Version:    version,
		BuildTime:  buildTime,
		SystemInfo: systemInfo,
	}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for API info.
type HandlerImpl struct {
	Version    string
	BuildTime  string
	SystemInfo bool
}

func (h *HandlerImpl) GetInfo(w http.ResponseWriter, r *http.Request) {
	api := &models.InfoAPI{
		Version: h.Version,
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", h.BuildTime)
	if err == nil {
		api.BuildDate = strfmt.DateTime(date)
	} else {
		fmt.Println(err.Error())
	}

	sys := &models.InfoSystem{}

	if h.SystemInfo {
		hName, err := os.Hostname()
		if err == nil {
			sys.Hostname = hName
		}

		sys.MemInfo = &models.InfoSystemMemInfo{}
		sys.CPUInfo = &models.InfoSystemCPUInfo{}

		memInfo, err := mem.VirtualMemory()
		if err == nil {
			sys.MemInfo.TotalMemory = int64(memInfo.Total)
			sys.MemInfo.FreeMemory = int64(memInfo.Free)
		}
		if uptime, err := host.Uptime(); err == nil {
			uptimeInt64 := int64(uptime)
			sys.Uptime = &uptimeInt64
		}

		sys.CPUInfo.NumCpus = int64(runtime.NumCPU())

		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)
		sys.MemInfo.DataplaneapiMemory = int64(m.Sys)

		sys.CPUInfo.Model = parseCPUModel()

		uName := &unix.Utsname{}
		err = unix.Uname(uName)
		if err == nil {
			sys.OsString = string(bytes.Trim(uName.Sysname[:], "\x00")) + " " + string(bytes.Trim(uName.Release[:], "\x00")) + " " + string(bytes.Trim(uName.Version[:], "\x00"))
		}
		sys.Time = time.Now().Unix()
	}

	respond.JSON(w, http.StatusOK, &models.Info{API: api, System: sys})
}

func parseCPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "model name") {
			s := strings.Split(l, ":")
			return strings.TrimSpace(strings.Join(s[1:], ":"))
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}
