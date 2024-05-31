#!/usr/bin/env bash
#
# Copyright 2019 HAProxy Technologies
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
set -eo pipefail

export BASE_PATH="/v3"

HAPROXY_VERSION=${HAPROXY_VERSION:-2.9}
DOCKER_BASE_IMAGE="${DOCKER_BASE_IMAGE:-haproxytech/haproxy-debian}:${HAPROXY_VERSION}"
DOCKER_CONTAINER_NAME="dataplaneapi-e2e"
export DOCKER_CONTAINER_NAME

ROOT_DIR=$(git rev-parse --show-toplevel)
export E2E_PORT=${E2E_PORT:-8042}
export E2E_DIR=${ROOT_DIR}/e2e
export LOCAL_IP_ADDRESS=${LOCAL_IP_ADDRESS:-127.0.0.1}

source "${E2E_DIR}/libs/cleanup.bash"

if ! docker version > /dev/null 2>&1; then
  echo '>>> Docker is not installed: cannot proceed for e2e test suite'
fi

if ! docker inspect "${DOCKER_BASE_IMAGE}" > /dev/null 2>&1; then
  echo '>>> Downloading base Docker image'
  docker pull "${DOCKER_BASE_IMAGE}"
fi

if [ ! -z $PREWIPE ] && [ "$PREWIPE" == "y" ]; then
   cleanup ${DOCKER_CONTAINER_NAME}
fi

# Custom configuration to run tests with the master socket.
IFS='.' read -ra version_parts <<< "$HAPROXY_VERSION"
major="${version_parts[0]}"
minor="${version_parts[1]}"

if [[ "$major" -eq "2"  &&  "$minor" -ge "7" || "$major" -gt "2" ]] ; then
  HAPROXY_FLAGS="-W -db -S /var/lib/haproxy/master -f /usr/local/etc/haproxy/haproxy.cfg"
  VARIANT="-master-socket"
fi


if [ ! -z $(docker ps -q -f name=${DOCKER_CONTAINER_NAME}) ]; then
    echo ">>> Skipping provisioning the e2e environment, ${DOCKER_CONTAINER_NAME} already present"
else
    echo '>>> Provisioning the e2e environment'
    docker run \
      --rm \
      --detach \
      --name ${DOCKER_CONTAINER_NAME} \
      --publish "${E2E_PORT}":8080 \
      --security-opt seccomp=unconfined \
      "${DOCKER_BASE_IMAGE}" $HAPROXY_FLAGS > /dev/null 2>&1
    docker cp "${ROOT_DIR}/build/dataplaneapi" ${DOCKER_CONTAINER_NAME}:/usr/local/bin/dataplaneapi
    docker cp "${E2E_DIR}/fixtures/dataplaneapi${VARIANT}.yaml" ${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplaneapi.yaml
    docker cp "${E2E_DIR}/fixtures/haproxy.cfg" ${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg
    docker cp "${E2E_DIR}/fixtures/userlist.cfg" ${DOCKER_CONTAINER_NAME}:/etc/haproxy/userlist.cfg
    docker exec -d ${DOCKER_CONTAINER_NAME} sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
fi

echo '>>> Waiting dataplane API to be up and running'
count=1
DATAPLANE_USER=$(grep insecure-password ${E2E_DIR}/fixtures/userlist.cfg | awk '{print $2}')
DATAPLANE_PASS=$(grep insecure-password ${E2E_DIR}/fixtures/userlist.cfg | awk '{print $4}')
until curl -s "${DATAPLANE_USER}:${DATAPLANE_PASS}@${LOCAL_IP_ADDRESS}":"${E2E_PORT}${BASE_PATH}/specification" 2>&1 1>/dev/null; do
    sleep 1;
    ((count++))
    if [ $count -eq 10 ]; then
        echo ">>> timeout waiting for dataplaneapi to start"
        exit 1
    fi
done

# deferring Docker container removal
# shellcheck disable=SC1090

if [ ! -z $SKIP_CLEANUP ] && [ "$SKIP_CLEANUP" == "y" ]; then
    echo ">>> Container will be left running: ${DOCKER_CONTAINER_NAME}"
else
    trap 'cleanup ${DOCKER_CONTAINER_NAME}' EXIT
fi

echo '>>> Starting test suite'
if [ ! -z $TESTNAME ]; then
    bats -t "${E2E_DIR}"/tests/${TESTNAME}
elif [ ! -z $TESTPART ]; then
    set +e
    echo $TESTPART | grep -q -e "[[:digit:]]/[[:digit:]]"
    if [ $? != 0 ]; then
        echo "invalid TESTPART argument: ${TESTPART}"
        exit 1
    fi
    set -e
    PARALLEL_RUNS=$(echo $TESTPART | cut -d\/ -f 2)
    THIS_RUN=$(($(echo $TESTPART | cut -d\/ -f 1) - 1))

    declare -a SELECTED_TESTS

    echo ">>> Selected partial run via TESTPART variable, running ${TESTPART} of tests:"
    ALL_TESTS=($(ls "${E2E_DIR}"/tests))
    ALL_TESTS_COUNT=${#ALL_TESTS[@]}
    for TESTNR in $(seq 0 $(( $ALL_TESTS_COUNT - 1 )) ); do
        if [ $(($TESTNR % $PARALLEL_RUNS)) == $THIS_RUN ]; then
            echo ">>> -> ${ALL_TESTS[$TESTNR]}"
            SELECTED_TESTS+=("${E2E_DIR}"/tests/${ALL_TESTS[$TESTNR]})
        fi
    done
    bats -t ${SELECTED_TESTS[*]}
else
    bats -t "${E2E_DIR}"/tests/*
fi
