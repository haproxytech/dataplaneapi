# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

# Data Plane API internal storage (dataplane_storage_dir)


## Dataplane API internal storage location

Location of Dataplane API storage is managed by :
- dataplane API start argument `dataplane-storage-dir` (default value: `/etc/haproxy/dataplane`)


It can be overriden by a field in Daplane API configuration file:
- `dataplaneapi.resources.dataplane_storage_dir`: refer to  [example-full](examples/example-full.yaml)

## Files stored in Dataplane API internal storage location

| File name | content           | example | Comment |
|-----------|-------------------|---------|---------|
| cluster.json| Cluster configuration + Cluster users | [cluster.json](./examples/dapi-storage/cluster.json) | Prior to 3.0 was in Dapi config file: `cluster`, `dataplaneapi.user`, `status` sections|
| service_discovery/consul.json | Consul configuration| | Prior to 3.0 was in Dapi config file `service_discovery.consuls` section |
| service_discovery/aws.json | AWS configuration| | Prior to 3.0 was in Dapi config file `service_discovery.aws_regions` section |
