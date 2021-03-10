name = "famous_condor"

dataplaneapi {
  host = "0.0.0.0"
  port = 8080

  user "admin" {
    insecure = true
    password = "adminpwd"
  }

  resources {
    maps-dir      = "/etc/haproxy/maps"
    ssl-certs-dir = "/etc/haproxy/ssl"
    spoe-dir      = "/etc/haproxy/spoe"
  }
}

haproxy {
  config-file = "/etc/haproxy/haproxy.cfg"
  haproxy-bin = "/usr/local/sbin/haproxy"

  reload {
    reload-cmd  = "kill -SIGUSR2 1"
    restart-cmd = "kill -SIGUSR2 1"
  }
}
