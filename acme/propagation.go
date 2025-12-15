// Copyright 2025 HAProxy Technologies
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

// This file contains code adapted from certmagic by Matt Holt.
// https://github.com/caddyserver/certmagic
//
// It has been modified.

package acme

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var dnsTimeout = 10 * time.Second

// FindZoneByFQDN determines the zone apex for the given fully-qualified
// domain name (FQDN) by recursing up the domain labels until the nameserver
// returns a SOA record in the answer section.
func FindZoneByFQDN(ctx context.Context, fqdn string, nameservers []string) (string, error) {
	if !strings.HasSuffix(fqdn, ".") {
		fqdn += "."
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	soa, err := fetchSoaByFqdn(ctx, fqdn, nameservers)
	if err != nil {
		return "", err
	}

	return soa.Hdr.Name, nil
}

func fetchSoaByFqdn(ctx context.Context, fqdn string, nameservers []string) (*dns.SOA, error) {
	var err error
	var in *dns.Msg

	labelIndexes := dns.Split(fqdn)
	for _, index := range labelIndexes {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		domain := fqdn[index:]

		in, err = dnsQuery(ctx, domain, dns.TypeSOA, nameservers, true)
		if err != nil {
			continue
		}
		if in == nil {
			continue
		}

		switch in.Rcode {
		case dns.RcodeSuccess:
			// Check if we got a SOA RR in the answer section
			if len(in.Answer) == 0 {
				continue
			}

			// CNAME records cannot/should not exist at the root of a zone.
			// So we skip a domain when a CNAME is found.
			if dnsMsgContainsCNAME(in) {
				continue
			}

			for _, ans := range in.Answer {
				if soa, ok := ans.(*dns.SOA); ok {
					return soa, nil
				}
			}
		case dns.RcodeNameError:
			// NXDOMAIN
		default:
			// Any response code other than NOERROR and NXDOMAIN is treated as error
			return nil, fmt.Errorf("unexpected response code '%s' for %s", dns.RcodeToString[in.Rcode], domain)
		}
	}

	return nil, fmt.Errorf("could not find the start of authority for %s%s", fqdn, formatDNSError(in, err))
}

func formatDNSError(msg *dns.Msg, err error) string {
	var parts []string
	if msg != nil {
		parts = append(parts, dns.RcodeToString[msg.Rcode])
	}
	if err != nil {
		parts = append(parts, err.Error())
	}
	if len(parts) > 0 {
		return ": " + strings.Join(parts, " ")
	}
	return ""
}

// checkDNSPropagation checks if the expected record has been propagated to all authoritative nameservers.
func checkDNSPropagation(ctx context.Context, fqdn string, recType uint16, expectedValue string, checkAuthoritativeServers bool, resolvers []string) (bool, error) {
	if !strings.HasSuffix(fqdn, ".") {
		fqdn += "."
	}

	// Initial attempt to resolve at the recursive NS - but do not actually
	// dereference (follow) a CNAME record if we are targeting a CNAME record
	// itself
	if recType != dns.TypeCNAME {
		r, err := dnsQuery(ctx, fqdn, recType, resolvers, true)
		if err != nil {
			return false, fmt.Errorf("CNAME dns query: %v", err)
		}
		if r.Rcode == dns.RcodeSuccess {
			fqdn = updateDomainWithCName(r, fqdn)
		}
	}

	if checkAuthoritativeServers {
		authoritativeServers, err := lookupNameservers(ctx, fqdn, resolvers)
		if err != nil {
			return false, fmt.Errorf("looking up authoritative nameservers: %v", err)
		}
		populateNameserverPorts(authoritativeServers)
		resolvers = authoritativeServers
	}

	return checkAuthoritativeNss(ctx, fqdn, recType, expectedValue, resolvers)
}

// checkAuthoritativeNss queries each of the given nameservers for the expected record.
func checkAuthoritativeNss(ctx context.Context, fqdn string, recType uint16, expectedValue string, nameservers []string) (bool, error) {
	for _, ns := range nameservers {
		r, err := dnsQuery(ctx, fqdn, recType, []string{ns}, true)
		if err != nil {
			return false, fmt.Errorf("querying authoritative nameservers: %v", err)
		}

		if r.Rcode != dns.RcodeSuccess {
			if r.Rcode == dns.RcodeNameError || r.Rcode == dns.RcodeServerFailure {
				// if Present() succeeded, then it must show up eventually, or else
				// something is really broken in the DNS provider or their API;
				// no need for error here, simply have the caller try again
				return false, nil
			}
			return false, fmt.Errorf("NS %s returned %s for %s", ns, dns.RcodeToString[r.Rcode], fqdn)
		}

		for _, rr := range r.Answer {
			switch recType {
			case dns.TypeTXT:
				if txt, ok := rr.(*dns.TXT); ok {
					record := strings.Join(txt.Txt, "")
					if record == expectedValue {
						return true, nil
					}
				}
			case dns.TypeCNAME:
				if cname, ok := rr.(*dns.CNAME); ok {
					// TODO: whether a DNS provider assumes a trailing dot or not varies, and we may have to standardize this in libdns packages
					if strings.TrimSuffix(cname.Target, ".") == strings.TrimSuffix(expectedValue, ".") {
						return true, nil
					}
				}
			default:
				return false, fmt.Errorf("unsupported record type: %d", recType)
			}
		}
	}

	return false, nil
}

// lookupNameservers returns the authoritative nameservers for the given fqdn.
func lookupNameservers(ctx context.Context, fqdn string, resolvers []string) ([]string, error) {
	var authoritativeNss []string

	zone, err := FindZoneByFQDN(ctx, fqdn, resolvers)
	if err != nil {
		return nil, fmt.Errorf("could not determine the zone for '%s': %w", fqdn, err)
	}

	r, err := dnsQuery(ctx, zone, dns.TypeNS, resolvers, true)
	if err != nil {
		return nil, fmt.Errorf("querying NS resolver for zone '%s' recursively: %v", zone, err)
	}

	for _, rr := range r.Answer {
		if ns, ok := rr.(*dns.NS); ok {
			authoritativeNss = append(authoritativeNss, strings.ToLower(ns.Ns))
		}
	}

	if len(authoritativeNss) > 0 {
		return authoritativeNss, nil
	}
	return nil, errors.New("could not determine authoritative nameservers")
}

// Update FQDN with CNAME if any
func updateDomainWithCName(r *dns.Msg, fqdn string) string {
	for _, rr := range r.Answer {
		if cn, ok := rr.(*dns.CNAME); ok {
			if cn.Hdr.Name == fqdn {
				return cn.Target
			}
		}
	}
	return fqdn
}

// dnsMsgContainsCNAME checks for a CNAME answer in msg
func dnsMsgContainsCNAME(msg *dns.Msg) bool {
	for _, ans := range msg.Answer {
		if _, ok := ans.(*dns.CNAME); ok {
			return true
		}
	}
	return false
}

func dnsQuery(ctx context.Context, fqdn string, rtype uint16, nameservers []string, recursive bool) (*dns.Msg, error) {
	m := createDNSMsg(fqdn, rtype, recursive)
	var in *dns.Msg
	var err error
	for _, ns := range nameservers {
		in, err = sendDNSQuery(ctx, m, ns)
		if err == nil && len(in.Answer) > 0 {
			break
		}
	}
	return in, err
}

func createDNSMsg(fqdn string, rtype uint16, recursive bool) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(fqdn, rtype)

	// See: https://caddy.community/t/hard-time-getting-a-response-on-a-dns-01-challenge/15721/16
	m.SetEdns0(1232, false)
	if !recursive {
		m.RecursionDesired = false
	}
	return m
}

func sendDNSQuery(ctx context.Context, m *dns.Msg, ns string) (*dns.Msg, error) {
	udp := &dns.Client{Net: "udp", Timeout: dnsTimeout}
	in, _, err := udp.ExchangeContext(ctx, m, ns)
	// two kinds of errors we can handle by retrying with TCP:
	// truncation and timeout; see https://github.com/caddyserver/caddy/issues/3639
	truncated := in != nil && in.Truncated
	timeoutErr := err != nil && strings.Contains(err.Error(), "timeout")
	if truncated || timeoutErr {
		tcp := &dns.Client{Net: "tcp", Timeout: dnsTimeout}
		in, _, err = tcp.ExchangeContext(ctx, m, ns)
	}
	return in, err
}

// RecursiveNameservers are used to pre-check DNS propagation. It
// picks user-configured nameservers (custom) OR the defaults
// obtained from resolv.conf and defaultNameservers if none is
// configured and ensures that all server addresses have a port value.
func RecursiveNameservers(custom []string) []string {
	var servers []string
	if len(custom) == 0 {
		servers = systemOrDefaultNameservers(defaultResolvConf, defaultNameservers)
	} else {
		servers = make([]string, len(custom))
		copy(servers, custom)
	}
	populateNameserverPorts(servers)
	return servers
}

// systemOrDefaultNameservers attempts to get system nameservers from the
// resolv.conf file given by path before falling back to hard-coded defaults.
func systemOrDefaultNameservers(path string, defaults []string) []string {
	config, err := dns.ClientConfigFromFile(path)
	if err != nil || len(config.Servers) == 0 {
		return defaults
	}
	return config.Servers
}

// populateNameserverPorts ensures that all nameservers have a port number.
// If not, the default DNS server port of 53 will be appended.
func populateNameserverPorts(servers []string) {
	for i := range servers {
		_, port, _ := net.SplitHostPort(servers[i])
		if port == "" {
			servers[i] = net.JoinHostPort(servers[i], "53")
		}
	}
}

var defaultNameservers = []string{
	"8.8.8.8:53",
	"8.8.4.4:53",
	"1.1.1.1:53",
	"1.0.0.1:53",
}

const defaultResolvConf = "/etc/resolv.conf"
