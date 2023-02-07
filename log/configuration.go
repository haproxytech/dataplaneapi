package log

const DefaultApacheLogFormat = "%h %l %u %t \"%m %U%q %H\" %>s %b \"%{Referer}i\" \"%{User-agent}i\""

type Target struct {
	LogTo          string   `yaml:"log_to,omitempty" description:"Log target, can be stdout, file, or syslog"`
	LogFile        string   `yaml:"log_file,omitempty" description:"Location of the log file"`
	LogLevel       string   `yaml:"log_level,omitempty" description:"Logging level, allowed values: trace|debug|info|warning|error"`
	LogFormat      string   `yaml:"log_format,omitempty" description:"Logging format, allowed values: text|JSON"`
	ACLFormat      string   `yaml:"acl_format,omitempty" description:"Apache Common Log Format to format the access log entries, default:\"%h %l %u %t \\\"%r\\\" %>s %b \\\"%{Referer}i\\\" \\\"%{User-agent}i\\\" %{us}T"`
	SyslogAddr     string   `yaml:"syslog_address,omitempty" description:"Syslog address (with port declaration in case of TCP type) where logs should be forwarded: accepting socket path in case of unix or unixgram"`
	SyslogProto    string   `yaml:"syslog_protocol,omitempty" description:"Syslog server protocol, allowed values: tcp|tcp4|tcp6|unix|unixgram"`
	SyslogTag      string   `yaml:"syslog_tag,omitempty" description:"String to tag the syslog messages"`
	SyslogLevel    string   `yaml:"syslog_level,omitempty" description:"Define the required syslog messages level, allowed values: debug|info|notice|warning|error|critical|alert|emergency"`
	SyslogFacility string   `yaml:"syslog_facility,omitempty" description:"Define the Syslog facility number, allowed values: kern|user|mail|daemon|auth|syslog|lpr|news|uucp|cron|authpriv|ftp|local0|local1|local2|local3|local4|local5|local6|local7"`
	SyslogMsgID    string   `yaml:"-"`
	LogTypes       []string `yaml:"log_types,omitempty" description:"Define which log types to log to this target, allowed values: app|access" save:"true"`
}

type Targets []Target

// Deprecated: use configuration file options instead
type LoggingOptions struct {
	LogTo     string `long:"log-to" description:"Log target, can be stdout, file, or syslog" default:"stdout" choice:"stdout" choice:"file" choice:"syslog" group:"log"`
	LogFile   string `long:"log-file" description:"Location of the log file" default:"/var/log/dataplaneapi/dataplaneapi.log" group:"log"`
	LogLevel  string `long:"log-level" description:"Logging level" default:"warning" choice:"trace" choice:"debug" choice:"info" choice:"warning" choice:"error" group:"log"`
	LogFormat string `long:"log-format" description:"Logging format" default:"text" choice:"text" choice:"JSON" group:"log"`
	ACLFormat string `long:"apache-common-log-format" description:"Apache Common Log Format to format the access log entries" default:"%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-agent}i\" %{us}T" group:"log"`
}

// Deprecated: use configuration file options instead
type SyslogOptions struct {
	SyslogAddr     string `long:"syslog-address" description:"Syslog address (with port declaration in case of TCP type) where logs should be forwarded: accepting socket path in case of unix or unixgram" default:"" group:"syslog"`
	SyslogProto    string `long:"syslog-protocol" description:"Syslog server protocol" default:"tcp" choice:"tcp" choice:"tcp4" choice:"tcp6" choice:"unix" choice:"unixgram" group:"syslog"`
	SyslogTag      string `long:"syslog-tag" description:"String to tag the syslog messages" default:"dataplaneapi" group:"syslog"`
	SyslogLevel    string `long:"syslog-level" description:"Define the required syslog messages level, allowed values: debug|info|notice|warning|error|critical|alert|emergency " default:"debug" group:"syslog"`
	SyslogFacility string `long:"syslog-facility" description:"Define the Syslog facility number, allowed values: kern|user|mail|daemon|auth|syslog|lpr|news|uucp|cron|authpriv|ftp|local0|local1|local2|local3|local4|local5|local6|local7" default:"local0" group:"syslog"`
	SyslogMsgID    string
}
