# ![HAProxy](assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API

**Data Plane API** is a sidecar process that runs next to HAProxy and provides API endpoints for managing HAProxy. It requires HAProxy version 1.9.0 or higher.

## Building the Data Plane API

In order to build the Data Plane API you need go 1.14 installed on your system with go modules support enabled, and execute the following steps:

1\. Clone dataplaneapi repository

```
git clone https://github.com/haproxytech/dataplaneapi.git
```

2\. Run make build:

```
make build
```

3\. You can find the built binary in /build directory.

## Running the Data Plane API
Basic usage:

```
Usage:
  dataplaneapi [OPTIONS]

API for editing and managing haproxy instances

Application Options:
      --scheme=                                    the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=                           grace period for which to wait before killing idle connections (default: 10s)
      --graceful-timeout=                          grace period for which to wait before shutting down the server (default: 15s)
      --max-header-size=                           controls the maximum number of bytes the server will read parsing the request header's keys and values,
                                                   including the request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=                               the unix socket to listen on (default: /var/run/data-plane.sock)
      --host=                                      the IP to listen on (default: localhost) [$HOST]
      --port=                                      the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=                              limit the number of outstanding requests
      --keep-alive=                                sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing
                                                   laptop mid-download) (default: 3m)
      --read-timeout=                              maximum duration before timing out read of the request (default: 30s)
      --write-timeout=                             maximum duration before timing out write of the response (default: 60s)
      --tls-host=                                  the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=                                  the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=                           the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=                                   the private key to use for secure connections [$TLS_PRIVATE_KEY]
      --tls-ca=                                    the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=                          limit the number of outstanding requests
      --tls-keep-alive=                            sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing
                                                   laptop mid-download)
      --tls-read-timeout=                          maximum duration before timing out read of the request
      --tls-write-timeout=                         maximum duration before timing out write of the response

HAProxy options:
  -c, --config-file=                               Path to the haproxy configuration file (default: /etc/haproxy/haproxy.cfg)
  -u, --userlist=                                  Userlist in HAProxy configuration to use for API Basic Authentication (default: controller)
  -b, --haproxy-bin=                               Path to the haproxy binary file (default: haproxy)
  -d, --reload-delay=                              Minimum delay between two reloads (in s) (default: 5)
  -r, --reload-cmd=                                Reload command
  -s, --restart-cmd=                               Restart command
      --reload-retention=                          Reload retention in days, every older reload id will be deleted (default: 1)
  -t, --transaction-dir=                           Path to the transaction directory (default: /tmp/haproxy)
  -n, --backups-number=                            Number of backup configuration files you want to keep, stored in the config dir with version number suffix
                                                   (default: 0)
  -m, --master-runtime=                            Path to the master Runtime API socket
  -i, --show-system-info                           Show system info on info endpoint
  -f=                                              Path to the dataplane configuration file
      --userlist-file=                             Path to the dataplaneapi userlist file. By default userlist is read from HAProxy conf. When specified
                                                   userlist would be read from this file
      --fid=                                       Path to file that will dataplaneapi use to write its id (not a pid) that was given to him after joining a
                                                   cluster
  -p, --maps-dir=                                  Path to directory of map files managed by dataplane (default: /etc/haproxy/maps)
      --ssl-certs-dir=                             Path to SSL certificates directory (default: /etc/haproxy/ssl)
      --update-map-files                           Flag used for syncing map files with runtime maps values
      --update-map-files-period=                   Elapsed time in seconds between two maps syncing operations (default: 10)
      --cluster-tls-dir=                           Path where cluster tls certificates will be stored. Defaults to same directory as dataplane configuration file
      --spoe-dir=                                  Path to SPOE directory. (default: /etc/haproxy/spoe)
      --spoe-transaction-dir=                      Path to the SPOE transaction directory (default: /tmp/spoe-haproxy)
      --master-worker-mode                         Flag to enable helpers when running within HAProxy
      --max-open-transactions=                     Limit for active transaction in pending state (default: 20)

Logging options:
      --log-to=[stdout|file]                       Log target, can be stdout or file (default: stdout)
      --log-file=                                  Location of the log file (default: /var/log/dataplaneapi/dataplaneapi.log)
      --log-level=[trace|debug|info|warning|error] Logging level (default: warning)
      --log-format=[text|JSON]                     Logging format (default: text)

API options:
      --api-address=                               Advertised API address
      --api-port=                                  Advertised API port

Show version:
  -v, --version                                    Version and build information

Help Options:
  -h, --help                                       Show this help message

```

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

For more docs how to use the Data Plane API check our [documentation](https://www.haproxy.com/documentation/hapee/2-2r1/reference/dataplaneapi/)

Alternatively, dataplaneapi serves its own interactive documentation relevant for the current build on the `/v2/docs` uri. Just point your browser to the host/port dataplane was started with (i.e. `http://localhost:5555/v2/docs`)

## Contributing

If you wish to contribute to this project please check [Contributing Guide](CONTRIBUTING.md)
