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
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v3/configuration"
	client_errors "github.com/haproxytech/client-native/v3/errors"
	"github.com/haproxytech/client-native/v3/models"
	jsoniter "github.com/json-iterator/go"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/rate"
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
	// ErrHTTPRateLimit HTTP status code 429
	ErrHTTPRateLimit = int64(429)
	// ErrHTTPOk HTTP status code 200
	ErrHTTPOk = int64(200)
)

func OutdatedTransactionError(id string) *models.Error {
	var code int64 = 406
	msg := fmt.Sprintf("transaction %s is outdated and cannot be committed", id)
	return &models.Error{
		Code:    &code,
		Message: &msg,
	}
}

func FailedTransactionError(id string) *models.Error {
	var code int64 = 406
	msg := fmt.Sprintf("transaction %s is failed and cannot be committed", id)
	return &models.Error{
		Code:    &code,
		Message: &msg,
	}
}

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
	case *rate.TransactionLimitReachedErr:
		httpCode := ErrHTTPRateLimit
		msg := t.Error()
		return &models.Error{Code: &httpCode, Message: &msg}
	default:
		msg := t.Error()
		code := ErrHTTPInternalServerError
		return &models.Error{Code: &code, Message: &msg}
	}
}

// HandleContainerGetError translates error codes from client native into models.Error with appropriate http status code. Intended for get requests on container endpoints.
func HandleContainerGetError(err error) *models.Error {
	if t, ok := err.(*configuration.ConfError); ok {
		if t.Code() == configuration.ErrParentDoesNotExist {
			code := ErrHTTPOk
			return &models.Error{Code: &code}
		}
	}
	return HandleError(err)
}

// DiscoverChildPaths return children models.Endpoints given path
func DiscoverChildPaths(path string, spec json.RawMessage) (models.Endpoints, error) {
	var m map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(spec, &m)
	if err != nil {
		return nil, err
	}
	es := make(models.Endpoints, 0, 1)
	paths := m["paths"].(map[string]interface{})
	for key, value := range paths {
		v := value.(map[string]interface{})
		if g, ok := v["get"].(map[string]interface{}); ok {
			title := ""
			if titleInterface, ok := g["summary"]; ok && titleInterface != nil {
				title = titleInterface.(string)
			}
			description := ""
			if descInterface, ok := g["description"]; ok && descInterface != nil {
				description = descInterface.(string)
			}

			if strings.HasPrefix(key, path) && key != path {
				resource := key[len(path):]
				if strings.HasPrefix(resource, "/") && len(strings.Split(resource[1:], "/")) == 1 {
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

// extractEnvVar extracts and returns env variable from HAProxy variable
// provided in "${SOME_VAR}" format
func ExtractEnvVar(pass string) string {
	return strings.TrimLeft(strings.TrimRight(pass, "\"}"), "\"${")
}

func HasOSArg(short, long, env string) bool {
	if short == "" && long == "" && env == "" {
		return false
	}
	target1 := "--" + long
	hasShort := short != ""
	target2 := "-" + short

	if env != "" {
		if os.Getenv(env) != "" {
			return true
		}
	}
	for _, arg := range os.Args {
		if hasShort && arg == target2 {
			return true
		}
		if arg == target1 {
			return true
		}
		if strings.HasPrefix(arg, target1) {
			p := strings.Split(arg, "=")
			if len(p) > 1 {
				return true
			}
		}
	}
	return false
}

func RandomString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	result := strings.ReplaceAll(base64.URLEncoding.EncodeToString(b), `=`, ``)
	result = strings.ReplaceAll(result, `-`, ``)
	result = strings.ReplaceAll(result, `_`, ``)
	return result[:size], err
}
