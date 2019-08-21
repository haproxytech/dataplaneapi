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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/information"
	"github.com/haproxytech/models"
)

//GetHaproxyProcessInfoHandlerImpl implementation of the GetHaproxyProcessInfoHandler interface using client-native client
type GetHaproxyProcessInfoHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetHaproxyProcessInfoHandlerImpl) Handle(params information.GetHaproxyProcessInfoParams, principal interface{}) middleware.Responder {
	info, err := h.Client.Runtime.GetInfo()
	if err != nil || len(info) == 0 {
		code := misc.ErrHTTPInternalServerError
		msg := err.Error()
		e := &models.Error{
			Code:    &code,
			Message: &msg,
		}
		return information.NewGetHaproxyProcessInfoDefault(int(misc.ErrHTTPInternalServerError)).WithPayload(e)
	}

	data := models.ProcessInfo{}
	data.Haproxy = &info[0]

	return information.NewGetHaproxyProcessInfoOK().WithPayload(&data)
}

//GetInfoHandlerImpl implementation of the GetInfoHandler interface
type GetInfoHandlerImpl struct {
	SystemInfo bool
	BuildTime  string
	Version    string
}

//Handle executing the request and returning a response
func (h *GetInfoHandlerImpl) Handle(params information.GetInfoParams, principal interface{}) middleware.Responder {
	api := &models.InfoAPI{
		Version: h.Version,
	}
	date, err := time.Parse("2006-01-02T15:04:05", h.BuildTime)
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

		sysInfo := &unix.Sysinfo_t{}
		err = unix.Sysinfo(sysInfo)

		sys.MemInfo = &models.InfoSystemMemInfo{}
		sys.CPUInfo = &models.InfoSystemCPUInfo{}

		if err == nil {
			sys.Uptime = &sysInfo.Uptime
			sys.MemInfo.TotalMemory = int64(sysInfo.Totalram * uint64(sysInfo.Unit))
			sys.MemInfo.FreeMemory = int64(sysInfo.Freeram * uint64(sysInfo.Unit))
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

	return information.NewGetInfoOK().WithPayload(&models.Info{API: api, System: sys})
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

func cutTrailingBytes(arr []byte) []byte {
	r := make([]byte, len(arr))
	for _, b := range arr {
		if b == 0 {
			continue
		}
		r = append(r, b)
	}
	return r
}
