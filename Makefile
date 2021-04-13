DATAPLANEAPI_PATH?=${PWD}
GIT_REPO?=$(shell git config --get remote.origin.url)
GIT_HEAD_COMMIT=$(shell git rev-parse --short HEAD)
GIT_LAST_TAG=$(shell git describe --abbrev=0 --tags)
GIT_TAG_COMMIT=$(shell git rev-parse --short ${GIT_LAST_TAG})
GIT_MODIFIED1=$(shell git diff "${GIT_HEAD_COMMIT}" "${GIT_TAG_COMMIT}" --quiet || echo .dev)
GIT_MODIFIED2=$(shell git diff --quiet || echo .dirty)
GIT_MODIFIED=${GIT_MODIFIED1}${GIT_MODIFIED2}
BUILD_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')
CGO_ENABLED?=0
GOOS?=linux

all: update clean build

update:
	go get -v

clean:
	rm -rf ${DATAPLANEAPI_PATH}/build

.PHONY: build
build:
	mkdir -p ${DATAPLANEAPI_PATH}/build
	GOOS=$(GOOS) CGO_ENABLED=$(CGO_ENABLED) go build -trimpath -ldflags "-X main.GitRepo=${GIT_REPO} -X main.GitTag=${GIT_LAST_TAG} -X main.GitCommit=${GIT_HEAD_COMMIT} -X main.GitDirty=${GIT_MODIFIED} -X main.BuildTime=${BUILD_DATE}" -o ${DATAPLANEAPI_PATH}/build/dataplaneapi ${DATAPLANEAPI_PATH}/cmd/dataplaneapi/

.PHONY: e2e
e2e: build
	TESTNAME=$(TESTNAME) TESTNUMBER=$(TESTNUMBER) TESTDESCRIPTION="$(TESTDESCRIPTION)" SKIP_CLEANUP=$(SKIP_CLEANUP) PREWIPE=$(PREWIPE) HAPROXY_VERSION=$(HAPROXY_VERSION) ./e2e/run.bash
