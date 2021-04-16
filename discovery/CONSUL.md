# CONSUL service discovery

Dataplane enables connecting to a existing consul service discovery instance and maps services into backends with the nodes as servers.

# Connecting to consul

In order to connect dataplane to consul an ACL token is required with appropriate policies/roles set. To register a consul services with
dataplane a request to the `service_discovery/consul` needs to be sent:

```json
curl -XPOST "http://localhost:5555/v2/service_discovery/consul" -H 'content-type: application/json' -d @/path/to/payload.json
{
  "token": "****************",
  "enabled": true,
  "name": "my-consul-service",
  "description": "Consul test system",
  "address": "127.0.0.1",
  "port" : 2222,
  "retry_timeout": 60
}
```
# Filtering

Using the `service_allowlist` and `service_denylist` fields a list of services to be tracked or ignored respectively can be provided. The `service_allowlist` option has precedence. In both cases, or if neither option is used, the `consul` service is ignored by dataplane.

# Backend example

The backend in the haproxy configuration associated with the consul service `my-service` will have the following format:

```
backend consul-backend-127.0.0.1-2222-my-service 
  server SRV_rfBd5 127.0.0.11:8088 weight 128
  server SRV_6ti2S 127.0.0.23:8081 weight 128
  server SRV_MtYvS 127.0.0.1:80 disabled weight 128
  server SRV_gD5xA 127.0.0.1:80 disabled weight 128
  server SRV_V0YU9 127.0.0.1:80 disabled weight 128
  server SRV_9zamp 127.0.0.1:80 disabled weight 128
  server SRV_ta7Z7 127.0.0.1:80 disabled weight 128
  server SRV_S575K 127.0.0.1:80 disabled weight 128
  server SRV_LkIZ9 127.0.0.1:80 disabled weight 128
  server SRV_PYkL1 127.0.0.1:80 disabled weight 128
```

The backend name will always have the following format: `consul-backend-<service-ip>-<service-port>-<service-name>`.
Each node associated with the services is represented by a server line with the following format:
`server SRV_<rand_id> <node-ip>:<node-port> weight 128`
In addition some additional disabled server lines will be added for future nodes to be used.

# Updating

After a consul service discovery is connected to dataplane its data will be polled every `retry_timeout` amount of seconds. Consuls
health endpoint i used to check if there have been changes for each individual service and if there were the server list for that
service gets updated. The changes to all services in a single update cycle are done in a single transaction. If the consul instance 
gets updated and the `enabled` parameter gets set to `false` data for that instance will not be pulled anymore but the existing configuration
for that service persists in its current state.

# Server slot scaling

The amount of server lines in a backend can be controlled using the `server_slots_base`, `server_slots_growth_type` and `server_slots_growth_increment` parameters. The `server_slots_base` option indicates how many servers are created on backend creation. Once the base slots used up the servers
expand based on the `server_slots_growth_type`: 
- `linear` : on each expansion the number of servers increases by the amount specified in `server_slots_growth_increment`
- `exponential` : the number of servers doubles on each server count increase
In case nodes get removed from the service the number of servers in the backend will decrease based on the specified growth type.

