# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API

**Data Plane API** is a sidecar process that runs next to HAProxy and provides API endpoints for managing HAProxy. It requires HAProxy version 1.9.0 or higher.

## API Specification

Data Plane API is built using [go-swagger](https://github.com/go-swagger/go-swagger) from the swagger spec found [here](http://github.com/haproxytech/haproxy-open-api-spec/blob/2.0/build/haproxy_spec.yaml) using the following command.

```
./swagger generate server -f ~/projects/haproxy-api/haproxy-open-api-spec/build/haproxy_spec.yaml \
    -A "Data Plane" \
    -t $GOPATH/src/github.com/haproxytech/ \
    --existing-models github.com/haproxytech/models \
    --exclude-main \
    --skip-models \
    -s dataplaneapi \
    --tags=Discovery \
    --tags=Information \
    --tags=Specification \
    --tags=Transactions \
    --tags=Sites \
    --tags=Stats \
    --tags=Global \
    --tags=Frontend \
    --tags=Backend \
    --tags=Bind \
    --tags=Server \
    --tags=Configuration \
    --tags=HTTPRequestRule \
    --tags=HTTPResponseRule \
    --tags=BackendSwitchingRule \
    --tags=ServerSwitchingRule \
    --tags=TCPResponseRule \
    --tags=TCPRequestRule \
    --tags=Filter \
    --tags=StickRule \
    --tags=LogTarget \
    --tags=Reloads \
    --tags=ACL
```

This command generates some of the files in this project, which are marked with //DO NOT EDIT comments at the top of the files. These are not to be edited, as they are overwritten when specification is changed and the above-mentioned command is run.

## Dependencies

The project depends on the following internal projects:
- [models](http://github.com/haproxytech/models)
- [client-native](http://github.com/haproxytech/client-native)
- [config-parser](http://github.com/haproxytech/config-parser)

External dependecies:
- [go-openapi](https://github.com/go-openapi)
- [go-units](https://github.com/docker/go-units)
- [go-flags](https://github.com/jessevdk/go-flags)
- [cors] (https://github.com/rs/cors)
- [graceful] (https://github.com/tylerb/graceful)
- [logrus] (https://github.com/sirupsen/logrus)
- [uuid] (https://github.com/google/uuid)
- [govalidator](https://github.com/asaskevich/govalidator)
- [mgo](https://github.com/globalsign/mgo)
- [easyjson](https://github.com/mailru/easyjson)
- [mapstructure](https://github.com/mitchellh/mapstructure)
- [crypt] (https://github.com/GehirnInc/crypt)

## Building

Following steps are required for building:

1\. Set your GOPATH variable

2\. Clone dataplaneapi repository into $GOPATH/src/github.com/haproxy/controller

```
cd $GOPATH/src/github.com/haproxy/controller
git clone http://github.com/haproxytech/dataplaneapi.git
```

3\. Run make build:

```
make build
```

4\. You can find the built binary in $GOPATH/bin directory.

## Running the Data Plane API

Basic usage:

```
./dataplaneapi --help
Usage:
  dataplaneapi [OPTIONS]

API for editing and managing HAPEE instances

Application Options:
      --scheme=                                    the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=                           grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=                           controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=                               the unix socket to listen on (default: /var/run/data-plane.sock)
      --host=                                      the IP to listen on (default: localhost) [$HOST]
      --port=                                      the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=                              limit the number of outstanding requests
      --keep-alive=                                sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=                              maximum duration before timing out read of the request (default: 30s)
      --write-timeout=                             maximum duration before timing out write of the response (default: 60s)
      --tls-host=                                  the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=                                  the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=                           the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=                                   the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=                                    the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=                          limit the number of outstanding requests
      --tls-keep-alive=                            sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=                          maximum duration before timing out read of the request
      --tls-write-timeout=                         maximum duration before timing out write of the response

HAProxy options:
  -c, --config-file=                               Path to the haproxy configuration file (default: /etc/haproxy/haproxy.cfg)
  -u, --userlist=                                  Userlist in HAProxy configuration to use for API Basic Authentication (default: controller)
  -b, --haproxy-bin=                               Path to the haproxy binary file (default: haproxy)
  -d, --reload-delay=                              Minimum delay between two reloads (in s)
  -r, --reload-cmd=                                Reload command
      --reload-retention=                          Reload retention in days, every older reload id will be deleted (default: 1)
  -t, --transaction-dir=                           Path to the transaction directory (default: /tmp/haproxy)

Logging options:
      --log-to=[stdout|file]                       Log target, can be stdout or file (default: stdout)
      --log-file=                                  Location of the log file (default: /var/log/dataplaneapi/dataplaneapi.log)
      --log-level=[trace|debug|info|warning|error] Logging level (default: warning)
      --log-format=[text|JSON]                     Logging format (default: text)

Help Options:
  -h, --help                                       Show this help message
```

## Example 

You can test it by simply running:

```
./dataplaneapi --port 5555 -b /usr/sbin/haproxy -c /etc/haproxy/haproxy.cfg  -d 5 -r "service haproxy reload" -u dataplaneapi -t /tmp/haproxy
```

Test it out with curl, note that you need user/pass combination setup in HAProxy userlist in global.cfg (in above example: /etc/haproxy/global.cfg, userlist controller):

```
curl -u <user>:<pass> -H "Content-Type: application/json" "http://127.0.0.1:5555/v1/"
```

If you are using secure passwords, supported algorithms are: md5, sha-256 and sha-512.