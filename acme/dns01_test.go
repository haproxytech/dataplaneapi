package acme

import (
	"reflect"
	"testing"

	"github.com/libdns/cloudflare"
	"github.com/libdns/googleclouddns"
	"github.com/libdns/ovh"
)

func TestNewDNSProvider(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		want    DNSProvider
		wantErr bool
	}{
		{
			name: "cloudflare",
			args: map[string]any{
				"api_token":  "foo",
				"zone_token": "bar",
			},
			want:    &cloudflare.Provider{APIToken: "foo", ZoneToken: "bar"},
			wantErr: false,
		},
		{
			name: "googleclouddns",
			args: map[string]any{
				"gcp_project":             "Project X",
				"gcp_application_default": `{"gcp_account_id":"j9h8hl094756h98990h"}`,
			},
			want:    &googleclouddns.Provider{Project: "Project X", ServiceAccountJSON: `{"gcp_account_id":"j9h8hl094756h98990h"}`},
			wantErr: false,
		},
		{
			name: "ovh",
			args: map[string]any{
				"endpoint":        "/lol",
				"application_key": "foobar",
			},
			want:    &ovh.Provider{Endpoint: "/lol", ApplicationKey: "foobar"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDNSProvider(tt.name, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDNSProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDNSProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guessZone(t *testing.T) {
	tests := []struct {
		fqdn string
		want string
	}{
		{
			fqdn: "",
			want: ".",
		},
		{
			fqdn: "haproxy.org",
			want: "haproxy.org.",
		},
		{
			fqdn: "foo.haproxy.org",
			want: "haproxy.org.",
		},
		{
			fqdn: "*.haproxy.org.",
			want: "haproxy.org.",
		},
		{
			fqdn: "*.foo.haproxy.org",
			want: "haproxy.org.",
		},
		{
			fqdn: "localhost",
			want: "localhost.",
		},
		{
			fqdn: "very.long.sub.domain.name.haproxy.lol",
			want: "haproxy.lol.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.fqdn, func(t *testing.T) {
			got := guessZone(tt.fqdn)
			if got != tt.want {
				t.Errorf("guessZone() = %v, want %v", got, tt.want)
			}
		})
	}
}
