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

# NOTE: in order to use this haproxy.cfg must be created in data test folder

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
  local haproxy_cfg_file="${BATS_TEST_DIRNAME}/data/haproxy_*.cfg"
  local haproxy_file_version=""
  local copy_haproxy_file=true

  if ls $haproxy_cfg_file 1> /dev/null 2>&1; then
    haproxy_file_version=$(echo $haproxy_cfg_file | sed 's/.*_\([0-9]\+\.[0-9]\+\)\.cfg/\1/')
    major_file_version=$(echo $haproxy_file_version | cut -d '.' -f 1)
    minor_file_version=$(echo $haproxy_file_version | cut -d '.' -f 2)
    major_cfg_version=$(echo $HAPROXY_VERSION | cut -d '.' -f 1)
    minor_cfg_version=$(echo $HAPROXY_VERSION | cut -d '.' -f 2)
    if [[ -f "${BATS_TEST_DIRNAME}/data/haproxy_${haproxy_file_version}.cfg" ]] ; then
        if [[  $major_cfg_version -eq $major_file_version && $minor_cfg_version -ge $minor_file_version || $major_cfg_version -gt $major_file_version ]] ; then
            run docker cp "${BATS_TEST_DIRNAME}/data/haproxy_${haproxy_file_version}.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
            copy_haproxy_file=false
            assert_success
        fi
    fi
  fi

  if [[ "$copy_haproxy_file" = true ]] ; then
    if [ -f "${BATS_TEST_DIRNAME}/data/haproxy.cfg" ]; then
        run docker cp "${BATS_TEST_DIRNAME}/data/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
        assert_success
    else
        run docker cp "${E2E_DIR}/fixtures/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
        assert_success
    fi
  fi

  # replace the default dataplaneapi config file
  if [ -f "${BATS_TEST_DIRNAME}/dataplaneapi.yaml" ]; then
      run docker cp "${BATS_TEST_DIRNAME}/dataplaneapi.yaml" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplaneapi.yaml"
  else
      # TODO: we should respect $VARIANT here
      run docker cp "${E2E_DIR}/fixtures/dataplaneapi.yaml" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplaneapi.yaml"
  fi
  assert_success

  if [ -d "${BATS_TEST_DIRNAME}/data/container" ]; then
      run docker cp "${BATS_TEST_DIRNAME}/data/container/." "${DOCKER_CONTAINER_NAME}:/"
      assert_success
  fi

  run dpa_docker_exec 'kill -s 12 1'
  assert_success

  if [ -z "$DONT_RESTART_DPAPI" ]; then
    run dpa_docker_exec 'pkill -9 dataplaneapi'
    assert_success

    if [ -x "${BATS_TEST_DIRNAME}/custom_dataplane_launch.sh" ]; then
        run "${BATS_TEST_DIRNAME}/custom_dataplane_launch.sh"
    else
        run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
    fi
    assert_success
  fi

  local restart_retry_count=0

  until dpa_curl GET "/info"; do
      # 5 seconds wait
      if [[ restart_retry_count -eq 50 ]]; then
          echo -e "\nerror restarting: no response from dataplane on /info"
          echo -e "\nhaproxy stderr:"
          docker logs -n 42 ${DOCKER_CONTAINER_NAME}
          echo -e "\ndataplane log output:"
          docker exec ${DOCKER_CONTAINER_NAME} sh -c "tail -n 42 /var/log/dataplaneapi.log"
          exit 1
      fi

      sleep 0.1
      restart_retry_count=$((restart_retry_count+1))
  done
}

# teardown returns original configuration to dataplane
teardown() {
  run docker cp "${E2E_DIR}/fixtures/haproxy.cfg" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/haproxy.cfg"
  assert_success

  run docker cp "${E2E_DIR}/fixtures/dataplaneapi.yaml" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplaneapi.yaml"
  assert_success

  run dpa_docker_exec 'kill -s 12 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
  assert_success

  local restart_retry_count=0

  until dpa_curl GET "/info"; do
      # 5 seconds wait
      if [[ restart_retry_count -eq 50 ]]; then
          echo -e "\nerror restarting: no response from dataplane on /info"
          echo -e "\nhaproxy stderr:"
          docker logs -n 42 ${DOCKER_CONTAINER_NAME}
          echo -e "\ndataplane log output:"
          docker exec ${DOCKER_CONTAINER_NAME} sh -c "tail -n 42 /var/log/dataplaneapi.log"
          exit 1
      fi

      sleep 0.1
      restart_retry_count=$((restart_retry_count+1))
  done
}
