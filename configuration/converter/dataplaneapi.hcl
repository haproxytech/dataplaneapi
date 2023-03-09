dataplaneapi {
  host = "0.0.0.0"
  port = 9555

  user "admin" {
    insecure = true
    password = "adminpwd"
  }

  resources {
    maps_dir            = "/etc/hapee-2.6/maps"
    ssl_certs_dir       = "/etc/hapee-2.6/ssl"
    spoe_dir            = "/etc/hapee-2.6/spoe"
    waf_dir             = "/etc/hapee-2.6/waf"
    general_storage_dir = "/etc/hapee-2.6/general"
  }

  advertised {
    api_address = "10.254.0.12"
    api_port    = 9555
  }
}

haproxy {
  config_file = "/etc/hapee-2.6/hapee-lb.cfg"
  haproxy_bin = "/opt/hapee-2.6/sbin/hapee-lb"

  reload {
    reload_delay = 5
    reload_strategy = "s6"
    service_name = "/var/run/s6/services/haproxy"
  }
}

keepalived {
  config_file     = "/etc/keepalived/default.conf"
  start_cmd       = "/usr/bin/systemctl start keepalived.service"
  reload_cmd      = "/usr/bin/systemctl reload keepalived.service"
  restart_cmd     = "/usr/bin/systemctl restart keepalived.service"
  stop_cmd        = "/usr/bin/systemctl stop keepalived.service"
  status_cmd      = "/usr/bin/systemctl status keepalived.service"
  validate_cmd    = "/usr/bin/keepalived --log-console --config-test=$KEEPALIVED_TRANSACTION_FILE"
  reload_delay    = 5
  transaction_dir = "/etc/keepalived/transactions"
  backups_dir     = "/etc/keepalived/backups"
  backups_number  = 10
}

log_targets = [
  {
    log_to     = "stdout"
    log_level  = "info"
    log_format = "text"
    log_types = [
      "access",
      "app",
    ]
  }
]
