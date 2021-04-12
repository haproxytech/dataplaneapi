config_version = 2
name =  "famous_condor"
mode =  "single"
status =  "null"

dataplaneapi {
  scheme =  ["http"]
  cleanup-timeout =  "10s"
  graceful-timeout =  "15s"
  max-header-size =  "1MiB"
  socket-path =  "/var/run/data-plane.sock"
  host =  "localhost"
  port =  "80"
  listen-limit =  "null"
  keep-alive =  "3m"
  read-timeout =  "30s"
  write-timeout =  "60s"
  show-system-info =  "false"
  disable-inotify =  "false"
  pid-file =  "/tmp/dataplane.pid"

  tls {
    tls-host =  "null"
    tls-port =  "null"
    tls-certificate =  "null"
    tls-key =  "null"
    tls-ca =  "null"
    tls-listen-limit =  "null"
    tls-keep-alive =  "null"
    tls-read-timeout =  "null"
    tls-write-timeout =  "null"
  }

  user "admin" {
    insecure =  "true"
    password =  "adminpwd"
  }

  userlist {
    userlist =  "controller"
    userlist-file =  "null"
  }

  transaction {
    transaction-dir =  "/tmp/haproxy"
    backups-number =  "0"
    backups-dir =  "/tmp/backups"
    max-open-transactions =  "20"
  }

  resources {
    maps-dir =  "/etc/haproxy/maps"
    ssl-certs-dir =  "/etc/haproxy/ssl"
    update-map-files =  "false"
    update-map-files-period =  "10"
    spoe-dir =  "/etc/haproxy/spoe"
    spoe-transaction-dir =  "/tmp/spoe-haproxy"
  }

  advertised {
    api-address =  "10.2.3.4"
    api-port =  "80"
  }
}

haproxy {
  config-file =  "/etc/haproxy/haproxy.cfg"
  haproxy-bin =  "haproxy"
  master-runtime =  "null"
  fid =  "null"
  master-worker-mode =  "false"

  reload {
    reload-delay =  "5"
    reload-cmd =  "null"
    restart-cmd =  "null"
    reload-retention =  "1"
    validate-cmd =  "null"
  }
}

cluster {
  cluster-tls-dir =  "null"
  id =  "null"
  bootstrap_key =  "null"
  active_bootstrap_key =  "null"
  token =  "null"
  url =  "null"
  port =  "80"
  api_base_path =  "null"
  api_nodes_path =  "null"
  api_register_path =  "null"
  storage-dir =  "null"
  cert-path =  "null"
  cert-fetched =  "null"
  name =  "null"
  description =  "null"
}

service_discovery {
  consuls =  "null"
  aws-regions =  "null"
}

log {
  log-to =  "stdout"
  log-file =  "/var/log/dataplaneapi/dataplaneapi.log"
  log-level =  "warning"
  log-format =  "text"
  apache-common-log-format =  "%h"

  syslog {
    syslog-address =  "null"
    syslog-protocol =  "tcp"
    syslog-tag =  "dataplaneapi"
    syslog-level =  "debug"
    syslog-facility =  "local0"
  }
}
