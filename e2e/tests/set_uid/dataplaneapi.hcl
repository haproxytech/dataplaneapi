name = "famous_condor"

dataplaneapi {
  host = "0.0.0.0"
  port = 8080
  uid = 1500

  userlist {
    userlist_file = "/etc/haproxy/userlist.cfg"
  }

  resources {
    maps_dir                  = "/home/testuiduser/maps"
    ssl_certs_dir             = "/home/testuiduser/ssl"
    general_storage_dir       = "/home/testuiduser/general"
    spoe_dir                  = "/home/testuiduser/spoe"
    spoe_transaction_dir      = "/home/testuiduser/spoe-td"
  }
}

haproxy {
  config_file = "/home/testuiduser/haproxy.cfg"
  haproxy_bin = "/usr/local/sbin/haproxy"

  reload {
    reload_cmd  = "kill -SIGUSR2 1"
    restart_cmd = "kill -SIGUSR2 1"
  }
}
