# AWS Service Discovery

Data Plane API allows performing EC2 instances discovery, self-registering IP addresses as backend servers.

## Required tags

All instances must be tagged with the following tags:

- `HAProxy:Service:Name`: the service name will compose the HAProxy `backend` name.
- `HAProxy:Service:Port`: the default service port is listening to.

> The said tags are mandatory, otherwise, the instance will be ignored.

An additional tag is provided, in case of override for the single instance

- `HAProxy:Instance:Port`: allows to override the default Service port.

## Filtering

By default, all instances in the selected AWS region will be considered.

Selection of specific instances can be achieved using the `allowlist` functionality, specifying the desired EC2 filter to consider according to the [AWS documentation](https://docs.aws.amazon.com/cli/latest/reference/ec2/describe-instances.html#options).

```yaml
  aws_regions:
  - description: "Allowlist example"
    allowlist:
      - key: tag-key
        value: "Must:Have:This:Tag:Key"
    enabled: false
    id: 96b14c57-b011-42e5-8d01-b58feba07319
    name: john.doe
    region: "us-east-1"
    retrytimeout: 10
    serverslotsbase: 10
    serverslotsgrowthincrement: 10
    serverslotsgrowthtype: exponential
```

As `allowlist`, the `denylist` option allows to filter out specific instances matching the desired filters.

```yaml
  aws_regions:
  - description: "Denylist example"
    allowlist:
      - key: tag-key
        value: "Must:Have:This:Tag:Key"
    denylist:
      - key: "tag:Environment
        value: Development"
    enabled: false
    id: 96b14c57-b011-42e5-8d01-b58feba07319
    name: john.doe
    region: "us-east-1"
    retrytimeout: 10
    serverslotsbase: 10
    serverslotsgrowthincrement: 10
    serverslotsgrowthtype: exponential
```

## Authorization

Data Plane API needs the plain AWS credentials to interact with it.

```yaml
  aws_regions:
  - description: "Credentials example"
    secretaccesskey: "************************************soLl"
    accesskeyid: "****************L7GT"
    enabled: false
    id: 96b14c57-b011-42e5-8d01-b58feba07319
    name: john.doe
    region: "us-east-1"
    retrytimeout: 10
    serverslotsbase: 10
    serverslotsgrowthincrement: 10
    serverslotsgrowthtype: exponential
```

> In case of Data Plane API running in an EC2 with a IAM Role attached (as [`AmazonEC2ReadOnlyAccess`](https://console.aws.amazon.com/iam/home#/policies/arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess$serviceLevelSummary)), there's no need for additional credentials.

## Server options

Upon a Service discovery, Data Plane API will create the corresponding `backend` section using the following options:

- `ServerSlotsBase`: the minimum amount of `server` entries per `backend`
- `ServerSlotsGrowthIncrement`: the additional slots allocating for `server` in case of additional entries
- `ServerSlotsGrowthType`: the function type to implement in case of `server` slots growth

## Instances IP address

Using the JSON `ipv4_address` option you can specify which kind of IP address Data Plane API has to consider for the backend `server`.

Available values can be `private` (for the private one, reachable inside the AWS VPC) or `public`.

> If the instances doesn't have a public IPv4 address, and the service discovery configuration claims the `public` type, In case of `public` type, the EC2 will be ignored.

## Retry timeout

With the JSON `retry_timeout` option you can specify the interval of time elapsing between the reconciliation and the following.

Unit is expressed in __seconds__.

# Examples

## Creating a discovery on a selected AWS region

```json
// curl -XPOST "http://localhost:5555/v2/service_discovery/aws" -H 'content-type: application/json' -d @/path/to/payload.json
{
  "access_key_id": "****************L7GT",
  "enabled": true,
  "name": "my-service-discovery",
  "region": "us-east-1",
  "secret_access_key": "****************soLl",
  "ipv4_address": "private",
  "retry_timeout": 60
}
```

```yaml
  aws_regions:
  - enabled: true
    name: john.doe
    region: "us-east-1"
    secretaccesskey: "************************************soLl"
    ipv4address: private
    retrytimeout: 60
```

The resulting output will be the following, YMMV.

```
backend aws-us-east-1-my-service-discovery-my-service-name-8080
  server SRV_L17LT 172.31.68.158:8080 check weight 128
  server SRV_lsVqM 127.0.0.1:80 disabled weight 128
  server SRV_NTZyL 127.0.0.1:80 disabled weight 128
  server SRV_KMIFS 127.0.0.1:80 disabled weight 128
  server SRV_D2x28 127.0.0.1:80 disabled weight 128
  server SRV_MlgPJ 127.0.0.1:80 disabled weight 128
  server SRV_0SDZV 127.0.0.1:80 disabled weight 128
  server SRV_HnHJP 127.0.0.1:80 disabled weight 128
  server SRV_xMKi0 127.0.0.1:80 disabled weight 128
  server SRV_tWxu3 127.0.0.1:80 disabled weight 128
```

The `backend` name pattern is built with the following format:
`aws-<REGION>-<SERVICE_DISCOVERY_CONFIGURATION_NAME>-<SERVICE_NAME>-<SERVICE_PORT>`

## Pausing the discovery on a selected AWS region

```json
// curl -XPUT "http://localhost:5555/v2/service_discovery/aws/96b14c57-b011-42e5-8d01-b58feba07319" -H 'content-type: application/json' -d @/path/to/payload.json
{
  "access_key_id": "****************L7GT",
  "enabled": false,
  "name": "my-service-discovery",
  "region": "us-east-1",
  "secret_access_key": "****************soLl",
  "ipv4_address": "private",
  "retry_timeout": 60
}
```

```yaml
  aws_regions:
  - accesskeyid: "****************L7GT"
    enabled: false
    name: my-service-discovery
    region: "us-east-1"
    secretaccesskey: "****************soLl"
    ipv4address: private
    retrytimeout: 60
```

As a result of this action, Data Plane API will not update the discovered `backend` sections and their `server` entries: no data will be lost.

> Potentially, due to the spawn of newer EC2 instances or reboots with a change of the IPv4 address, data could be outdated.
