module github.com/haproxytech/dataplaneapi

go 1.26

require (
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/aws/aws-sdk-go-v2 v1.42.1
	github.com/aws/aws-sdk-go-v2/config v1.32.30
	github.com/aws/aws-sdk-go-v2/credentials v1.19.29
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.69.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.316.1
	github.com/docker/go-units v0.5.0
	github.com/dustinkirkland/golang-petname v0.0.0-20260215035315-f0c533e9ce9b
	github.com/fsnotify/fsnotify v1.10.1
	github.com/getkin/kin-openapi v0.142.0
	github.com/go-openapi/errors v0.22.8
	github.com/go-openapi/loads v0.24.0
	github.com/go-openapi/runtime v0.32.5
	github.com/go-openapi/spec v0.22.6
	github.com/go-openapi/strfmt v0.27.0
	github.com/go-openapi/swag v0.27.0
	github.com/go-openapi/swag/cmdutils v0.27.0
	github.com/go-openapi/validate v0.26.0
	github.com/google/renameio v1.0.1
	github.com/google/uuid v1.6.0
	github.com/haproxytech/client-native/v5 v5.1.26
	github.com/jessevdk/go-flags v1.6.1
	github.com/json-iterator/go v1.1.12
	github.com/kr/pretty v0.3.1
	github.com/lestrrat-go/apache-logformat v0.0.0-20210106032603-24d066f940f8
	github.com/maruel/panicparse/v2 v2.5.0
	github.com/nathanaelle/syslog5424/v2 v2.0.5
	github.com/rs/cors v1.11.1
	github.com/rubyist/circuitbreaker v2.2.1+incompatible
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.11.1
	golang.org/x/net v0.57.0
	golang.org/x/sys v0.47.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.30 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.32.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.37.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.44.1 // indirect
	github.com/aws/smithy-go v1.27.4 // indirect
	github.com/cenk/backoff v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.25.3 // indirect
	github.com/go-openapi/jsonpointer v1.0.0 // indirect
	github.com/go-openapi/jsonreference v1.0.0 // indirect
	github.com/go-openapi/runtime/server-middleware v0.30.0 // indirect
	github.com/go-openapi/swag/conv v0.27.0 // indirect
	github.com/go-openapi/swag/fileutils v0.27.0 // indirect
	github.com/go-openapi/swag/jsonname v0.27.0 // indirect
	github.com/go-openapi/swag/jsonutils v0.27.0 // indirect
	github.com/go-openapi/swag/loading v0.27.0 // indirect
	github.com/go-openapi/swag/mangling v0.27.0 // indirect
	github.com/go-openapi/swag/netutils v0.27.0 // indirect
	github.com/go-openapi/swag/stringutils v0.27.0 // indirect
	github.com/go-openapi/swag/typeutils v0.27.0 // indirect
	github.com/go-openapi/swag/yamlutils v0.27.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gofrs/flock v0.13.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/haproxytech/go-logger v1.1.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lestrrat-go/strftime v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/oasdiff/yaml v0.1.1 // indirect
	github.com/oasdiff/yaml3 v0.0.14 // indirect
	github.com/oklog/ulid/v2 v2.1.1 // indirect
	github.com/peterbourgon/g2s v0.0.0-20170223122336-d4e7ad98afea // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.15.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
