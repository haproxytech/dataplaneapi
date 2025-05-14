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
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server"
)

// GetRuntimeServerHandlerImpl implementation of the GetRuntimeServerHandler interface using client-native client
type GetRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetAllRuntimeServerHandlerImpl implementation of the GetRuntimeServersHandler interface using client-native client
type GetAllRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceRuntimeServerHandlerImpl implementation of the ReplaceRuntimeServerHandler interface using client-native client
type ReplaceRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// AddRuntimeServerHandlerImpl implementation of the ServerAddRuntimeServerHandler interface using client-native client
type AddRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// DeleteRuntimeServerHandlerImpl implementation of the ServerDeleteRuntimeServerHandler interface using client-native client
type DeleteRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *GetRuntimeServerHandlerImpl) Handle(params server.GetRuntimeServerParams, principal any) middleware.Responder {
	rn, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := rn.GetServerState(params.ParentName, params.Name)
	if err != nil {
		if isNotFoundError(err) {
			code := int64(404)
			msg := err.Error()
			return server.NewGetRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
		}
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if rs == nil {
		code := int64(404)
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.ParentName)
		return server.NewGetRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	return server.NewGetRuntimeServerOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *GetAllRuntimeServerHandlerImpl) Handle(params server.GetAllRuntimeServerParams, principal any) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetAllRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := runtime.GetServersState(params.ParentName)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server.NewGetAllRuntimeServerOK().WithPayload(models.RuntimeServers{})
		}
		return server.NewGetAllRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewGetAllRuntimeServerOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *ReplaceRuntimeServerHandlerImpl) Handle(params server.ReplaceRuntimeServerParams, principal any) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := runtime.GetServerState(params.ParentName, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if rs == nil {
		code := int64(404)
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.ParentName)
		return server.NewReplaceRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	// change operational state
	if params.Data.OperationalState != "" && rs.OperationalState != params.Data.OperationalState {
		err = runtime.SetServerHealth(params.ParentName, params.Name, params.Data.OperationalState)
		if err != nil {
			e := misc.HandleError(err)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	// change admin state
	if params.Data.AdminState != "" && rs.AdminState != params.Data.AdminState {
		err = runtime.SetServerState(params.ParentName, params.Name, params.Data.AdminState)
		if err != nil {
			e := misc.HandleError(err)

			// try to revert operational state and fall silently
			//nolint:errcheck
			runtime.SetServerHealth(params.ParentName, params.Name, rs.OperationalState)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	rs, err = runtime.GetServerState(params.ParentName, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewReplaceRuntimeServerOK().WithPayload(rs)
}

// Adds a new server dynamically without modifying the configuration.
// Warning: this only works if you have not defined a `default_server` in the defaults
// or in the current `backend` section.
func (h *AddRuntimeServerHandlerImpl) Handle(params server.AddRuntimeServerParams, principal any) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewAddRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if params.Data.Name == "" {
		code := int64(400)
		msg := "the new server must have a name"
		return server.NewAddRuntimeServerBadRequest().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	haversion, err := runtime.GetVersion()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewAddRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.AddServer(params.ParentName, params.Data.Name, SerializeRuntimeAddServer(params.Data, &haversion))
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "No such backend"):
			code := int64(404)
			return server.NewAddRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
		case strings.Contains(msg, "Already exists"):
			code := int64(409)
			return server.NewAddRuntimeServerConflict().WithPayload(&models.Error{Code: &code, Message: &msg})
		default:
			e := misc.HandleError(err)
			return server.NewAddRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	return server.NewAddRuntimeServerCreated().WithPayload(params.Data)
}

func isNotFoundError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "No such backend") ||
		strings.Contains(msg, native_errors.ErrNotFound.Error())
}

// Deletes a server from a backend immediately, without waiting for connections to drain.
func (h *DeleteRuntimeServerHandlerImpl) Handle(params server.DeleteRuntimeServerParams, principal any) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	// Check if this server exists.
	rs, err := runtime.GetServerState(params.ParentName, params.Name)
	if err != nil {
		if isNotFoundError(err) {
			code := int64(404)
			msg := err.Error()
			return server.NewDeleteRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
		}
		e := misc.HandleError(err)
		return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	// Put the server in maintenance state before deleting it.
	if rs.AdminState != "maint" {
		err = runtime.DisableServer(params.ParentName, params.Name)
		if err != nil {
			e := misc.HandleError(err)
			return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	// TODO: wait for connections to drain. This is not yet possible with HAProxy 2.6.

	err = runtime.DeleteServer(params.ParentName, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewDeleteRuntimeServerNoContent()
}

// SerializeRuntimeAddServer returns a string in the HAProxy config format, suitable
// for the "add server" operation over the control socket.
// Not all the Server attributes are available in this case.
func SerializeRuntimeAddServer(srv *models.RuntimeAddServer, version *cn_runtime.HAProxyVersion) string { //nolint: cyclop,maintidx
	b := &strings.Builder{}

	push := func(s string) {
		b.WriteByte(' ')
		b.WriteString(s)
	}
	pushi := func(key string, val *int64) {
		fmt.Fprintf(b, " %s %d", key, *val)
	}
	// push a quoted string
	pushq := func(key, val string) {
		fmt.Fprintf(b, ` %s %s`, key, val)
	}
	enabled := func(s string) bool {
		return s == "enabled"
	}
	disabled := func(s string) bool {
		return s == "disabled"
	}

	// Address is mandatory and must come first, with an optional port number.
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
