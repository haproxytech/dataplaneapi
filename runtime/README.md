# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

# Runtime command socket for debugging purpose

## Goal
Dataplaneapi provides a command socket for debugging purpose.

The path for this socket can be defined :
- as dataplaneapi process argument (*--debug-socket-path*). Refer to [dataplaneapi README](../README.md)
- in the dataplaneapi configuration file. Refer to [configuration file](../configuration/README.md)

Note : if --debug-socket-path* is not set, the command socket will not run.

## Check the socket path at start up

You can check the socket path when starting dataplaneapi in the logs:
```
"level":"info","msg":"-- command socket Starting on /path/to/dataplane-debug.sock","time":"2023-09-13T11:44:08+02:00"}
``````

## Help

To display the available comands:

```
echo "help" | socat  - /path/to/dataplane-debug.sock
```

will output:
```
Dataplaneapi runtime commands:
           conf   show HAProxy configuration
     goroutines   display number of goroutines
          stack   output stack trace
        version   show dataplaneapi version
          pprof   pprof dumps
       dapiconf   show dataplaneapi configuration

type help <command> for more info
```


To get more info on a command, for example on the `conf` command:
```
echo "help conf" | socat  - /path/to/dataplane-debug.sock
```
will output:
```
Dataplaneapi runtime

command conf:
conf                                      show HAProxy current raw configuration
conf raw version [transactionID]          show HAProxy raw configuration for version (transactionID is optional default "" (for a transactionID, put '0' for version))
conf structured                           show HAProxy current structured configuration
conf structured version [transactionID]   show HAProxy structured configuration for version (transactionID is optional default "" (for a transactionID, put '0' for version))

```
## Available commands

-`conf` to show the HAProxy configuration (possibly for version/transaction):
   - `conf`: for current raw version
   - `conf raw`: for current raw version
   - `conf structured` for current structured version
Adding `version [transactionID]` to `conf raw` or `conf structured` allow to show the configuration for a given version or transaction.
Note that for a transactionID, version should be 0, for example : *conf raw 0  123-456-789*
- `dapiconf` for the dataplaneapi configuration
- `gouroutines` for the number of goroutines
- `stack` for the stack trace
- `version` for the dataplaneapi version
- `pprof` for pprof dumps


### More info for : `pprof`
To have some examples on how to use pprof, the following command will give you some good guidance:

```
echo "pprof examples" | socat  - /path/to/dataplane-debug.sock
```

### More info for: `dapiconf`
`dapiconf` dumps the *Configuration*. You can find an example of the output, when the dataplane is run with ` -f example-full.yaml`  [dapiconf output](./example/dapiconf_output.txt)
