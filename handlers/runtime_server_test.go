package handlers

import (
	"testing"

	"github.com/haproxytech/client-native/v6/models"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/dataplaneapi/misc"
)

func TestSerializeRuntimeAddServer(t *testing.T) {
	tests := []struct {
		name    string
		srv     *models.RuntimeAddServer
		version *cn_runtime.HAProxyVersion
		want    string
	}{
		{
			name: "basic server",
			srv: &models.RuntimeAddServer{
				Address: "127.0.0.1",
				Port:    misc.Int64P(8080),
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    " 127.0.0.1:8080",
		},
		{
			name: "server with weight",
			srv: &models.RuntimeAddServer{
				Address: "192.168.1.100",
				Port:    misc.Int64P(9000),
				Weight:  misc.Int64P(50),
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    " 192.168.1.100:9000 weight 50",
		},
		{
			name: "server with maintenance",
			srv: &models.RuntimeAddServer{
				Address:     "10.0.0.10",
				Maintenance: "enabled",
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    " 10.0.0.10 disabled",
		},
		{
			name: "server with maintenance disabled",
			srv: &models.RuntimeAddServer{
				Address:     "10.0.0.10",
				Maintenance: "disabled",
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    " 10.0.0.10 enabled",
		},
		{
			name: "server with maintenance disabled, HAProxy 3.0",
			srv: &models.RuntimeAddServer{
				Address:     "10.0.0.10",
				Maintenance: "disabled",
			},
			version: func() *cn_runtime.HAProxyVersion {
				v := new(cn_runtime.HAProxyVersion)
				v.ParseHAProxyVersion("3.0.0")
				return v
			}(),
			want: " 10.0.0.10",
		},
		{
			name: "server with agent check",
			srv: &models.RuntimeAddServer{
				Address:    "172.16.0.5",
				AgentCheck: "enabled",
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    " 172.16.0.5 agent-check",
		},
		{
			name: "server with agent addr",
			srv: &models.RuntimeAddServer{
				Address:   "172.16.0.6",
				AgentAddr: "127.0.0.1",
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    ` 172.16.0.6 agent-addr 127.0.0.1`,
		},
		{
			name: "server with multiple options",
			srv: &models.RuntimeAddServer{
				Address:            "10.1.1.10",
				Port:               misc.Int64P(80),
				Weight:             misc.Int64P(10),
				Check:              "enabled",
				Backup:             "enabled",
				Maintenance:        "enabled",
				AgentCheck:         "enabled",
				AgentAddr:          "127.0.0.1",
				AgentPort:          misc.Int64P(5000),
				HealthCheckAddress: "127.0.0.2",
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    ` 10.1.1.10:80 agent-check agent-addr 127.0.0.1 agent-port 5000 backup check addr 127.0.0.2 disabled weight 10`,
		},
		{
			name: "server with all fields",
			srv: &models.RuntimeAddServer{
				Address:            "10.1.1.10",
				Port:               misc.Int64P(80),
				Weight:             misc.Int64P(10),
				Check:              "enabled",
				Backup:             "enabled",
				Maintenance:        "enabled",
				AgentCheck:         "enabled",
				AgentAddr:          "127.0.0.1",
				AgentPort:          misc.Int64P(5000),
				AgentInter:         misc.Int64P(1000),
				AgentSend:          "foobar",
				Allow0rtt:          true,
				Alpn:               "h2,http/1.1",
				CheckAlpn:          "h2",
				CheckProto:         "HTTP",
				CheckSendProxy:     "enabled",
				CheckSni:           "example.com",
				CheckSsl:           "enabled",
				CheckViaSocks4:     "enabled",
				Ciphers:            "HIGH:!aNULL:!MD5",
				Ciphersuites:       "TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256",
				CrlFile:            "/path/to/crl.pem",
				SslCertificate:     "/path/to/cert.pem",
				Downinter:          misc.Int64P(2000),
				ErrorLimit:         misc.Int64P(10),
				Fall:               misc.Int64P(2),
				Fastinter:          misc.Int64P(500),
				ForceSslv3:         "enabled",
				ForceTlsv10:        "enabled",
				ForceTlsv11:        "enabled",
				ForceTlsv12:        "enabled",
				ForceTlsv13:        "enabled",
				HealthCheckAddress: "127.0.0.2",
				HealthCheckPort:    misc.Int64P(8080),
				Inter:              misc.Int64P(3000),
				Maxconn:            misc.Int64P(100),
				Maxqueue:           misc.Int64P(200),
				Minconn:            misc.Int64P(50),
				Rise:               misc.Int64P(1),
			},
			version: &cn_runtime.HAProxyVersion{},
			want:    ` 10.1.1.10:80 agent-check agent-addr 127.0.0.1 agent-port 5000 agent-inter 1000 agent-send foobar allow-0rtt alpn h2,http/1.1 backup check check-alpn h2 addr 127.0.0.2 port 8080 check-proto HTTP check-send-proxy check-sni example.com check-ssl check-via-socks4 ciphers HIGH:!aNULL:!MD5 ciphersuites TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256 crl-file /path/to/crl.pem crt /path/to/cert.pem disabled downinter 2000 error-limit 10 fall 2 fastinter 500 force-sslv3 force-tlsv10 force-tlsv11 force-tlsv12 force-tlsv13 inter 3000 maxconn 100 maxqueue 200 minconn 50 rise 1 weight 10`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.version == nil {
				tt.version = &cn_runtime.HAProxyVersion{}
			}
			if got := SerializeRuntimeAddServer(tt.srv, tt.version); got != tt.want {
				t.Errorf("SerializeRuntimeAddServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
