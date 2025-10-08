module github.com/haproxytech/dataplaneapi

go 1.24.0

require (
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5
	github.com/KimMachineGun/automemlimit v0.7.4
	github.com/aws/aws-sdk-go-v2 v1.39.2
	github.com/aws/aws-sdk-go-v2/config v1.31.12
	github.com/aws/aws-sdk-go-v2/credentials v1.18.16
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.56.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.254.1
	github.com/docker/go-units v0.5.0
	github.com/dustinkirkland/golang-petname v0.0.0-20240428194347-eebcea082ee0
	github.com/fsnotify/fsnotify v1.9.0
	github.com/getkin/kin-openapi v0.133.0
	github.com/go-openapi/errors v0.22.3
	github.com/go-openapi/loads v0.23.1
	github.com/go-openapi/runtime v0.29.0
	github.com/go-openapi/spec v0.22.0
	github.com/go-openapi/strfmt v0.24.0
	github.com/go-openapi/swag v0.25.1
	github.com/go-openapi/swag/cmdutils v0.25.1
	github.com/go-openapi/swag/mangling v0.25.1
	github.com/go-openapi/validate v0.25.0
	github.com/google/go-cmp v0.7.0
	github.com/google/renameio v1.0.1
	github.com/google/uuid v1.6.0
	github.com/haproxytech/client-native/v6 v6.1.8
	github.com/jessevdk/go-flags v1.6.1
	github.com/joho/godotenv v1.5.1
	github.com/json-iterator/go v1.1.12
	github.com/kr/pretty v0.3.1
	github.com/lestrrat-go/apache-logformat v0.0.0-20210106032603-24d066f940f8
	github.com/maruel/panicparse/v2 v2.5.0
	github.com/nathanaelle/syslog5424/v2 v2.0.5
	github.com/rs/cors v1.11.1
	github.com/rubyist/circuitbreaker v2.2.1+incompatible
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
	go.uber.org/automaxprocs v1.6.0
	golang.org/x/net v0.45.0
	golang.org/x/sys v0.36.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.6 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/cenk/backoff v2.2.1+incompatible // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/analysis v0.24.0 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.2 // indirect
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/fileutils v0.25.1 // indirect
	github.com/go-openapi/swag/jsonname v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/loading v0.25.1 // indirect
	github.com/go-openapi/swag/netutils v0.25.1 // indirect
	github.com/go-openapi/swag/stringutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/haproxytech/go-logger v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lestrrat-go/strftime v1.1.1 // indirect
	github.com/mailru/easyjson v0.9.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/peterbourgon/g2s v0.0.0-20170223122336-d4e7ad98afea // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/woodsbury/decimal128 v1.4.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.mongodb.org/mongo-driver v1.17.4 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
