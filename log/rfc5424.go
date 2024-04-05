// Copyright 2021 HAProxy Technologies
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

package log

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nathanaelle/syslog5424/v2"
	circuit "github.com/rubyist/circuitbreaker"
	"github.com/sirupsen/logrus"
)

type RFC5424Hook struct {
	syslog *syslog5424.Syslog
	sender *syslog5424.Sender
	msgID  string

	// Use a circuit breaker to pause sending messages to the syslog target
	// in the presence of connection errors.
	cb *circuit.Breaker
}

func (r RFC5424Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (r RFC5424Hook) Fire(entry *logrus.Entry) (err error) {
	var sev syslog5424.Priority
	switch entry.Level {
	case logrus.PanicLevel:
		sev = syslog5424.LogALERT
	case logrus.FatalLevel:
		sev = syslog5424.LogCRIT
	case logrus.ErrorLevel:
		sev = syslog5424.LogERR
	case logrus.WarnLevel:
		sev = syslog5424.LogWARNING
	case logrus.InfoLevel:
		sev = syslog5424.LogINFO
	case logrus.DebugLevel, logrus.TraceLevel:
		sev = syslog5424.LogDEBUG
	}

	messages := []string{entry.Message}
	for k, v := range entry.Data {
		// TODO: we should deal with structured data properly
		messages = append(messages, fmt.Sprintf("%s=%v", k, v))
	}

	msg := strings.Join(messages, " ")

	// Do not perform any action unless the circuit breaker is either closed (reset), or is ready to retry.
	if r.cb.Ready() {
		r.syslog.Channel(sev).Msgid(r.msgID).Log(msg)
		// Register any call as successful to enable automatic resets.
		// Failures are registered asynchronously by the goroutine that consumes errors from the corresponding channel.
		r.cb.Success()
	}

	return
}

func NewRFC5424Hook(opts Target) (logrus.Hook, error) {
	if len(opts.SyslogAddr) == 0 {
		return nil, errors.New("no address has been declared")
	}

	priority := strings.Join([]string{opts.SyslogFacility, opts.SyslogLevel}, ".")
	var priorityParsed syslog5424.Priority
	if err := priorityParsed.Set(priority); err != nil {
		return nil, err
	}

	// syslog5424.Dial() returns an error channel, which needs to be drained
	// in order to avoid blocking.
	slConn, errCh, err := syslog5424.Dial(opts.SyslogProto, opts.SyslogAddr)
	if err != nil {
		return nil, err
	}

	syslogServer, err := syslog5424.New(slConn, priorityParsed, opts.SyslogTag)
	if err != nil {
		return nil, err
	}

	r := &RFC5424Hook{
		syslog: syslogServer, sender: slConn, msgID: opts.SyslogMsgID,
		// We can change the circuit breaker settings as desired - including making
		// them configurable and/or dynamically adjustable based on runtime conditions.
		//
		// Please note, however, that a 3-failure threshold breaker with default settings
		// was found to work well with varying load and different states of a log target.
		// Specifically, the breaker will remain tripped when sending messages to the target
		// that is consistently failing, and will reset quickly when delivery begins to succeed.
		cb: circuit.NewThresholdBreaker(3),
	}

	// A signal channel that is used to stop the goroutine reporting on circuit breaker state changes.
	doneCh := make(chan struct{})

	// Consume errors from errCh until it is closed.
	go func() {
		for {
			err, ok := <-errCh
			if err != nil {
				r.cb.Fail() // Register a failure with the circuit breaker.
			}
			if !ok {
				close(doneCh)
				return
			}
		}
	}()

	// Report on circuit breaker state changes.
	cbStateCh := r.cb.Subscribe()
	go func() {
		for {
			select {
			case e, ok := <-cbStateCh:
				if !ok {
					return
				}
				var state string
				switch e {
				case circuit.BreakerTripped:
					state = "too many connection errors, log delivery is stopped until this improves"
				case circuit.BreakerReset:
					state = "resuming log delivery"
				default:
					continue
				}
				fmt.Println(time.Now().Format(time.RFC3339), "syslog target", opts.SyslogAddr, "("+opts.SyslogTag+"):", state)
			case <-doneCh:
				return
			}
		}
	}()

	return r, nil
}

func (r RFC5424Hook) Close() error {
	r.sender.End() // This will also close errCh returned by syslog.Dial() in NewRFC5424Hook(), causing related goroutines to exit.
	return nil
}
