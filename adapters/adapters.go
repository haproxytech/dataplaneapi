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
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/haproxytech/client-native/v2/models"
	apachelog "github.com/lestrrat-go/apache-logformat"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

var (
	apacheLogWritter *io.PipeWriter
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
func RecoverMiddleware(entry logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					frames := callers()

					entry := entry.WithField("stack_trace", frames.String())
					entry.Error(fmt.Sprintf("Panic %v", err))

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

func UniqueIDMiddleware(entry *logrus.Entry) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-Id")
			if len(reqID) == 0 {
				t := time.Now()
				// we need performance and math/rand does the trick, although it's "unsafe":
				// speed is absolutely required here since we're going to generate an ULID for each request
				// that doesn't contain a prefilled X-Request-Id header.
				entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0) // nolint:gosec
				reqID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
			}
			*entry = *entry.WithField("request_id", reqID)

			h.ServeHTTP(w, r)
		})
	}
}

func ApacheLogMiddleware(entry *logrus.Entry, format *apachelog.ApacheLog) Adapter {
	return func(h http.Handler) http.Handler {
		if apacheLogWritter != nil {
			apacheLogWritter.Close()
		}
		apacheLogWritter = entry.Logger.Writer()
		return format.Wrap(h, apacheLogWritter)
	}
}

// LoggingMiddleware logs before and after the response to given logger
func LoggingMiddleware(entry *logrus.Entry) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logBefore(entry, r)
			res := newStatusResponseWriter(w)
			defer logAfter(entry, res, start)
			h.ServeHTTP(res, r)
		})
	}
}

func logBefore(entry *logrus.Entry, req *http.Request) {
	entry.WithFields(logrus.Fields{
		"request": req.RequestURI,
		"method":  req.Method,
		"remote": func() (remote string) {
			remote = req.RemoteAddr
			if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
				remote = realIP
			}
			return
		}(),
	}).Info("started handling request")
}

func logAfter(entry *logrus.Entry, res *statusResponseWriter, start time.Time) {
	entry.WithFields(logrus.Fields{
		"status": res.Status(),
		"length": units.HumanSize(float64(res.Length())),
		"took":   time.Since(start),
	}).Info("completed handling request")
}
