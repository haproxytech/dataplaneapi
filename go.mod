module github.com/haproxytech/dataplaneapi

go 1.21

require (
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.5
	github.com/aws/aws-sdk-go-v2/credentials v1.16.16
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.36.7
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.144.0
	github.com/docker/go-units v0.5.0
	github.com/dustinkirkland/golang-petname v0.0.0-20231002161417-6a283f1aaaf2
	github.com/fsnotify/fsnotify v1.7.0
	github.com/getkin/kin-openapi v0.122.0
	github.com/go-openapi/errors v0.21.0
	github.com/go-openapi/loads v0.21.5
	github.com/go-openapi/runtime v0.26.2
	github.com/go-openapi/spec v0.20.14
	github.com/go-openapi/strfmt v0.22.0
	github.com/go-openapi/swag v0.22.8
	github.com/go-openapi/validate v0.22.6
	github.com/google/renameio v1.0.1
	github.com/google/uuid v1.5.0
	github.com/haproxytech/client-native/v4 v4.1.12
	github.com/haproxytech/config-parser/v4 v4.1.5
	github.com/hashicorp/consul/api v1.27.0
	github.com/hashicorp/hcl v1.0.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/json-iterator/go v1.1.12
	github.com/lestrrat-go/apache-logformat v0.0.0-20210106032603-24d066f940f8
	github.com/nathanaelle/syslog5424/v2 v2.0.5
	github.com/rodaine/hclencoder v0.0.1
	github.com/rs/cors v1.10.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.20.0
	golang.org/x/sys v0.16.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7 // indirect
	github.com/aws/smithy-go v1.19.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.22.2 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/haproxytech/go-logger v1.1.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.mongodb.org/mongo-driver v1.13.1 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
