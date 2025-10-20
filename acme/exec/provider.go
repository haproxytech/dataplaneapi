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

package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	osexec "os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

type Provider struct {
	Command     string   `json:"command"`
	Environment []string `json:"environment,omitempty"`
}

func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	cmd := osexec.CommandContext(ctx, p.Command) // #nosec G204
	p.populateEnv(cmd, "get", zone, nil)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result []libdns.RR
	err = json.Unmarshal(out, &result)
	if err != nil || len(result) == 0 {
		return nil, err
	}

	records := make([]libdns.Record, 0, len(result))
	for _, rr := range result {
		record, err := rr.Parse()
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return p.doRecords(ctx, "append", zone, records)
}

func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return p.doRecords(ctx, "set", zone, records)
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	return p.doRecords(ctx, "delete", zone, records)
}

func (p *Provider) doRecords(ctx context.Context, action, zone string, records []libdns.Record) ([]libdns.Record, error) {
	for _, record := range records {
		rr := record.RR()
		cmd := osexec.CommandContext(ctx, p.Command) // #nosec G204
		p.populateEnv(cmd, action, zone, &rr)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if _, ok := err.(*osexec.ExitError); ok {
				if msg := strings.TrimSpace(stderr.String()); len(msg) > 0 {
					return nil, fmt.Errorf("%v (stderr: %s)", err, msg)
				}
			}
			return nil, err
		}
	}
	return records, nil
}

func (p *Provider) populateEnv(cmd *osexec.Cmd, action, zone string, rr *libdns.RR) {
	env := cmd.Environ()
	if len(p.Environment) > 0 {
		env = append(env, p.Environment...)
	}
	env = append(env, "ACTION="+action, "ZONE="+zone)
	if rr != nil {
		env = append(env, rr2env(rr)...)
	}
	cmd.Env = env
}

func rr2env(rr *libdns.RR) []string {
	return []string{
		"REC_NAME=" + rr.Name,
		"REC_TTL=" + strconv.Itoa(int(rr.TTL/time.Second)),
		"REC_TYPE=" + rr.Type,
		"REC_DATA=" + rr.Data,
	}
}
