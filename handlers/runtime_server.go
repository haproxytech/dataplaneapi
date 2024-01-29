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

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server"
)

// GetRuntimeServerHandlerImpl implementation of the GetRuntimeServerHandler interface using client-native client
type GetRuntimeServerHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetRuntimeServersHandlerImpl implementation of the GetRuntimeServersHandler interface using client-native client
type GetRuntimeServersHandlerImpl struct {
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
func (h *GetRuntimeServerHandlerImpl) Handle(params server.GetRuntimeServerParams, principal interface{}) middleware.Responder {
	rn, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := rn.GetServerState(params.Backend, params.Name)
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
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.Backend)
		return server.NewGetRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	return server.NewGetRuntimeServerOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *GetRuntimeServersHandlerImpl) Handle(params server.GetRuntimeServersParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetRuntimeServersDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := runtime.GetServersState(params.Backend)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server.NewGetRuntimeServersOK().WithPayload(models.RuntimeServers{})
		}
		return server.NewGetRuntimeServersDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewGetRuntimeServersOK().WithPayload(rs)
}

// Handle executing the request and returning a response
func (h *ReplaceRuntimeServerHandlerImpl) Handle(params server.ReplaceRuntimeServerParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	rs, err := runtime.GetServerState(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	if rs == nil {
		code := int64(404)
		msg := fmt.Sprintf("Runtime server %s not found in backend %s", params.Name, params.Backend)
		return server.NewReplaceRuntimeServerNotFound().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	// change operational state
	if params.Data.OperationalState != "" && rs.OperationalState != params.Data.OperationalState {
		err = runtime.SetServerHealth(params.Backend, params.Name, params.Data.OperationalState)
		if err != nil {
			e := misc.HandleError(err)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	// change admin state
	if params.Data.AdminState != "" && rs.AdminState != params.Data.AdminState {
		err = runtime.SetServerState(params.Backend, params.Name, params.Data.AdminState)
		if err != nil {
			e := misc.HandleError(err)

			// try to revert operational state and fall silently
			//nolint:errcheck
			runtime.SetServerHealth(params.Backend, params.Name, rs.OperationalState)
			return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	rs, err = runtime.GetServerState(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewReplaceRuntimeServerOK().WithPayload(rs)
}

// Adds a new server dynamically without modifying the configuration.
// Warning: this only works if you have not defined a `default_server` in the defaults
// or in the current `backend` section.
func (h *AddRuntimeServerHandlerImpl) Handle(params server.AddRuntimeServerParams, principal interface{}) middleware.Responder {
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

	err = runtime.AddServer(params.Backend, params.Data.Name, SerializeRuntimeAddServer(params.Data))
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
func (h *DeleteRuntimeServerHandlerImpl) Handle(params server.DeleteRuntimeServerParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	// Check if this server exists.
	rs, err := runtime.GetServerState(params.Backend, params.Name)
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
		err = runtime.DisableServer(params.Backend, params.Name)
		if err != nil {
			e := misc.HandleError(err)
			return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
		}
	}

	// TODO: wait for connections to drain. This is not yet possible with HAProxy 2.6.

	err = runtime.DeleteServer(params.Backend, params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteRuntimeServerDefault(int(*e.Code)).WithPayload(e)
	}

	return server.NewDeleteRuntimeServerNoContent()
}

// SerializeRuntimeAddServer returns a string in the HAProxy config format, suitable
// for the "add server" operation over the control socket.
// Not all the Server attributes are available in this case.
func SerializeRuntimeAddServer(srv *models.RuntimeAddServer) string { //nolint:cyclop,maintidx
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
		fmt.Fprintf(b, ` %s "%s"`, key, val)
	}
	enabled := func(s string) bool {
		return s == "enabled"
	}

	// Address is mandatory and must come first, with an optional port number.
	addr := srv.Address
	if srv.Port != nil {
		addr += fmt.Sprintf(":%d", *srv.Port)
	}
	push(addr)

	switch {
	case enabled(srv.AgentCheck):
		push("agent-check")
	case srv.AgentAddr != "":
		pushq("agent-addr", srv.AgentAddr)
	case srv.AgentPort != nil:
		pushi("agent-port", srv.AgentPort)
	case srv.AgentInter != nil:
		pushi("agent-inter", srv.AgentInter)
	case srv.AgentSend != "":
		pushq("agent-send", srv.AgentSend)
	case srv.Allow0rtt:
		push("allow-0rtt")
	case srv.Alpn != "":
		pushq("alpn", srv.Alpn)
	case enabled(srv.Backup):
		push("backup")
	case srv.SslCafile != "":
		pushq("ca-file", srv.SslCafile)
	case enabled(srv.Check):
		push("check")
	case srv.CheckAlpn != "":
		pushq("check-alpn", srv.CheckAlpn)
	case srv.HealthCheckAddress != "":
		pushq("addr", srv.HealthCheckAddress)
	case srv.HealthCheckPort != nil:
		pushi("port", srv.HealthCheckPort)
	case srv.CheckProto != "":
		pushq("check-proto", srv.CheckProto)
	case enabled(srv.CheckSendProxy):
		push("check-send-proxy")
	case srv.CheckSni != "":
		pushq("check-sni", srv.CheckSni)
	case enabled(srv.CheckSsl):
		push("check-ssl")
	case enabled(srv.CheckViaSocks4):
		push("check-via-socks4")
	case srv.Ciphers != "":
		pushq("ciphers", srv.Ciphers)
	case srv.Ciphersuites != "":
		pushq("ciphersuites", srv.Ciphersuites)
	case srv.CrlFile != "":
		pushq("crl-file", srv.CrlFile)
	case srv.SslCertificate != "":
		pushq("crt", srv.SslCertificate)
	case enabled(srv.Maintenance):
		push("disabled")
	case srv.Downinter != nil:
		pushi("downinter", srv.Downinter)
	case !enabled(srv.Maintenance):
		push("enabled")
	case srv.ErrorLimit != nil:
		pushi("error-limit", srv.ErrorLimit)
	case srv.Fall != nil:
		pushi("fall", srv.Fall)
	case srv.Fastinter != nil:
		pushi("fastinter", srv.Fastinter)
	case enabled(srv.ForceSslv3):
		push("force-sslv3")
	case enabled(srv.ForceTlsv10):
		push("force-tlsv10")
	case enabled(srv.ForceTlsv11):
		push("force-tlsv11")
	case enabled(srv.ForceTlsv12):
		push("force-tlsv12")
	case enabled(srv.ForceTlsv13):
		push("force-tlsv13")
	case srv.ID != "":
		pushq("id", srv.ID)
	case srv.Inter != nil:
		pushi("inter", srv.Inter)
	case srv.Maxconn != nil:
		pushi("maxconn", srv.Maxconn)
	case srv.Maxqueue != nil:
		pushi("maxqueue", srv.Maxqueue)
	case srv.Minconn != nil:
		pushi("minconn", srv.Minconn)
	case !enabled(srv.SslReuse):
		push("no-ssl-reuse")
	case enabled(srv.NoSslv3):
		push("no-sslv3")
	case enabled(srv.NoTlsv10):
		push("no-tlsv10")
	case enabled(srv.NoTlsv11):
		push("no-tlsv11")
	case enabled(srv.NoTlsv12):
		push("no-tlsv12")
	case enabled(srv.NoTlsv13):
		push("no-tlsv13")
	case !enabled(srv.TLSTickets):
		push("no-tls-tickets")
	case srv.Npn != "":
		pushq("npm", srv.Npn)
	case srv.Observe != "":
		pushq("observe", srv.Observe)
	case srv.OnError != "":
		pushq("on-error", srv.OnError)
	case srv.OnMarkedDown != "":
		pushq("on-marked-down", srv.OnMarkedDown)
	case srv.OnMarkedUp != "":
		pushq("on-marked-up", srv.OnMarkedUp)
	case srv.PoolLowConn != nil:
		pushi("pool-low-conn", srv.PoolLowConn)
	case srv.PoolMaxConn != nil:
		pushi("pool-max-conn", srv.PoolMaxConn)
	case srv.PoolPurgeDelay != nil:
		pushi("pool-purge-delay", srv.PoolPurgeDelay)
	case srv.Proto != "":
		pushq("proto", srv.Proto)
	case len(srv.ProxyV2Options) > 0:
		pushq("proxy-v2-options", strings.Join(srv.ProxyV2Options, ","))
	case srv.Rise != nil:
		pushi("rise", srv.Rise)
	case enabled(srv.SendProxy):
		push("send-proxy")
	case enabled(srv.SendProxyV2):
		push("send-proxy-v2")
	case enabled(srv.SendProxyV2Ssl):
		push("send-proxy-v2-ssl")
	case enabled(srv.SendProxyV2SslCn):
		push("send-proxy-v2-ssl-cn")
	case srv.Slowstart != nil:
		pushi("slowstart", srv.Slowstart)
	case srv.Sni != "":
		pushq("sni", srv.Sni)
	case srv.Source != "":
		pushq("source", srv.Source)
	case enabled(srv.Ssl):
		push("ssl")
	case srv.SslMaxVer != "":
		pushq("ssl-max-ver", srv.SslMaxVer)
	case srv.SslMinVer != "":
		pushq("ssl-min-ver", srv.SslMinVer)
	case enabled(srv.Tfo):
		push("tfo")
	case enabled(srv.TLSTickets):
		push("tls-tickets")
	case srv.Track != "":
		pushq("track", srv.Track)
	/* XXX usesrc is not supported */
	case srv.Verify != "":
		pushq("verify", srv.Verify)
	case srv.Verifyhost != "":
		pushq("verifyhost", srv.Verifyhost)
	case srv.Weight != nil:
		pushi("weight", srv.Weight)
	case srv.Ws != "":
		pushq("ws", srv.Ws)
	}

	return b.String()
}
