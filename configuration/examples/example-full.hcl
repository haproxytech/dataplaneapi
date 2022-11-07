config_version = 2
name = "famous_condor"
mode = "single"
status = "null"

dataplaneapi {
  scheme = ["http"]
  cleanup_timeout = "10s"
  graceful_timeout = "15s"
  max_header_size = "1MiB"
  socket_path = "/var/run/data-plane.sock"
  host = "localhost"
  port = 80
  listen_limit = 1024
  keep_alive = "3m"
  read_timeout = "30s"
  write_timeout = "60s"
  show_system_info = false
  disable_inotify = false
  pid_file = "/tmp/dataplane.pid"
  uid = 1000
  gid = 1000

  tls {
    tls_host = "null"
    tls_port = 6443
    tls_certificate = "null"
    tls_key = "null"
    tls_ca = "null"
    tls_listen_limit = 10
    tls_keep_alive = "1m"
    tls_read_timeout = "10s"
    tls_write_timeout = "10s"
  }

  user "admin" {
    insecure = true
    password = "adminpwd"
  }

  userlist {
    userlist = "controller"
    userlist_file = "null"
  }

  transaction {
    transaction_dir = "/tmp/haproxy"
    backups_number = 0
    backups_dir = "/tmp/backups"
    max_open_transactions = 20
  }

  resources {
    maps_dir = "/etc/haproxy/maps"
    ssl_certs_dir = "/etc/haproxy/ssl"
    update_map_files = false
    update_map_files_period = 10
    spoe_dir = "/etc/haproxy/spoe"
    spoe_transaction_dir = "/tmp/spoe-haproxy"
  }

  advertised {
    api_address = "10.2.3.4"
    api_port = 80
  }
}

haproxy {
  config_file = "/etc/haproxy/haproxy.cfg"
  haproxy_bin = "haproxy"
  master_runtime = "null"
  fid = "null"
  master_worker_mode = false

  reload {
    reload_delay = 5
    reload_cmd = "systemctl reload haproxy"
    restart_cmd = "systemctl restart haproxy"
    status_cmd = "systemctl status haproxy"
    service_name = "haproxy.service"
    reload_retention = 1
    reload_strategy = "custom"
    validate_cmd = "null"
  }
}

cluster {
  cluster_tls_dir = "null"
  id = "null"
  bootstrap_key = "null"
  active_bootstrap_key = "null"
  token = "null"
  url = "null"
  port = 80
  api_base_path = "null"
  api_nodes_path = "null"
  api_register_path = "null"
  storage_dir = "null"
  cert_path = "null"
  cert_fetched = false
  name = "null"
  description = "null"
}

service_discovery {
  consuls = []
  aws_regions = []
}

log_targets = [
  {
    log_to           = "stdout"
    log_level        = "debug"
    log_format       = "json"
    log_types = [
      "access",
      "app",
    ]
  },
  {
    log_to           = "file"
    log_file         = "/var/log/dataplanepi.log"
    log_level        = "info"
    log_format       = "text"
    log_types        = ["app"]
  },
  {
    log_to           = "syslog"
    log_level        = "info"
    syslog_address   = "127.0.0.1"
    syslog_protocol  = "tcp"
    syslog_tag       = "dataplaneapi"
    syslog_level     = "debug"
    syslog_facillity = "local0"
    log_types        = ["access"]
  },
]

# Deprecated: use log_targets instead
log {
  log_to = "stdout"
  log_file = "/var/log/dataplaneapi/dataplaneapi.log"
  log_level = "warning"
  log_format = "text"
  apache_common_log_format = "%h"

  syslog {
    syslog_address = "null"
    syslog_protocol = "tcp"
    syslog_tag = "dataplaneapi"
    syslog_level = "debug"
    syslog_facility = "local0"
  }
}
