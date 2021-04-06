# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API configuration file

you can select between two formats, yaml & hcl

examples can be seen [hcl](examples/example-dataplaneapi.hcl) and [yaml](examples/example-dataplaneapi.yaml)
full examples of configuration also can be seen at [hcl](examples/example-full.hcl) and [yaml](examples/example-full.yaml)

## Converting between formats

you can convert from one format to another with

```bash
go run configuration/converter/converter.go original.cfgfile.x converted.x
```
