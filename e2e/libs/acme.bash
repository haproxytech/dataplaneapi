#!/usr/bin/env bash
#
# Copyright 2025 HAProxy Technologies
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This module will launch 2 extra containers needed for ACME tests:
# - pebble - the "official" ACME test server from LetsEncrypt
# - challserv - pebble's companion used to fake DNS challenges
#
# Those containers can be accessed via their respective variables:
# - ACME_DIR_CONTAINER_NAME=pebble
# - ACME_DNS_CONTAINER_NAME=challtestsrv

load '../../libs/haproxy_version'

# ACME is available starting with HAProxy 3.2,
# but became compatible with pebble in version 3.3.
ACME=false
if haproxy_version_ge 3.3; then
    ACME=true
    ACME_DIR_CONTAINER_NAME=pebble
    ACME_DNS_CONTAINER_NAME=challtestsrv
    PEBBLE_IMAGE="ghcr.io/letsencrypt/pebble:latest"
    CHALLTESTSRV_IMAGE="ghcr.io/letsencrypt/pebble-challtestsrv:latest"
fi

setup_file() {
    $ACME || return 0

    echo '>>> Setting up ACME challtestsrv' >&3
    # Start challtestsrv -- https://github.com/letsencrypt/pebble/blob/main/cmd/pebble-challtestsrv/README.md
    docker pull "$CHALLTESTSRV_IMAGE" >/dev/null
    docker run -d --name "$ACME_DNS_CONTAINER_NAME" --net "$DOCKER_NETWORK" \
        -p 5001:5001 -p 5002:5002 -p 5003:5003 -p 8053:8053 -p 8055:8055 \
        "$CHALLTESTSRV_IMAGE"

    echo '>>> Setting up ACME directory: pebble' >&3
    docker pull "$PEBBLE_IMAGE" >/dev/null
    docker run -d --name "$ACME_DIR_CONTAINER_NAME" --net "$DOCKER_NETWORK" \
        -p 14000:14000 -p 15000:15000  -v "${E2E_DIR}/fixtures/pebble:/mnt:ro" \
        -e PEBBLE_VA_NOSLEEP=1 -e PEBBLE_VA_ALWAYS_VALID=1 -e PEBBLE_WFE_NONCEREJECT=0 \
        "$PEBBLE_IMAGE" -config /mnt/pebble-config.json -strict -dnsserver $ACME_DNS_CONTAINER_NAME:8053
    sleep 1

    # Allow HAProxy's HTTPS client to connect to Pebble (the fake ACME server)
    docker cp "${E2E_DIR}/fixtures/pebble/pebble.minica.pem" ${DOCKER_CONTAINER_NAME}:/var/lib/haproxy/pebble.minica.pem
}

teardown_file() {
    $ACME || return 0

    for c in "$ACME_DIR_CONTAINER_NAME" "$ACME_DNS_CONTAINER_NAME"; do
        echo ">>> Shutting down ACME container $c" >&3
        mkdir -p "$E2E_DIR/logs"
        docker logs "$c" &> "$E2E_DIR/logs/$c.log"
        docker stop "$c"
        docker rm -f "$c"
    done
}
