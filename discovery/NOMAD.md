# Nomad service discovery

Dataplane enables connecting to an existing Nomad cluster and map Service Registrations into HAProxy backends.

# Connecting to Nomad

To connect Dataplane to an existing Nomad cluster, an HTTP request to `service_discovery/nomad` needs to be sent:

```json
curl -XPOST "http://localhost:5555/v2/service_discovery/nomad" -H 'content-type: application/json' -d @/path/to/payload.json
{
  "secret_id": "****************",
  "enabled": true,
  "name": "my-nomad-service",
  "description": "Nomad test system",
  "address": "http://127.0.0.1",
  "port" : 4646,
  "namespace": "default",
  "retry_timeout": 60
}
// Response:
{
    "address": "http://127.0.0.1",
    "allowlist": null,
    "denylist": null,
    "description": "Nomad test system",
    "enabled": true,
    "id": "931e5ed0-9ac1-434c-b026-4d02f5c2fadb",
    "name": "my-nomad-service",
    "namespace": "*",
    "port": 4646,
    "retry_timeout": 60,
    "secret_id": "****************",
    "server_slots_base": 10,
    "server_slots_growth_increment": 10,
    "server_slots_growth_type": "linear"
}
```

**NOTE**

- If there's no ACL on Nomad cluster then `secret_id` can be left empty.
- To discover services across all Nomad namespaces, `*` (wildcard) operator can be used.

Alternatively, the configuration can be set in a file either in `.hcl` or `.yaml` form:

```hcl
service_discovery {
  nomads = [
    {
      Address                    = "http://127.0.0.1"
      Description                = "Nomad test system"
      Enabled                    = true
      ID                         = "b40eb63b-2fb2-4996-b870-20f50ca173de"
      Name                       = "my-nomad-service"
      Namespace                  = "default"
      Port                       = 4646
      RetryTimeout               = 15
      SecretID                   = ""
      ServerSlotsBase            = 10
      ServerSlotsGrowthIncrement = 10
      ServerSlotsGrowthType      = "linear"
    },
  ]
}
```

```yaml
service_discovery:
  nomads:
  - address: http://127.0.0.1
    allowlist: []
    denylist: []
    description: Nomad test system
    enabled: true
    id: b40eb63b-2fb2-4996-b870-20f50ca173de
    name: my-nomad-service
    namespace: '*'
    port: 4646
    retrytimeout: 15
    secretid: ""
    serverslotsbase: 10
    serverslotsgrowthincrement: 10
    serverslotsgrowthtype: linear
```

## Filtering Services

By default, dataplane will attempt to fetch all services for a given namespace. Using the `allowlist` and `denylist` fields a list of services to be tracked or ignored respectively can be provided. The `allowlist` option has precedence. Both of these filtering options can be used simultaneously.

The service name is of the format: `<namespace>-<service-name>`.

To filter for `good-svc` registered inside `default` namespace the following example can be used:

```hcl
service_discovery {
  nomads = [
    {
      ...
      Allowlist = ["default-good-svc"]
      Denylist  = ["default-bad-svc"]
      ...
    },
  ]
}
```

## Discovering all services

Namespace filtering can be done by providing the `namespace` option. Only one namespace can be specified and only services from that 
namespaces will be tracked. Specifying `*` would return all services across all the authorized namespaces.

## Authorization

If ACL is enabled on Nomad cluster, then a valid ACL token should be configured with the `secret_id` field in Nomad service discovery configuration. A valid ACL policy should also be attached to this ACL token.

An example of an ACL policy which allows to discover services across all namespaces:

```hcl
namespace "*" {
  policy = "read-job"
}
```

## Updating

After a nomad service discovery is connected to dataplane it's data will be polled every `retry_timeout` amount of seconds. Nomad Services HTTP API 
is used to check if there have been changes for each individual service and if there were the server list for that
service gets updated. The changes to all services in a single update cycle are done in a single transaction. If the nomad instance 
gets updated and the `enabled` parameter gets set to `false` data for that instance will not be pulled anymore but the existing configuration
for that service persists in its current state.

## Server slot scaling

The amount of server lines in a backend can be controlled using the `server_slots_base`, `server_slots_growth_type` and `server_slots_growth_increment` parameters. The `server_slots_base` option indicates how many servers are created on backend creation. Once the base slots used up the servers
expand based on the `server_slots_growth_type`: 
- `linear` : on each expansion the number of servers increases by the amount specified in `server_slots_growth_increment`
- `exponential` : the number of servers doubles on each server count increase
In case nodes get removed from the service the number of servers in the backend will decrease based on the specified growth type.

# Examples

## Creating a discovery on a selected Nomad namespace for all services

The following example shows how to fetch all services inside `orange` namespace and reconcilie every 15s.

```hcl
service_discovery {
  nomads = [
    {
      Address                    = "http://127.0.0.1"
      Description                = "Nomad Service discovery for orange namespace"
      Enabled                    = true
      ID                         = "b40eb63b-2fb2-4996-b870-20f50ca173de"
      Name                       = "example-nomad-svc"
      Namespace                  = "orange"
      Port                       = 4646
      RetryTimeout               = 15
    },
  ]
}
```

The backend in the haproxy configuration associated with the nomad service `web-fruits` in `orange` namespace will have the following format:

```
backend nomad-backend-orange-web-fruits
  server SRV_fDxC4 192.168.29.76:31513 check weight 128
  server SRV_m12Um 192.168.29.76:26146 check weight 128
  server SRV_Yi8sO 192.168.29.76:31528 check weight 128
  server SRV_nQOhc 192.168.29.76:27266 check weight 128
  server SRV_II8XR 127.0.0.1:80 disabled weight 128
  server SRV_YV0co 127.0.0.1:80 disabled weight 128
  server SRV_A8zUR 127.0.0.1:80 disabled weight 128
  server SRV_khA3j 127.0.0.1:80 disabled weight 128
  server SRV_w1B3d 127.0.0.1:80 disabled weight 128
  server SRV_ZxeUk 127.0.0.1:80 disabled weight 128
```

The backend name will always have the following format: `nomad-backend-<namespace>-<service-name>`.

Each node associated with the services is represented by a server line with the following format: `server SRV_<rand_id> <node-ip>:<node-port> weight 128`

In addition some additional disabled server lines will be added for future nodes to be used.
