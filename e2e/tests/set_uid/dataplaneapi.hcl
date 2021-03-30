name = "famous_condor"

dataplaneapi {
  host = "0.0.0.0"
  port = 8080
  uid = 1500

  userlist {
    userlist-file = "/etc/haproxy/userlist.cfg"
  }

  resources {
    maps-dir                  = "/home/testuiduser/maps"
    ssl-certs-dir             = "/home/testuiduser/ssl"
    spoe-dir                  = "/home/testuiduser/spoe"
    spoe-transaction-dir      = "/home/testuiduser/spoe-td"
  }
}

haproxy {
  config-file = "/home/testuiduser/haproxy.cfg"
  haproxy-bin = "/usr/local/sbin/haproxy"

  reload {
    reload-cmd  = "kill -SIGUSR2 1"
    restart-cmd = "kill -SIGUSR2 1"
  }
}
