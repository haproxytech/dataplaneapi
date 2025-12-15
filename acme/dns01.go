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

package acme

//go:generate go run gen_constructor.go -i dns01-providers.txt -t constructor.tmpl -o constructor.go

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/miekg/dns"
)

const (
	// TTL of the temporary DNS record used for DNS-01 validation.
	DefaultTTL = 30 * time.Second
	// Typical negative response TTL defined in the SOA.
	defaultDNSPropagationTimeout = 300 * time.Second
)

// DNSProvider defines the operations required for dns-01 challenges.
type DNSProvider interface {
	libdns.RecordAppender
	libdns.RecordDeleter
}

// A DNS01Solver uses a DNSProvider to actually solve the challenge.
type DNS01Solver struct {
	provider DNSProvider
	TTL      time.Duration

	// How long to wait before starting propagation checks.
	// Default: 0 (no wait).
	PropagationDelay time.Duration

	// Maximum time to wait for temporary DNS record to appear.
	// Set to -1 to disable propagation checks.
	// Default: 2 minutes.
	PropagationTimeout time.Duration

	// Preferred DNS resolver(s) to use when doing DNS lookups.
	Resolvers []string
}

func NewDNS01Solver(name string, params map[string]any, ttl ...time.Duration) (*DNS01Solver, error) {
	prov, err := NewDNSProvider(name, params)
	if err != nil {
		return nil, err
	}

	recordTTL := DefaultTTL
	if len(ttl) > 0 {
		recordTTL = ttl[0]
	}

	return &DNS01Solver{provider: prov, TTL: recordTTL}, nil
}

// Present creates the DNS TXT record for the given ACME challenge.
func (s *DNS01Solver) Present(ctx context.Context, domain, zone, keyAuth string) error {
	rec := makeRecord(domain, keyAuth, s.TTL)

	if zone == "" {
		zone = GuessZone(domain)
	} else {
		zone = rooted(zone)
	}

	results, err := s.provider.AppendRecords(ctx, zone, []libdns.Record{rec})
	if err != nil {
		return fmt.Errorf("adding temporary record for zone %q: %w", zone, err)
	}
	if len(results) != 1 {
		return fmt.Errorf("expected one record, got %d: %v", len(results), results)
	}

	return nil
}

// Wait blocks until the TXT record created in Present() appears in
// authoritative lookups, i.e. until it has propagated, or until
// timeout, whichever is first.
func (s *DNS01Solver) Wait(ctx context.Context, domain, zone, keyAuth string) error {
	// if configured to, pause before doing propagation checks
	// (even if they are disabled, the wait might be desirable on its own)
	if s.PropagationDelay > 0 {
		select {
		case <-time.After(s.PropagationDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// skip propagation checks if configured to do so
	if s.PropagationTimeout == -1 {
		return nil
	}

	// timings
	timeout := s.PropagationTimeout
	if timeout == 0 {
		timeout = defaultDNSPropagationTimeout
	}
	const interval = 5 * time.Second

	// how we'll do the checks
	checkAuthoritativeServers := len(s.Resolvers) == 0
	resolvers := RecursiveNameservers(s.Resolvers)

	absName := strings.Trim(domain, ".")

	var err error
	start := time.Now()
	for time.Since(start) < timeout {
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return ctx.Err()
		}

		var ready bool
		ready, err = checkDNSPropagation(ctx, absName, dns.TypeTXT, keyAuth, checkAuthoritativeServers, resolvers)
		if err != nil {
			return fmt.Errorf("checking DNS propagation of %q (resolvers=%v): %w", absName, resolvers, err)
		}
		if ready {
			return nil
		}
	}

	return fmt.Errorf("DNS propagation timed out. Last error: %v", err)
}

// CleanUp deletes the DNS TXT record created in Present().
func (s *DNS01Solver) CleanUp(ctx context.Context, domain, zone, keyAuth string) error {
	rr := makeRecord(domain, keyAuth, s.TTL)

	if zone == "" {
		zone = GuessZone(domain)
	} else {
		zone = rooted(zone)
	}

	_, err := s.provider.DeleteRecords(ctx, zone, []libdns.Record{rr})
	if err != nil {
		return fmt.Errorf("deleting temporary record for name %q in zone %q: %w", zone, rr, err)
	}

	return nil
}

// Assemble a TXT Record suited for DNS-01 challenges.
func makeRecord(fqdn, keyAuth string, ttl time.Duration) libdns.RR {
	return libdns.RR{
		Type: "TXT",
		Name: "_acme-challenge." + trimWildcard(fqdn),
		Data: keyAuth,
		TTL:  ttl,
	}
}

// Guess the root zone for a domain when we cannot use a better method.
func GuessZone(fqdn string) string {
	fqdn = trimWildcard(fqdn)
	parts := make([]string, 0, 8)
	strings.SplitSeq(fqdn, ".")(func(part string) bool {
		if part != "" {
			parts = append(parts, part)
		}
		return true
	})

	n := len(parts)
	if n < 3 {
		return rooted(fqdn)
	}
	return rooted(strings.Join(parts[n-2:], "."))
}

// Remove the wildcard from a domain so it can be used in a record name.
func trimWildcard(fqdn string) string {
	fqdn = strings.TrimSpace(fqdn)
	return strings.TrimPrefix(fqdn, "*.")
}

// Ensures a domain name has its final dot (the root zone).
func rooted(domain string) string {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	return domain
}
