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

package syslog

import (
	"fmt"
	"strings"
	"time"

	"github.com/nathanaelle/syslog5424/v2"
	"github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi/configuration"
)

type RFC5424Hook struct {
	syslog *syslog5424.Syslog
}

func (r RFC5424Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (r RFC5424Hook) Fire(entry *logrus.Entry) (err error) {
	var sev syslog5424.Priority
	switch entry.Level {
	case logrus.PanicLevel:
		sev = syslog5424.LogCRIT
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
		messages = append(messages, fmt.Sprintf("%s=%v", k, v))
	}

	msg := strings.Join(messages, " ")

	_, err = r.syslog.Channel(sev).Write([]byte(msg))

	return
}

func NewRFC5424Hook(opts configuration.SyslogOptions) (logrus.Hook, error) {
	if len(opts.SyslogSrv) == 0 {
		return nil, fmt.Errorf("no server has been declared")
	}

	var priority syslog5424.Priority
	switch strings.ToLower(opts.SyslogPriority) {
	case "debug":
		priority = syslog5424.LogDEBUG
	case "info":
		priority = syslog5424.LogINFO
	case "notice":
		priority = syslog5424.LogNOTICE
	case "warning":
		priority = syslog5424.LogWARNING
	case "error":
		priority = syslog5424.LogERR
	case "critical":
		priority = syslog5424.LogCRIT
	case "alert":
		priority = syslog5424.LogALERT
	case "emergency":
		priority = syslog5424.LogEMERG
	default:
		return nil, fmt.Errorf("unrecognized severity: %s", opts.SyslogPriority)
	}

	var err error
	var facility syslog5424.Priority
	if err = facility.Set(opts.SyslogFacility); err != nil {
		return nil, err
	}

	connector := syslog5424.TCPConnector(opts.SyslogProto, fmt.Sprintf("%s:%d", opts.SyslogSrv, opts.SyslogPort))
	slConn, chErr := syslog5424.NewSender(connector, syslog5424.TransportRFC5425, time.NewTicker(500*time.Millisecond).C)

	go func(ch <-chan error) {
		for i := range ch {
			fmt.Printf("Error received from the syslog server: %s", i.Error())
		}
	}(chErr)

	syslogServer, err := syslog5424.New(slConn, facility|priority, opts.SyslogTag)
	if err != nil {
		return nil, err
	}

	return &RFC5424Hook{syslog: syslogServer}, nil
}
