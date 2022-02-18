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

package adapters

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	clientnative "github.com/haproxytech/client-native/v3"
	"github.com/haproxytech/client-native/v3/models"
	"github.com/haproxytech/dataplaneapi/log"
)

var configVersion string

func ConfigVersion() string {
	return configVersion
}

// Adapter is just a wrapper over http handler function
type Adapter func(http.Handler) http.Handler

type frame struct {
	File string
	Line int
	Name string
}

type stack []frame

func (f frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Name)
}

func (s stack) String() string {
	str := ""
	for _, f := range s {
		str = fmt.Sprintf("%v\n%v", str, f.String())
	}
	return str
}

func callers() stack {
	pcs := make([]uintptr, 32)
	num := runtime.Callers(5, pcs)
	st := make(stack, num)
	for i, pc := range pcs[:num] {
		fun := runtime.FuncForPC(pc)
		file, line := fun.FileLine(pc - 1)
		st[i].File = file
		st[i].Line = line
		st[i].Name = stripPackage(fun.Name())
	}
	return st
}

func stripPackage(n string) string {
	slashI := strings.LastIndex(n, "/")
	if slashI == -1 {
		slashI = 0
	}
	dotI := strings.Index(n[slashI:], ".")
	if dotI == -1 {
		return n
	}
	return n[slashI+dotI+1:]
}

// RecoverMiddleware used for recovering from panic, logs the panic to given logger and returns 500
func RecoverMiddleware(logger *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					frames := callers()

					fields := make(map[string]interface{})
					fields["stack_trace"] = frames.String()
					logger.WithFieldsf(fields, log.ErrorLevel, "Panic %v", err)

					w.WriteHeader(http.StatusInternalServerError)

					code := int64(http.StatusInternalServerError)
					msg := fmt.Sprintf("%v: %v", err, frames[0].String())
					e := &models.Error{
						Code:    &code,
						Message: &msg,
					}

					errMsg, _ := e.MarshalJSON()
					ct := r.Header.Get("Content-Type")
					if strings.HasPrefix(ct, "application/json") {
						w.Header().Set("Content-Type", "application/json")
					}
					// nolint:errcheck
					w.Write(errMsg)
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func ApacheLogMiddleware(logger *log.ACLLogger) Adapter {
	return func(h http.Handler) http.Handler {
		apacheLogWritter := logger.Writer()
		return logger.ApacheLog.Wrap(h, apacheLogWritter)
	}
}

func ConfigVersionMiddleware(client clientnative.HAProxyClient) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			qs := r.URL.Query()
			tID := qs.Get("transaction_id")
			configuration, err := client.Configuration()
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotImplemented)
			}
			tr, _ := configuration.GetTransaction(tID)
			var v int64
			if tr != nil && tr.Status == models.TransactionStatusInProgress {
				v, err = configuration.GetConfigurationVersion(tr.ID)
			} else {
				v, err = configuration.GetConfigurationVersion("")
			}
			if err == nil {
				configVersion = strconv.FormatInt(v, 10)
				w.Header().Add("Configuration-Version", configVersion)
			}
			h.ServeHTTP(w, r)
		})
	}
}
