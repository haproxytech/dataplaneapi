# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API

**Data Plane API** is a sidecar process that runs next to HAProxy and provides API endpoints for managing HAProxy.

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
```

This command generates some of the files in this project, which are marked with //DO NOT EDIT comments at the top of the files. These are not to be edited, as they are overwritten when specification is changed and the above-mentioned command is run.

## Dependencies

The project depends on the following internal projects:
- [models](http://github.com/haproxytech/models)
- [client-native](http://github.com/haproxytech/client-native)
- [config-parser](http://github.com/haproxytech/config-parser)

**client-native** project currently depends on an internal bash/awk project that manipulates HAProxy configuration: [lbctl](http://github.com/HAPEE/lbctl). This is planned to be fully replaced by native golang parser **config-parser** in the near future.

External dependecies:
- [interpose](https://github.com/carbocation/interpose)
- [martini](https://github.com/go-martini/martini)
- [inject](https://github.com/codegangsta/inject)
- [negroni](https://github.com/urfave/negroni)
- [recover](https://github.com/dre1080/recover)
- [stack](https://github.com/facebookgo/stack)
- [color](github.com/fatih/color)
- [go-openapi errors](https://github.com/go-openapi/errors)
- [go-openapi runtime](https://github.com/go-openapi/runtime)
- [go-openapi strfmt](https://github.com/go-openapi/strfmt)
- [govalidator](https://github.com/asaskevich/govalidator)
- [mgo](https://github.com/globalsign/mgo)
- [easyjson](https://github.com/mailru/easyjson)
- [mapstructure](https://github.com/mitchellh/mapstructure)
- [go-openapi swag](https://github.com/go-openapi/swag)

## Building

Following steps are required for building:

```
cd $GOPATH/src
git clone git@github.com:haproxy-controller/dataplaneapi.git
git clone git@github.com:haproxy-controller/models.git
git clone git@github.com:haproxy-controller/client-native.git
git clone git@github.com:haproxy-controller/config-parser.git

cd dataplaneapi
go get -v -insecure

cd cmd/dataplaneapi
go build
```

Currently you should also clone **lbctl** project (controller-dev branch) somewhere and install it:

```
git clone git@github.com:HAPEE/lbctl.git
cd lbctl
git checkout controller-dev
make install
cd /opt/lbctl/scripts/
chmod +x lbctl
```

## Running the Data Plane API

Basic usage:

```
./dataplaneapi --help
Usage:
  dataplaneapi [OPTIONS]

API for editing and managing HAPEE instances

Application Options:
      --scheme=                the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=       grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=       controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=           the unix socket to listen on (default: /var/run/dataplaneapi.sock)
      --host=                  the IP to listen on (default: localhost) [$HOST]
      --port=                  the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=          limit the number of outstanding requests
      --keep-alive=            sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=          maximum duration before timing out read of the request (default: 30s)
      --write-timeout=         maximum duration before timing out write of the response (default: 60s)
      --tls-host=              the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=              the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=       the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=               the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=                the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=      limit the number of outstanding requests
      --tls-keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=      maximum duration before timing out read of the request
      --tls-write-timeout=     maximum duration before timing out write of the response

HAProxy options:
  -c, --config-file=           Path to the haproxy configuration file (default: /etc/haproxy/haproxy.cfg)
  -g, --global-config-file=    Path to the haproxy global section configuration file (default: /etc/haproxy/haproxy-global.cfg)
  -u, --userlist=              Userlist in HAProxy configuration to use for API Basic Authentication (default: controller)
  -b, --haproxy-bin=           Path to the haproxy binary file (default: haproxy)
  -d, --reload-delay=          Minimum delay between two reloads (in s)
  -r, --reload-cmd=            Reload command
  -l, --lbctl-path=            Path to the lbctl script (default: lbctl)
  -t, --lbctl-transaction-dir= Path to the lbctl transaction directory (default: /tmp/lbctl)

Help Options:
  -h, --help                   Show this help message
```

## Example 

You can test it by simply running:

```
./dataplaneapi --port 5555 -b /usr/sbin/haproxy -c /etc/haproxy/haproxy.cfg  -g /etc/haproxy/global.cfg -d 5 -r "service reload haproxy" -u dataplaneapi -l /opt/lbctl/scripts/lbctl -t /tmp/lbctl
```

Test it out with curl, note that you need user/pass combination setup in HAProxy userlist in global.cfg (in above example: /etc/haproxy/global.cfg, userlist controller):

```
curl -u <user>:<pass> -H "Content-Type: application/json" "http://127.0.0.1:5555/v1/"
```