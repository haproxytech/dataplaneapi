// Code generated from 'constructor.tmpl'; DO NOT EDIT.

package acme

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/libdns/azure"
	"github.com/libdns/cloudflare"
	"github.com/libdns/cloudns"
	"github.com/libdns/digitalocean"
	"github.com/haproxytech/dataplaneapi/acme/exec"
	"github.com/libdns/gandi"
	"github.com/libdns/godaddy"
	"github.com/libdns/googleclouddns"
	"github.com/libdns/hetzner"
	"github.com/libdns/infomaniak"
	"github.com/libdns/inwx"
	"github.com/libdns/ionos"
	"github.com/libdns/linode"
	"github.com/libdns/namecheap"
	"github.com/libdns/netcup"
	"github.com/libdns/ovh"
	"github.com/libdns/porkbun"
	"github.com/libdns/rfc2136"
	"github.com/libdns/route53"
	"github.com/libdns/scaleway"
	"github.com/libdns/vultr/v2"
)

func NewDNSProvider(name string, params map[string]any) (DNSProvider, error) {
	var prov DNSProvider

	switch name {
	case "azure":
		prov = &azure.Provider{}
	case "cloudflare":
		prov = &cloudflare.Provider{}
	case "cloudns":
		prov = &cloudns.Provider{}
	case "digitalocean":
		prov = &digitalocean.Provider{}
	case "exec":
		prov = &exec.Provider{}
	case "gandi":
		prov = &gandi.Provider{}
	case "godaddy":
		prov = &godaddy.Provider{}
	case "googleclouddns":
		prov = &googleclouddns.Provider{}
	case "hetzner":
		prov = &hetzner.Provider{}
	case "infomaniak":
		prov = &infomaniak.Provider{}
	case "inwx":
		prov = &inwx.Provider{}
	case "ionos":
		prov = &ionos.Provider{}
	case "linode":
		prov = &linode.Provider{}
	case "namecheap":
		prov = &namecheap.Provider{}
	case "netcup":
		prov = &netcup.Provider{}
	case "ovh":
		prov = &ovh.Provider{}
	case "porkbun":
		prov = &porkbun.Provider{}
	case "rfc2136":
		prov = &rfc2136.Provider{}
	case "route53":
		prov = &route53.Provider{}
	case "scaleway":
		prov = &scaleway.Provider{}
	case "vultr":
		prov = &vultr.Provider{}
	default:
		return nil, fmt.Errorf("invalid DNS provider name: '%s'", name)
	}

	jsoni := jsoniter.ConfigCompatibleWithStandardLibrary
	js, err := jsoni.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params for DNS provider %s: %w", name, err)
	}
	if err = jsoni.Unmarshal(js, prov); err != nil {
		return nil, fmt.Errorf("invalid params for DNS provider %s: %w", name, err)
	}

	return prov, nil
}
