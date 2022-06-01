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
	"fmt"
	"strings"

	"github.com/nathanaelle/syslog5424/v2"
	"github.com/sirupsen/logrus"
)

type RFC5424Hook struct {
	syslog *syslog5424.Syslog
	sender *syslog5424.Sender
	msgID  string
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

	r.syslog.Channel(sev).Msgid(r.msgID).Log(msg)

	return
}

func NewRFC5424Hook(opts Target) (logrus.Hook, error) {
	if len(opts.SyslogAddr) == 0 {
		return nil, fmt.Errorf("no address has been declared")
	}

	priority := strings.Join([]string{opts.SyslogFacility, opts.SyslogLevel}, ".")
	var priorityParsed syslog5424.Priority
	if err := priorityParsed.Set(priority); err != nil {
		return nil, err
	}

	slConn, _, err := syslog5424.Dial(opts.SyslogProto, opts.SyslogAddr)
	if err != nil {
		return nil, err
	}

	syslogServer, err := syslog5424.New(slConn, priorityParsed, opts.SyslogTag)
	if err != nil {
		return nil, err
	}

	return &RFC5424Hook{syslog: syslogServer, sender: slConn, msgID: opts.SyslogMsgID}, nil
}

func (r RFC5424Hook) Close() error {
	r.sender.End()
	return nil
}
