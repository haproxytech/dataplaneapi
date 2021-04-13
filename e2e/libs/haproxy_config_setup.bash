#!/usr/bin/env bash
#
# Copyright 2020 HAProxy Technologies
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

# NOTE: in order to use this haproxy.cfg must be created in test folder

# setup puts configuration from test folder as the one active in dataplane
setup() {

  # allow for running just one test in a directory
  if [[ -n $TESTNUMBER ]] && [[ "$BATS_TEST_NUMBER" -ne $TESTNUMBER ]]; then
      skip
  fi

  if [[ -n "$TESTDESCRIPTION" ]] && [[ "$BATS_TEST_DESCRIPTION" != "$TESTDESCRIPTION" ]]; then
      skip
  fi

  # replace the default haproxy config file
  if [ -f "${BATS_TEST_DIRNAME}/haproxy.cfg" ]; then
      run docker cp "${BATS_TEST_DIRNAME}/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
  else
      run docker cp "${E2E_DIR}/fixtures/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
  fi
  assert_success

  if [ -d "${BATS_TEST_DIRNAME}/data/container" ]; then
      run docker cp "${BATS_TEST_DIRNAME}/data/container/." "${DOCKER_CONTAINER_NAME}:/"
      assert_success
  fi

  run dpa_docker_exec 'kill -SIGUSR2 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  if [ -x "${BATS_TEST_DIRNAME}/custom_dataplane_launch.sh" ]; then
      run "${BATS_TEST_DIRNAME}/custom_dataplane_launch.sh"
  else
      run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /usr/local/bin/dataplaneapi.hcl"
  fi
  assert_success
  until dpa_curl GET "/info"; do
      sleep 0.1
  done
}

# teardown returns original configuration to dataplane
teardown() {
  run docker cp "${E2E_DIR}/fixtures/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
  assert_success

  run dpa_docker_exec 'kill -SIGUSR2 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /usr/local/bin/dataplaneapi.hcl"
  assert_success
  until dpa_curl GET "/info"; do
      sleep 0.1
  done
}
