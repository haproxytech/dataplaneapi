# ![HAProxy](../assets/images/haproxy-weblogo-210x49.png "HAProxy")

## HAProxy Data Plane API configuration file

documentation can be seen at [doc](doc/README.md)

example can be seen [yaml](examples/example-dataplaneapi.yaml)

full examples of configuration also can be seen at [yaml](examples/example-full.yaml)

# IMPORTANT information regarding migration from 2.x to 3.0 version

Some fields in dataplane configuration file are *deprecated* starting from version 3.0.

They are now moved to Data Plane API internal storage (`dataplane_storage_dir`): [read here](../storage/README.md)



Those fields are the one that were not really dataplane configuration attributes but *dynamic* data (cluster, users...). They are moved to a dataplane internal storage. Hence the dataplane configuration file is not any more updated with states and data and only contains configuration.

- `dataplaneapi.user` (**only cluster mode users, whose name are starting with `dpapi-c` are migrated**)
- `cluster`
- `service_discovery`
- `status`

Those data are moved to dataplane internal storage.



## Migration from Dataplane API 2.x to 3.0 behavior

Data (`user`, `cluster`, `service_discovery`, `status`) that were in Dataplane API configuration file are automatically created in:
- `<dataplane_storage_dir>/service_discovery.json`
- `<dataplane_storage_dir>/cluster.json`


## General behavior regarding the deprecated section

If some of those data are manually updated in Dataplane configuration file (**this should not be done, fields are deprecated**), or if it's the first start of a *3.0* dataplane API:
- A warning log is issued (search for logs having `"[CFG DEPRECATED]`) with `[SKIP]` or `[MIGRATE]` and the category:
  -  `[User]` or
  - `[Cluster]` or
  - `[Consul]` or
  - `[AWS Region]` or
  - `[Status]`

Below an an example for cluster users:

```
time="2024-07-22T09:12:09+02:00" level=warning msg="[CFG DEPRECATED] [SKIP] [User] [dpapi-c-Abr8s1V]: already migrated. Old location [/home/helene/go/src/gitlab.int.haproxy.com/dataplaneapi-ee/.test/etc/dataplaneapi-cluster.yml] New location [/home/helene/go/src/gitlab.int.haproxy.com/dataplaneapi-ee/.test/storage/dataplane/users.json]. Use only new location"
```
```
time="2024-07-22T09:12:09+02:00" level=warning msg="[CFG DEPRECATED] [MIGRATE] [User] [dpapi-c-8Mk2Z5UK]: migrating. Old location [/home/helene/go/src/gitlab.int.haproxy.com/dataplaneapi-ee/.test/etc/dataplaneapi-cluster.yml] New location [/home/helene/go/src/gitlab.int.haproxy.com/dataplaneapi-ee/.test/storage/dataplane/users.json]. Use only new location"
```

- **After migration, data are removed from the dataplane configuration file.**
- **If a data is added again in the dataplane configuration file, and has already been migrated to the internal storage, then the value from the configuration file is ignored. Only the value from the storage is kept.**



**IMPORTANT NOTE: only cluster mode users, whose name are starting with `dpapi-` are migrated**



## Cluster Mode: deprecation

The dataplane api configuration field:
- `mode` (values = `single` or `cluster`)

is now deprecated and this value is removed from dataplane api configuration file after migration.

The way `cluster` vs `single` is now handled is as following:

| Mode |     <dataplane_storage_dir>/cluster.json content    |
|-----------|-------------------|
| Single| `cluster` attribute is empty|
| Cluster| `cluster` attribute is not empty |
