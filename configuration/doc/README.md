# HAProxy Data Plane API Configuration

## Overview

The primary configuration for the Data Plane API is managed through a YAML file. This file allows you to set various parameters that control the behavior of the API and its interaction with HAProxy.

## Top-Level Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `config_version` | integer | No | Configuration version |
| `name` | string | No | Name |

## Structure

The configuration file is structured into several sections, each corresponding to a different aspect of the Data Plane API's functionality. Here's a breakdown of the main sections and their options:

### dataplaneapi Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `write_timeout` | string | No | Specifies the write timeout for API requests |
| `graceful_timeout` | string | No | Sets the graceful shutdown timeout |
| `show_system_info` | boolean | No | Enables or disables the display of system information on the info endpoint |
| `max_header_size` | string | No | Defines the maximum header size for API requests |
| `socket_path` | string | No | Specifies the path to the API's Unix socket |
| `debug_socket_path` | string | No | Defines the path for the debugging command socket |
| `host` | string | No | Sets the host address for the API |
| `port` | integer | No | Defines the port number for the API |
| `listen_limit` | integer | No | Sets the maximum number of connections the API can handle |
| `disable_inotify` | boolean | No | Disables the inotify watcher for the configuration file |
| `read_timeout` | string | No | Specifies the read timeout for API requests |
| `advertised` | object | No | Contains settings for the advertised API address and port. [See advertised options](#dataplaneapiadvertised-configuration-options) |
| `cleanup_timeout` | string | No | Sets the timeout for cleanup operations |
| `keep_alive` | string | No | Defines the keep-alive timeout for connections |
| `pid_file` | string | No | Specifies the path to the file where the API's PID will be written |
| `uid` | integer | No | User id value to set on start |
| `gid` | integer | No | Group id value to set on start |
| `tls` | object | No | Contains settings for TLS encryption. [See TLS options](#dataplanapitls-configuration-options) |
| `scheme` | array of strings | No | Enabled listeners |
| `userlist` | object | No | Contains settings for userlist. [See userlist options](#dataplanapituserlist-configuration-options) |
| `transaction` | object | No | Contains settings for transactions. [See transaction options](#dataplanapitransaction-configuration-options) |
| `resources` | object | No | Contains settings for resources. [See resources options](#dataplaneapiresources-configuration-options) |
| `user` | array of objects | No | List of users. [See user options](#dataplaneapiuser-configuration-options) |

### dataplaneapi.advertised Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `api_address` | string | No | The advertised API address |
| `api_port` | integer | No | The advertised API port |

### dataplaneapi.tls Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `tls_host` | string | No | The host address for TLS connections |
| `tls_port` | integer | No | The port number for TLS connections |
| `tls_certificate` | string | No | Path to the TLS certificate file |
| `tls_key` | string | No | Path to the TLS certificate key file |
| `tls_ca` | string | No | Path to the TLS CA certificate file |
| `tls_listen_limit` | integer | No | The maximum number of TLS connections |
| `tls_keep_alive` | string | No | The keep-alive timeout for TLS connections |
| `tls_read_timeout` | string | No | The read timeout for TLS connections |
| `tls_write_timeout` | string | No | The write timeout for TLS connections |

### dataplaneapi.userlist Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `userlist` | string | No | Userlist in HAProxy configuration to use for API Basic Authentication |
| `userlist_file` | string | No | Path to the dataplaneapi userlist file |

### dataplaneapi.transaction Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `transaction_dir` | string | No | Path to the transaction directory |
| `backups_number` | integer | No | Number of backup configuration files to keep |
| `backups_dir` | string | No | Path to directory in which to place backup files |
| `max_open_transactions` | integer | No | Limit for active transaction in pending state |

### dataplaneapi.resources Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `maps_dir` | string | No | Path to directory of map files managed by dataplane |
| `ssl_certs_dir` | string | No | Path to SSL certificates directory |
| `general_storage_dir` | string | No | Path to general storage directory |
| `dataplane_storage_dir` | string | No | Path to dataplane internal storage directory |
| `update_map_files` | boolean | No | Flag used for syncing map files with runtime maps values |
| `update_map_files_period` | integer | No | Elapsed time in seconds between two maps syncing operations |
| `spoe_dir` | string | No | Path to SPOE directory |
| `spoe_transaction_dir` | string | No | Path to the SPOE transaction directory |

### dataplaneapi.user Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | string | Yes | User name |
| `password` | string | Yes | Password |
| `insecure` | boolean | No | Insecure password |

### haproxy Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `config_file` | string | No | Path to the HAProxy configuration file |
| `haproxy_bin` | string | No | Path to the HAProxy binary file |
| `master_runtime` | string | No | Path to the master Runtime API socket |
| `fid` | string | No | Path to file that will dataplaneapi use to write its id |
| `master_worker_mode` | boolean | No | Flag to enable helpers when running HAProxy in master worker mode |
| `reload` | object | No | Contains settings for reloading HAProxy. [See reload options](#haproxyreload-configuration-options) |
| `delayed_start_max` | string | No | Maximum duration to wait for the haproxy runtime socket to be ready |
| `delayed_start_tick` | string | No | Duration between checks for the haproxy runtime socket to be ready |

### haproxy.reload Configuration Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `reload_delay` | integer | No | Minimum delay between two reloads (in s) |
| `reload_cmd` | string | No | Reload command |
| `restart_cmd` | string | No | Restart command |
| `status_cmd` | string | No | Status command |
| `service_name` | string | No | Name of the HAProxy service |
| `reload_retention` | integer | No | Reload retention in days |
| `reload_strategy` | string | No | Either systemd, s6 or custom |
| `validate_cmd` | string | No | Executes a custom command to perform the HAProxy configuration check |

### log_targets Configuration Options

This section contains settings related to log targets.

-   **`log_targets`**: (array of objects, optional) - List of log targets.

    The `log_targets` option allows you to define multiple destinations for log messages. Each element in the `log_targets` array is an object that defines a specific log target. Each log target object can have the following properties:

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `log_to` | string | Yes | Type of log target: 'file', 'syslog', or 'stdout' |
| `log_format` | string | No | Log format for this target: 'text', 'json', or 'apache_common'. Default is 'text' |
| `log_level` | string | No | Log level: 'trace', 'debug', 'info', 'warning', 'error', or 'critical'. Default is 'info' |
| `log_file` | string | No | If target is 'file', specifies the path to the log file |
| `acl_format` | string | No | Apache Common Log Format to format the access log entries, default:\"%h %l %u %t \\\"%r\\\" %>s %b \\\"%{Referer}i\\\" \\\"%{User-agent}i\\\" %{us}T" |
| `syslog_address` | string | No | If target is 'syslog', specifies the syslog server address |
| `syslog_protocol` | string | No | If target is 'syslog', specifies the syslog protocol (e.g., 'udp', 'tcp') |
| `syslog_tag` | string | No | If target is 'syslog', specifies the syslog tag |
| `syslog_facility` | string | No | If target is 'syslog', define the Syslog facility number, allowed values: 'kern', 'user', 'mail', 'daemon', 'auth', 'syslog', 'lpr', 'news', 'uucp', 'cron', 'authpriv', 'ftp', 'local0', 'local1', 'local2', 'local3', 'local4', 'local5', 'local6', 'local7' |
| `syslog_level` | string | No | If target is 'syslog', define the required syslog messages level, allowed values: 'debug', 'info', 'notice', 'warning', 'error', 'critical', 'alert', 'emergency' |
| `log_types` | array of strings | No | Define which log types to log to this target, allowed values: 'app', 'access' |

**Example of `log_targets`:**

```yaml
log_targets:
  - log_target: file
    log_file: /var/log/dataplaneapi.log
    log_format: json
    log_level: debug
    log_types:
      - access
      - app
  - log_target: syslog
    syslog_address: 127.0.0.1:514
    syslog_protocol: udp
    syslog_tag: dataplaneapi
    syslog_facility: local0
    log_level: info
    log_types:
      - access
  - log_target: stdout
    log_format: text
    log_level: warning
    log_types:
      - access
      - app
```

## Command-Line Overrides

Many of the configuration options available in the YAML file can be overridden by command-line arguments when starting the Data Plane API. This allows for flexibility in deployment and testing.

### CLI Argument Precedence

When a configuration option is set both in the YAML file and via a command-line argument, the command-line argument takes precedence.

### Common CLI Flags

-   `-c`, `--config-file`: Overrides the `haproxy.config_file` option.
-   `-b`, `--haproxy-bin`: Overrides the `haproxy.haproxy_bin` option.
-   `-m`, `--master-runtime`: Overrides the `haproxy.master_runtime` option.
-   `-u`, `--userlist`: Overrides the `dataplaneapi.userlist.userlist` option.
-   `-d`, `--reload-delay`: Overrides the `haproxy.reload.reload_delay` option.
-   `-r`, `--reload-cmd`: Overrides the `haproxy.reload.reload_cmd` option.
-   `-s`, `--restart-cmd`: Overrides the `haproxy.reload.restart_cmd` option.
-   `-t`, `--transaction-dir`: Overrides the `dataplaneapi.transaction.transaction_dir` option.
-   `-n`, `--backups-number`: Overrides the `dataplaneapi.transaction.backups_number` option.
-   `-p`, `--maps-dir`: Overrides the `dataplaneapi.resources.maps_dir` option.
-   `-i`, `--show-system-info`: Overrides the `dataplaneapi.show_system_info` option.
-   `--disable-inotify`: Overrides the `dataplaneapi.disable_inotify` option.
-   `--pid-file`: Overrides the `dataplaneapi.pid_file` option.
-   `--debug-socket-path`: Overrides the `dataplaneapi.debug_socket_path` option.
-   `--uid`: Overrides the `dataplaneapi.uid` option.
-   `--gid`: Overrides the `dataplaneapi.gid` option.
-   `--userlist-file`: Overrides the `dataplaneapi.userlist.userlist_file` option.
-   `--status-cmd`: Overrides the `haproxy.reload.status_cmd` option.
-   `--service`: Overrides the `haproxy.reload.service_name` option.
-   `--reload-retention`: Overrides the `haproxy.reload.reload_retention` option.
-   `--reload-strategy`: Overrides the `haproxy.reload.reload_strategy` option.
-   `--validate-cmd`: Overrides the `haproxy.reload.validate_cmd` option.
-   `--backups-dir`: Overrides the `dataplaneapi.transaction.backups_dir` option.
-   `--max-open-transactions`: Overrides the `dataplaneapi.transaction.max_open_transactions` option.
-   `--ssl-certs-dir`: Overrides the `dataplaneapi.resources.ssl_certs_dir` option.
-   `--general-storage-dir`: Overrides the `dataplaneapi.resources.general_storage_dir` option.
-   `--dataplane-storage-dir`: Overrides the `dataplaneapi.resources.dataplane_storage_dir` option.
-   `--update-map-files`: Overrides the `dataplaneapi.resources.update_map_files` option.
-   `--update-map-files-period`: Overrides the `dataplaneapi.resources.update_map_files_period` option.
-   `--spoe-dir`: Overrides the `dataplaneapi.resources.spoe_dir` option.
-   `--spoe-transaction-dir`: Overrides the `dataplaneapi.resources.spoe_transaction_dir` option.
-   `--api-address`: Overrides the `dataplaneapi.advertised.api_address` option.
-   `--api-port`: Overrides the `dataplaneapi.advertised.api_port` option.
-   `--syslog-address`: Overrides the `log.syslog.syslog_address` option.
-   `--syslog-protocol`: Overrides the `log.syslog.syslog_protocol` option.
-   `--syslog-tag`: Overrides the `log.syslog.syslog_tag` option.
-   `--syslog-level`: Overrides the `log.syslog.syslog_level` option.
-   `--syslog-facility`: Overrides the `log.syslog.syslog_facility` option.
-   `--log-to`: Overrides the `log.log_to` option.
-   `--log-file`: Overrides the `log.log_file` option.
-   `--log-level`: Overrides the `log.log_level` option.
-   `--log-format`: Overrides the `log.log_format` option.
-   `--apache-common-log-format`: Overrides the `log.apache_common_log_format` option.
-   `--master-worker-mode`: Overrides the `haproxy.master_worker_mode` option.
-   `--delayed-start-max`: Overrides the `haproxy.delayed_start_max` option.
-   `--delayed-start-tick`: Overrides the `haproxy.delayed_start_tick` option.
-   `--fid`: Overrides the `haproxy.fid` option.

## Example Configuration

```yaml
config_version: 2
name: haproxy-dataplaneapi
dataplaneapi:
  host: 0.0.0.0
  port: 5555
  scheme:
  - http
  users:
    name: admin
    passowrd: adminpwd
    insecure: true
  transaction:
    transaction_dir: /etc/haproxy/transactions
    backups_number: 10
    backups_dir: /etc/haproxy/backups
    max_open_transactions: 10
  resources:
    maps_dir: /etc/haproxy/maps
    ssl_certs_dir: /etc/haproxy/ssl
    general_storage_dir: /etc/haproxy/general
    spoe_dir: /etc/haproxy/spoe
    dataplane_storage_dir: /etc/haproxy/dataplane
haproxy:
  config_file: /etc/haproxy/haproxy.cfg
  haproxy_bin: /usr/sbin/haproxy
  master_worker_mode: true
  master_runtime: /var/run/haproxy/master.sock
  reload:
    reload_delay: 5
    service_name: haproxy
    reload_strategy: systemd
log_targets:
- log_to: file
  log_file: /var/log/haproxy-dataplaneapi.log
  log_level: info
  log_types:
  - access
  - app
```
