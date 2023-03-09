name = "famous_condor"

dataplaneapi {
  host = "0.0.0.0"
  port = 8080
  scheme = ["http"]

  userlist {
    userlist_file = "/etc/haproxy/userlist.cfg"
  }

  resources {
    maps_dir      = "/etc/haproxy/maps"
    ssl_certs_dir = "/etc/haproxy/ssl"
    general_storage_dir = "/etc/haproxy/general"
    spoe_dir      = "/etc/haproxy/spoe"
  }

}

log {
  log_to = "file"
  log_file = "/var/log/dataplaneapi.log"
  log_level = "debug"
}

haproxy {
  config_file = "/etc/haproxy/haproxy.cfg"
  haproxy_bin = "/usr/local/sbin/haproxy"

  reload {
    reload_cmd  = "kill -SIGUSR2 1"
    restart_cmd = "kill -SIGUSR2 1"
  }
}
