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

package misc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/haproxytech/dataplaneapi/haproxy"

	"github.com/haproxytech/client-native/v2/configuration"
	client_errors "github.com/haproxytech/client-native/v2/errors"
	"github.com/haproxytech/models/v2"
)

const (
	// ErrHTTPNotFound HTTP status code 404
	ErrHTTPNotFound = int64(404)
	// ErrHTTPConflict HTTP status code 409
	ErrHTTPConflict = int64(409)
	// ErrHTTPInternalServerError HTTP status code 500
	ErrHTTPInternalServerError = int64(500)
	// ErrHTTPBadRequest HTTP status code 400
	ErrHTTPBadRequest = int64(400)
)

// HandleError translates error codes from client native into models.Error with appropriate http status code
func HandleError(err error) *models.Error {
	switch t := err.(type) {
	case *configuration.ConfError:
		msg := t.Error()
		httpCode := ErrHTTPInternalServerError
		switch t.Code() {
		case configuration.ErrObjectDoesNotExist:
			httpCode = ErrHTTPNotFound
		case configuration.ErrObjectAlreadyExists, configuration.ErrVersionMismatch, configuration.ErrTransactionAlreadyExists:
			httpCode = ErrHTTPConflict
		case configuration.ErrObjectIndexOutOfRange, configuration.ErrValidationError, configuration.ErrBothVersionTransaction,
			configuration.ErrNoVersionTransaction, configuration.ErrNoParentSpecified, configuration.ErrParentDoesNotExist,
			configuration.ErrTransactionDoesNotExist, configuration.ErrGeneralError:
			httpCode = ErrHTTPBadRequest
		}
		return &models.Error{Code: &httpCode, Message: &msg}
	case *haproxy.ReloadError:
		httpCode := ErrHTTPBadRequest
		msg := t.Error()
		return &models.Error{Code: &httpCode, Message: &msg}
	default:
		msg := t.Error()
		code := ErrHTTPInternalServerError
		return &models.Error{Code: &code, Message: &msg}
	}
}

// DiscoverChildPaths return children models.Endpoints given path
func DiscoverChildPaths(path string, spec json.RawMessage) (models.Endpoints, error) {
	var m map[string]interface{}
	err := json.Unmarshal(spec, &m)
	if err != nil {
		return nil, err
	}
	es := make(models.Endpoints, 0, 1)
	paths := m["paths"].(map[string]interface{})
	for key, value := range paths {
		v := value.(map[string]interface{})
		if g, ok := v["get"].(map[string]interface{}); ok {
			title := g["summary"].(string)
			description := g["description"].(string)

			if strings.HasPrefix(key, path) && key != path {
				if len(strings.Split(key[len(path)+1:], "/")) == 1 {
					e := models.Endpoint{
						URL:         key,
						Title:       title,
						Description: description,
					}
					es = append(es, &e)
				}
			}
		}
	}
	return es, nil
}

func IsUnixSocketAddr(addr string) bool {
	if strings.HasPrefix(addr, "ipv4@") || strings.HasPrefix(addr, "ipv6@") {
		return false
	}

	// check if it has semicolon
	if strings.Contains(addr, ":") {
		return false
	}
	return true
}

func ParseTimeout(tOut string) *int64 {
	var v int64
	switch {
	case strings.HasSuffix(tOut, "ms"):
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "ms"), 10, 64)
	case strings.HasSuffix(tOut, "s"):
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "s"), 10, 64)
		v *= 1000
	case strings.HasSuffix(tOut, "m"):
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "m"), 10, 64)
		v = v * 1000 * 60
	case strings.HasSuffix(tOut, "h"):
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "h"), 10, 64)
		v = v * 1000 * 60 * 60
	case strings.HasSuffix(tOut, "d"):
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "d"), 10, 64)
		v = v * 1000 * 60 * 60 * 24
	default:
		v, _ = strconv.ParseInt(tOut, 10, 64)
	}
	if v != 0 {
		return &v
	}
	return nil
}

func GetHTTPStatusFromErr(err error) int {
	switch {
	case errors.Is(err, client_errors.ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, client_errors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, client_errors.ErrGeneral):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func SetError(code int, msg string) *models.Error {
	return &models.Error{
		Code:    Int64P(code),
		Message: StringP(msg),
	}
}

func StringP(s string) *string {
	return &s
}

func Int64P(i int) *int64 {
	i64 := int64(i)
	return &i64
}

//extractEnvVar extracts and returns env variable from HAProxy variable
//provided in "${SOME_VAR}" format
func ExtractEnvVar(pass string) string {
	return strings.TrimLeft(strings.TrimRight(pass, "\"}"), "\"${")
}
