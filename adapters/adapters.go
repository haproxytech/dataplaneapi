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
	"strings"
	"time"

	"github.com/docker/go-units"

	"github.com/haproxytech/models"

	"github.com/sirupsen/logrus"
)

// Adapter is just a wrapper over http handler function
type Adapter func(http.Handler) http.Handler

type statusResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (srw *statusResponseWriter) WriteHeader(s int) {
	srw.status = s
	srw.ResponseWriter.WriteHeader(s)
}

func (srw *statusResponseWriter) Write(b []byte) (int, error) {
	if srw.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		srw.status = 200
		srw.ResponseWriter.WriteHeader(srw.Status())
	}
	size, err := srw.ResponseWriter.Write(b)
	srw.length = size
	return size, err
}

func (srw *statusResponseWriter) Header() http.Header {
	return srw.ResponseWriter.Header()
}

func (srw *statusResponseWriter) Status() int {
	return srw.status
}

func (srw *statusResponseWriter) Length() int {
	return srw.length
}

func newStatusResponseWriter(rw http.ResponseWriter) *statusResponseWriter {
	nsrw := &statusResponseWriter{
		ResponseWriter: rw,
	}
	return nsrw
}

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
func RecoverMiddleware(logger *logrus.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					frames := callers()

					entry := logrus.NewEntry(logger)
					entry = entry.WithField("stack_trace", frames.String())
					entry.Error(fmt.Sprintf("Panic %v", err))

					w.WriteHeader(http.StatusInternalServerError)

					code := int64(http.StatusInternalServerError)
					msg := fmt.Sprintf("%v: %v", err, frames[0].String())
					e := &models.Error{
						Code:    &code,
						Message: &msg,
					}

					errMsg, _ := e.MarshalJSON()
					ct := r.Header.Get(http.CanonicalHeaderKey("Content-Type"))
					if strings.HasPrefix(ct, "application/json") {
						w.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json")
					}
					w.Write(errMsg)
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// LoggingMiddleware logs before and after the response to given logger
func LoggingMiddleware(logger *logrus.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logBefore(logger, r)
			res := newStatusResponseWriter(w)
			defer logAfter(logger, res, start)
			h.ServeHTTP(res, r)
		})
	}
}

func logBefore(logger *logrus.Logger, req *http.Request) {
	e := logrus.NewEntry(logger)
	if reqID := req.Header.Get("X-Request-Id"); reqID != "" {
		e = e.WithField("request_id", reqID)
	}
	e = e.WithField("request", req.RequestURI)
	e = e.WithField("method", req.Method)
	remote := req.RemoteAddr
	if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
		remote = realIP
	}
	e = e.WithField("remote", remote)

	e.Info("started handling request")
}

func logAfter(logger *logrus.Logger, res *statusResponseWriter, start time.Time) {
	latency := time.Since(start)
	e := logrus.NewEntry(logger)
	e = e.WithField("status", res.Status())
	e = e.WithField("length", units.HumanSize(float64(res.Length())))
	e = e.WithField("took", latency)
	e.Info("completed handling request")
}
