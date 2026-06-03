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

package servers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all servers routes onto r using spec-based request validation.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy runtime servers.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllRuntimeServer(w http.ResponseWriter, r *http.Request, parentName string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	rs, err := rt.GetServersState(parentName)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, models.RuntimeServers{})
			return
		}
		respond.JSON(w, int(*e.Code), e)
		return
	}
	respond.JSON(w, http.StatusOK, rs)
}

func (h *HandlerImpl) AddRuntimeServer(w http.ResponseWriter, r *http.Request, parentName string) {
	var data AddRuntimeServerJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	if data.Name == "" {
		respond.BadRequest(w, "the new server must have a name")
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	haversion, err := rt.GetVersion()
	if err != nil {
		respond.Error(w, err)
		return
	}

	err = rt.AddServer(parentName, data.Name, SerializeRuntimeAddServer(&data, &haversion))
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "No such backend"):
			code := int64(404)
			respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		case strings.Contains(msg, "Already exists"):
			code := int64(409)
			respond.JSON(w, http.StatusConflict, &models.Error{Code: &code, Message: &msg})
		default:
			respond.Error(w, err)
		}
		return
	}

	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteRuntimeServer(w http.ResponseWriter, r *http.Request, parentName string, name string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	rs, err := rt.GetServerState(parentName, name)
	if err != nil {
		if isNotFoundError(err) {
			msg := err.Error()
			code := int64(http.StatusNotFound)
			respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
			return
		}
		respond.Error(w, err)
		return
	}

	if rs.AdminState != "maint" {
		if err = rt.DisableServer(parentName, name); err != nil {
			respond.Error(w, err)
			return
		}
	}

	if err = rt.DeleteServer(parentName, name); err != nil {
		respond.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) GetRuntimeServer(w http.ResponseWriter, r *http.Request, parentName string, name string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	rs, err := rt.GetServerState(parentName, name)
	if err != nil {
		if isNotFoundError(err) {
			code := int64(404)
			msg := err.Error()
			respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
			return
		}
		respond.Error(w, err)
		return
	}

	if rs == nil {
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", name, parentName)
		code := int64(404)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}

	respond.JSON(w, http.StatusOK, rs)
}

func isNotFoundError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "No such backend") ||
		strings.Contains(msg, native_errors.ErrNotFound.Error())
}

func (h *HandlerImpl) ReplaceRuntimeServer(w http.ResponseWriter, r *http.Request, parentName string, name string) {
	var data ReplaceRuntimeServerJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	rs, err := rt.GetServerState(parentName, name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if rs == nil {
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", name, parentName)
		code := int64(404)
		respond.JSON(w, http.StatusNotFound, &models.Error{Code: &code, Message: &msg})
		return
	}

	// save original values for rollback
	origOperationalState := rs.OperationalState
	origAdminState := rs.AdminState
	origWeight := rs.Weight

	var changedOperational, changedAdmin, changedWeight bool

	// change operational state
	if data.OperationalState != "" && rs.OperationalState != data.OperationalState {
		if err = rt.SetServerHealth(parentName, name, data.OperationalState); err != nil {
			respond.Error(w, err)
			return
		}
		changedOperational = true
	}

	// change admin state
	if data.AdminState != "" && rs.AdminState != data.AdminState {
		if err = rt.SetServerState(parentName, name, data.AdminState); err != nil {
			if changedOperational {
				//nolint:errcheck
				rt.SetServerHealth(parentName, name, origOperationalState)
			}
			respond.Error(w, err)
			return
		}
		changedAdmin = true
	}

	// change weight
	if data.Weight != nil && (rs.Weight == nil || *data.Weight != *rs.Weight) {
		if err = rt.SetServerWeight(parentName, name, strconv.FormatInt(*data.Weight, 10)); err != nil {
			if changedAdmin {
				//nolint:errcheck
				rt.SetServerState(parentName, name, origAdminState)
			}
			if changedOperational {
				//nolint:errcheck
				rt.SetServerHealth(parentName, name, origOperationalState)
			}
			respond.Error(w, err)
			return
		}
		changedWeight = true
	}

	// change address/port
	addrChanged := data.Address != "" && rs.Address != data.Address
	portChanged := data.Port != nil && (rs.Port == nil || *data.Port != *rs.Port)
	if addrChanged || portChanged {
		newAddr := rs.Address
		if data.Address != "" {
			newAddr = data.Address
		}
		var newPort int
		if data.Port != nil {
			newPort = int(*data.Port)
		} else if rs.Port != nil {
			newPort = int(*rs.Port)
		}
		if err = rt.SetServerAddr(parentName, name, newAddr, newPort); err != nil {
			if changedWeight {
				//nolint:errcheck
				rt.SetServerWeight(parentName, name, formatWeightPtr(origWeight))
			}
			if changedAdmin {
				//nolint:errcheck
				rt.SetServerState(parentName, name, origAdminState)
			}
			if changedOperational {
				//nolint:errcheck
				rt.SetServerHealth(parentName, name, origOperationalState)
			}
			respond.Error(w, err)
			return
		}
	}

	rs, err = rt.GetServerState(parentName, name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, rs)
}

func formatWeightPtr(w *int64) string {
	if w == nil {
		return "0"
	}
	return strconv.FormatInt(*w, 10)
}

// SerializeRuntimeAddServer returns a string in the HAProxy config format, suitable
// for the "add server" operation over the control socket.
func SerializeRuntimeAddServer(srv *models.RuntimeAddServer, version *cn_runtime.HAProxyVersion) string { //nolint: cyclop,maintidx
	b := &strings.Builder{}

	push := func(s string) {
		b.WriteByte(' ')
		b.WriteString(s)
	}
	pushi := func(key string, val *int64) {
		fmt.Fprintf(b, " %s %d", key, *val)
	}
	pushq := func(key, val string) {
		fmt.Fprintf(b, ` %s %s`, key, val)
	}
	enabled := func(s string) bool {
		return s == "enabled"
	}
	disabled := func(s string) bool {
		return s == "disabled"
	}

	addr := srv.Address
	if srv.Port != nil {
		addr += fmt.Sprintf(":%d", *srv.Port)
	}
	push(addr)

	if enabled(srv.AgentCheck) {
		push("agent-check")
	}
	if srv.AgentAddr != "" {
		pushq("agent-addr", srv.AgentAddr)
	}
	if srv.AgentPort != nil {
		pushi("agent-port", srv.AgentPort)
	}
	if srv.AgentInter != nil {
		pushi("agent-inter", srv.AgentInter)
	}
	if srv.AgentSend != "" {
		pushq("agent-send", srv.AgentSend)
	}
	if srv.Allow0rtt {
		push("allow-0rtt")
	}
	if srv.Alpn != "" {
		pushq("alpn", srv.Alpn)
	}
	if enabled(srv.Backup) {
		push("backup")
	}
	if srv.SslCafile != "" {
		pushq("ca-file", srv.SslCafile)
	}
	if enabled(srv.Check) {
		push("check")
	}
	if srv.CheckAlpn != "" {
		pushq("check-alpn", srv.CheckAlpn)
	}
	if srv.HealthCheckAddress != "" {
		pushq("addr", srv.HealthCheckAddress)
	}
	if srv.HealthCheckPort != nil {
		pushi("port", srv.HealthCheckPort)
	}
	if srv.CheckProto != "" {
		pushq("check-proto", srv.CheckProto)
	}
	if enabled(srv.CheckSendProxy) {
		push("check-send-proxy")
	}
	if srv.CheckSni != "" {
		pushq("check-sni", srv.CheckSni)
	}
	if enabled(srv.CheckSsl) {
		push("check-ssl")
	}
	if enabled(srv.CheckViaSocks4) {
		push("check-via-socks4")
	}
	if srv.Ciphers != "" {
		pushq("ciphers", srv.Ciphers)
	}
	if srv.Ciphersuites != "" {
		pushq("ciphersuites", srv.Ciphersuites)
	}
	if srv.CrlFile != "" {
		pushq("crl-file", srv.CrlFile)
	}
	if srv.SslCertificate != "" {
		pushq("crt", srv.SslCertificate)
	}
	if enabled(srv.Maintenance) {
		push("disabled")
	}
	if srv.Downinter != nil {
		pushi("downinter", srv.Downinter)
	}
	if disabled(srv.Maintenance) {
		required := new(cn_runtime.HAProxyVersion)
		required.ParseHAProxyVersion("3.0.0")
		if !cn_runtime.IsBiggerOrEqual(required, version) {
			push("enabled")
		}
	}
	if srv.ErrorLimit != nil {
		pushi("error-limit", srv.ErrorLimit)
	}
	if srv.Fall != nil {
		pushi("fall", srv.Fall)
	}
	if srv.Fastinter != nil {
		pushi("fastinter", srv.Fastinter)
	}
	if enabled(srv.ForceSslv3) {
		push("force-sslv3")
	}
	if enabled(srv.ForceTlsv10) {
		push("force-tlsv10")
	}
	if enabled(srv.ForceTlsv11) {
		push("force-tlsv11")
	}
	if enabled(srv.ForceTlsv12) {
		push("force-tlsv12")
	}
	if enabled(srv.ForceTlsv13) {
		push("force-tlsv13")
	}
	if srv.ID != "" {
		pushq("id", srv.ID)
	}
	if srv.Inter != nil {
		pushi("inter", srv.Inter)
	}
	if srv.Maxconn != nil {
		pushi("maxconn", srv.Maxconn)
	}
	if srv.Maxqueue != nil {
		pushi("maxqueue", srv.Maxqueue)
	}
	if srv.Minconn != nil {
		pushi("minconn", srv.Minconn)
	}
	if disabled(srv.SslReuse) {
		push("no-ssl-reuse")
	}
	if enabled(srv.NoSslv3) {
		push("no-sslv3")
	}
	if enabled(srv.NoTlsv10) {
		push("no-tlsv10")
	}
	if enabled(srv.NoTlsv11) {
		push("no-tlsv11")
	}
	if enabled(srv.NoTlsv12) {
		push("no-tlsv12")
	}
	if enabled(srv.NoTlsv13) {
		push("no-tlsv13")
	}
	if disabled(srv.TLSTickets) {
		push("no-tls-tickets")
	}
	if srv.Npn != "" {
		pushq("npm", srv.Npn)
	}
	if srv.Observe != "" {
		pushq("observe", srv.Observe)
	}
	if srv.OnError != "" {
		pushq("on-error", srv.OnError)
	}
	if srv.OnMarkedDown != "" {
		pushq("on-marked-down", srv.OnMarkedDown)
	}
	if srv.OnMarkedUp != "" {
		pushq("on-marked-up", srv.OnMarkedUp)
	}
	if srv.PoolLowConn != nil {
		pushi("pool-low-conn", srv.PoolLowConn)
	}
	if srv.PoolMaxConn != nil {
		pushi("pool-max-conn", srv.PoolMaxConn)
	}
	if srv.PoolPurgeDelay != nil {
		pushi("pool-purge-delay", srv.PoolPurgeDelay)
	}
	if srv.Proto != "" {
		pushq("proto", srv.Proto)
	}
	if len(srv.ProxyV2Options) > 0 {
		pushq("proxy-v2-options", strings.Join(srv.ProxyV2Options, ","))
	}
	if srv.Rise != nil {
		pushi("rise", srv.Rise)
	}
	if enabled(srv.SendProxy) {
		push("send-proxy")
	}
	if enabled(srv.SendProxyV2) {
		push("send-proxy-v2")
	}
	if enabled(srv.SendProxyV2Ssl) {
		push("send-proxy-v2-ssl")
	}
	if enabled(srv.SendProxyV2SslCn) {
		push("send-proxy-v2-ssl-cn")
	}
	if srv.Slowstart != nil {
		pushi("slowstart", srv.Slowstart)
	}
	if srv.Sni != "" {
		pushq("sni", srv.Sni)
	}
	if srv.Source != "" {
		pushq("source", srv.Source)
	}
	if enabled(srv.Ssl) {
		push("ssl")
	}
	if srv.SslMaxVer != "" {
		pushq("ssl-max-ver", srv.SslMaxVer)
	}
	if srv.SslMinVer != "" {
		pushq("ssl-min-ver", srv.SslMinVer)
	}
	if enabled(srv.Tfo) {
		push("tfo")
	}
	if enabled(srv.TLSTickets) {
		push("tls-tickets")
	}
	if srv.Track != "" {
		pushq("track", srv.Track)
	}
	if srv.Verify != "" {
		pushq("verify", srv.Verify)
	}
	if srv.Verifyhost != "" {
		pushq("verifyhost", srv.Verifyhost)
	}
	if srv.Weight != nil {
		pushi("weight", srv.Weight)
	}
	if srv.Ws != "" {
		pushq("ws", srv.Ws)
	}
	return b.String()
}
