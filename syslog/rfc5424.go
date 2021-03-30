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

	"github.com/nathanaelle/syslog5424/v2"
	"github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi/configuration"
)

type RFC5424Hook struct {
	syslog *syslog5424.Syslog
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

func NewRFC5424Hook(opts configuration.SyslogOptions) (logrus.Hook, error) {
	if len(opts.SyslogAddr) == 0 {
		return nil, fmt.Errorf("no address has been declared")
	}

	var severity syslog5424.Priority
	switch strings.ToLower(opts.SyslogLevel) {
	case "debug":
		severity = syslog5424.LogDEBUG
	case "info":
		severity = syslog5424.LogINFO
	case "notice":
		severity = syslog5424.LogNOTICE
	case "warning":
		severity = syslog5424.LogWARNING
	case "error":
		severity = syslog5424.LogERR
	case "critical":
		severity = syslog5424.LogCRIT
	case "alert":
		severity = syslog5424.LogALERT
	case "emergency":
		severity = syslog5424.LogEMERG
	default:
		return nil, fmt.Errorf("unrecognized severity: %s", opts.SyslogLevel)
	}

	var facility syslog5424.Priority
	switch opts.SyslogFacility {
	case "kern":
		facility = syslog5424.LogKERN
	case "user":
		facility = syslog5424.LogUSER
	case "mail":
		facility = syslog5424.LogMAIL
	case "daemon":
		facility = syslog5424.LogDAEMON
	case "auth":
		facility = syslog5424.LogAUTH
	case "syslog":
		facility = syslog5424.LogSYSLOG
	case "lpr":
		facility = syslog5424.LogLPR
	case "news":
		facility = syslog5424.LogNEWS
	case "uucp":
		facility = syslog5424.LogUUCP
	case "cron":
		facility = syslog5424.LogCRON
	case "authpriv":
		facility = syslog5424.LogAUTHPRIV
	case "ftp":
		facility = syslog5424.LogFTP
	case "local0":
		facility = syslog5424.LogLOCAL0
	case "local1":
		facility = syslog5424.LogLOCAL1
	case "local2":
		facility = syslog5424.LogLOCAL2
	case "local3":
		facility = syslog5424.LogLOCAL3
	case "local4":
		facility = syslog5424.LogLOCAL4
	case "local5":
		facility = syslog5424.LogLOCAL5
	case "local6":
		facility = syslog5424.LogLOCAL6
	case "local7":
		facility = syslog5424.LogLOCAL7
	default:
		return nil, fmt.Errorf("unrecognized facility: %s", opts.SyslogFacility)
	}

	slConn, chErr, err := syslog5424.Dial(opts.SyslogProto, opts.SyslogAddr)
	if err != nil {
		fmt.Printf("error establishing syslog output: %s\n", err)
	}

	go func(ch <-chan error) {
		for i := range ch {
			fmt.Printf("Error received from the syslog server: %s\n", i.Error())
		}
	}(chErr)

	syslogServer, err := syslog5424.New(slConn, facility|severity, opts.SyslogTag)
	if err != nil {
		return nil, err
	}

	return &RFC5424Hook{syslog: syslogServer, msgID: opts.SyslogMsgID}, nil
}
