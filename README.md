# ![HAProxy](assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API

[![Contributors](https://img.shields.io/github/contributors/haproxytech/dataplaneapi?color=purple)](CONTRIBUTING.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**Data Plane API** is a sidecar process that runs next to HAProxy and provides API endpoints for managing HAProxy. It requires HAProxy version 1.9.0 or higher.

## Building the Data Plane API

In order to build the Data Plane API you need Go installed on your system with go modules support enabled, and execute the following steps:

1\. Clone dataplaneapi repository

```
git clone https://github.com/haproxytech/dataplaneapi.git
```

2\. Run make build:

```
make build
```

3\. You can find the built binary in /build directory. TEST

## Running the Data Plane API
Basic usage:

```
Usage:
  dataplaneapi [OPTIONS]

API for editing and managing haproxy instances

Application Options:
      --scheme=                                       the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=                              grace period for which to wait before killing idle connections (default: 10s)
      --graceful-timeout=                             grace period for which to wait before shutting down the server (default: 15s)
      --max-header-size=                              controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request
                                                      body. (default: 1MiB)
      --socket-path=                                  the unix socket to listen on (default: /var/run/data-plane.sock)
      --host=                                         the IP to listen on (default: localhost) [$HOST]
      --port=                                         the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=                                 limit the number of outstanding requests
      --keep-alive=                                   sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=                                 maximum duration before timing out read of the request (default: 30s)
      --write-timeout=                                maximum duration before timing out write of the response (default: 60s)
      --tls-host=                                     the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=                                     the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=                              the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=                                      the private key to use for secure connections [$TLS_PRIVATE_KEY]
      --tls-ca=                                       the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=                             limit the number of outstanding requests
      --tls-keep-alive=                               sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=                             maximum duration before timing out read of the request
      --tls-write-timeout=                            maximum duration before timing out write of the response
      --uid                                           user id value to set on start
      --gid                                           group id value to set on start

HAProxy options:
  -c, --config-file=                                  Path to the haproxy configuration file (default: /etc/haproxy/haproxy.cfg)
  -u, --userlist=                                     Userlist in HAProxy configuration to use for API Basic Authentication (default: controller)
  -b, --haproxy-bin=                                  Path to the haproxy binary file (default: haproxy)
  -d, --reload-delay=                                 Minimum delay between two reloads (in s) (default: 5)
  -r, --reload-cmd=                                   Reload command
  -s, --restart-cmd=                                  Restart command
      --reload-retention=                             Reload retention in days, every older reload id will be deleted (default: 1)
  -t, --transaction-dir=                              Path to the transaction directory (default: /tmp/haproxy)
  -n, --backups-number=                               Number of backup configuration files you want to keep, stored in the config dir with version number suffix (default: 0)
      --backups-dir=                                  Path to directory in which to place backup files
  -m, --master-runtime=                               Path to the master Runtime API socket
  -i, --show-system-info                              Show system info on info endpoint
  -f=                                                 Path to the dataplane configuration file (default: /etc/haproxy/dataplaneapi.yaml)
      --userlist-file=                                Path to the dataplaneapi userlist file. By default userlist is read from HAProxy conf. When specified userlist would be read from this file
      --fid=                                          Path to file that will dataplaneapi use to write its id (not a pid) that was given to him after joining a cluster
  -p, --maps-dir=                                     Path to directory of map files managed by dataplane (default: /etc/haproxy/maps)
      --ssl-certs-dir=                                Path to SSL certificates directory (default: /etc/haproxy/ssl)
      --update-map-files                              Flag used for syncing map files with runtime maps values
      --update-map-files-period=                      Elapsed time in seconds between two maps syncing operations (default: 10)
      --cluster-tls-dir=                              Path where cluster tls certificates will be stored. Defaults to same directory as dataplane configuration file
      --spoe-dir=                                     Path to SPOE directory. (default: /etc/haproxy/spoe)
      --spoe-transaction-dir=                         Path to the SPOE transaction directory (default: /tmp/spoe-haproxy)
      --master-worker-mode                            Flag to enable helpers when running within HAProxy
      --max-open-transactions=                        Limit for active transaction in pending state (default: 20)
      --validate-cmd=                                 Executes a custom command to perform the HAProxy configuration check
      --disable-inotify                               Disables inotify watcher watcher for the configuration file
      --pid-file=                                     Path to file that will dataplaneapi use to write its pid
      --debug-socket-path=                            Unix socket path for the debugging command socket
Logging options:
      --log-to=[stdout|file|syslog]                   Log target, can be stdout, file, or syslog (default: stdout)
      --log-file=                                     Location of the log file (default: /var/log/dataplaneapi/dataplaneapi.log)
      --log-level=[trace|debug|info|warning|error]    Logging level (default: warning)
      --log-format=[text|JSON]                        Logging format (default: text)
      --apache-common-log-format=                     Apache Common Log Format to format the access log entries (default: %h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i" %{us}T)

Syslog options:
      --syslog-address=                               Syslog address (with port declaration in case of TCP type) where logs should be forwarded: accepting socket path in case of unix or unixgram
      --syslog-protocol=[tcp|tcp4|tcp6|unix|unixgram] Syslog server protocol (default: tcp)
      --syslog-tag=                                   String to tag the syslog messages (default: dataplaneapi)
      --syslog-level=                                 Define the required syslog messages level, allowed values: debug|info|notice|warning|error|critical|alert|emergency  (default: debug)
      --syslog-facility=                              Define the Syslog facility number, allowed values: kern|user|mail|daemon|auth|syslog|lpr|news|uucp|cron|authpriv|ftp|local0|local1|local2|local3|local4|local5|local6|local7
                                                      (default: local0)

Show version:
  -v, --version                                       Version and build information

Help Options:
  -h, --help                                          Show this help message
```

Beside those options, everything can be defined in side of configuration file. See [configuration file](configuration/README.md)

## Example

You can test it by simply running:

```
./dataplaneapi --port 5555 -b /usr/sbin/haproxy -c /etc/haproxy/haproxy.cfg  -d 5 -r "service haproxy reload" -s "service haproxy restart" -u dataplaneapi -t /tmp/haproxy
```

Dataplaneapi will require write permissions to the haproxy configuration file and the directories containing additional managed files (maps, ssl, spoe). The default locations can be overriden with command-line options.
Test it out with curl, note that you need user/pass combination setup in HAProxy userlist in haproxy configuration (in above example: /etc/haproxy/haproxy.cfg, userlist controller):

```
curl -u <user>:<pass> -H "Content-Type: application/json" "http://127.0.0.1:5555/v2/"
```

If you are using secure passwords, supported algorithms are: md5, sha-256 and sha-512.

## Using the Data Plane API

For more docs how to use the Data Plane API check our [documentation](https://www.haproxy.com/documentation/hapee/latest/api/data-plane-api/)

Alternatively, dataplaneapi serves its own interactive documentation relevant for the current build on the `/v2/docs` uri. Just point your browser to the host/port dataplane was started with (i.e. `http://localhost:5555/v2/docs`)

## Service Discovery

Check the documentation in the [README](./discovery/README.md).

## Command socket for debugging purpose

Check the documentation in the [README](./runtime/README.md).

## Contributing

If you wish to contribute to this project please check [Contributing Guide](CONTRIBUTING.md)
