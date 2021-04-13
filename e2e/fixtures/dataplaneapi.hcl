name = "famous_condor"

dataplaneapi {
  host = "0.0.0.0"
  port = 8080

  userlist {
    userlist-file = "/etc/haproxy/userlist.cfg"
  }

  resources {
    maps-dir      = "/etc/haproxy/maps"
    ssl-certs-dir = "/etc/haproxy/ssl"
    spoe-dir      = "/etc/haproxy/spoe"
  }

}

log {
  log-to = "file"
  log-file = "/var/log/dataplaneapi.log"
  log-level = "debug"
}

haproxy {
  config-file = "/etc/haproxy/haproxy.cfg"
  haproxy-bin = "/usr/local/sbin/haproxy"

  reload {
    reload-cmd  = "kill -SIGUSR2 1"
    restart-cmd = "kill -SIGUSR2 1"
  }
}
