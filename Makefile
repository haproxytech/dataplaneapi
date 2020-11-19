
ROOT = $(shell pwd)
INSTALL_PREFIX = /usr/local/haproxy-dataplaneapi


DATAPLANEAPI_PATH=${PWD}
GIT_REPO=$(shell git config --get remote.origin.url)
GIT_HEAD_COMMIT=$(shell git rev-parse --short HEAD)
GIT_LAST_TAG=$(shell git describe --abbrev=0 --tags)
GIT_TAG_COMMIT=$(shell git rev-parse --short ${GIT_LAST_TAG})
GIT_MODIFIED1=$(shell git diff "${GIT_HEAD_COMMIT}" "${GIT_TAG_COMMIT}" --quiet || echo .dev)
GIT_MODIFIED2=$(shell git diff --quiet || echo .dirty)
GIT_MODIFIED=${GIT_MODIFIED1}${GIT_MODIFIED2}
BUILD_DATE=$(shell date '+%Y-%m-%dT%H:%M:%S')

all: clean build

update:
	go get -v

clean:
	rm -rf ${DATAPLANEAPI_PATH}/build

.PHONY: build
build:
	mkdir -p ${DATAPLANEAPI_PATH}/build
	CGO_ENABLED=0 go build -gcflags "-N -l" -ldflags "-X main.GitRepo=${GIT_REPO} -X main.GitTag=${GIT_LAST_TAG} -X main.GitCommit=${GIT_HEAD_COMMIT} -X main.GitDirty=${GIT_MODIFIED} -X main.BuildTime=${BUILD_DATE}" -o ${DATAPLANEAPI_PATH}/build/dataplaneapi ${DATAPLANEAPI_PATH}/cmd/dataplaneapi/


install:
	mkdir -p $(DESTDIR)/$(INSTALL_PREFIX)/bin/ 
	cp -v $(ROOT)/build/dataplaneapi $(DESTDIR)/$(INSTALL_PREFIX)/bin/

clean:
	#
	#$(MAKE) -C bindata clean

package:
	# Building deb
	# 
	$(MAKE) all
	dpkg-buildpackage -b
	# 
	# all done, here is ALL your package(s)
	#
	ls -l ../haproxy*.deb
