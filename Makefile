DATAPLANEAPI_PATH?=$(shell pwd)
GIT_REPO?=$(shell git config --get remote.origin.url)
GIT_HEAD_COMMIT=$(shell git rev-parse --short HEAD)
GIT_LAST_TAG=$(shell git describe --abbrev=0 --tags)
GIT_TAG_COMMIT=$(shell git rev-parse --short ${GIT_LAST_TAG})
GIT_MODIFIED1=$(shell git diff "${GIT_HEAD_COMMIT}" "${GIT_TAG_COMMIT}" --quiet || echo .dev)
GIT_MODIFIED2=$(shell git diff --quiet || echo .dirty)
GIT_MODIFIED=${GIT_MODIFIED1}${GIT_MODIFIED2}
SWAGGER_VERSION=${shell curl -s https://raw.githubusercontent.com/haproxytech/client-native/master/Makefile | grep SWAGGER_VERSION -m 1 | awk -F"=" '{print $$2}'}
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
CGO_ENABLED?=0
GOLANGCI_LINT_VERSION=1.51.1

all: update clean build

update:
	go mod tidy

.PHONY: lint
lint:
	cd bin;GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION} sh lint-check.sh
	bin/golangci-lint run --timeout 5m --color always --max-issues-per-linter 0 --max-same-issues 0

.PHONY: gofumpt
gofumpt:
	go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

clean:
	rm -rf ${DATAPLANEAPI_PATH}/build

.PHONY: build
build:
	mkdir -p ${DATAPLANEAPI_PATH}/build
	CGO_ENABLED=$(CGO_ENABLED) go build -trimpath -ldflags "-X \"main.GitRepo=${GIT_REPO}\" -X main.GitTag=${GIT_LAST_TAG} -X main.GitCommit=${GIT_HEAD_COMMIT} -X main.GitDirty=${GIT_MODIFIED} -X main.BuildTime=${BUILD_DATE}" -o ${DATAPLANEAPI_PATH}/build/dataplaneapi ${DATAPLANEAPI_PATH}/cmd/dataplaneapi/

.PHONY: e2e
e2e: build
	TESTNAME=$(TESTNAME) TESTNUMBER=$(TESTNUMBER) TESTDESCRIPTION="$(TESTDESCRIPTION)" SKIP_CLEANUP=$(SKIP_CLEANUP) PREWIPE=$(PREWIPE) HAPROXY_VERSION=$(HAPROXY_VERSION) ./e2e/run.bash

.PHONY: generate
generate:
	cd generate/swagger;docker build \
		--build-arg SWAGGER_VERSION=${SWAGGER_VERSION} \
		--build-arg UID=$(shell id -u) \
		--build-arg GID=$(shell id -g) \
		-t dataplaneapi-swagger-gen .
	docker run --rm -v "$(shell pwd)":/data dataplaneapi-swagger-gen
	generate/post_swagger.sh

.PHONY: generate-native
generate-native:
	generate/swagger/script.sh
	generate/post_swagger.sh
