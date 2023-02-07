#!/usr/bin/env bats
#
# Copyright 2021 HAProxy Technologies
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

load '../../libs/dataplaneapi'

setup() {
  run dpa_docker_exec 'kill -SIGUSR2 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  run dpa_docker_exec 'adduser -u 1500 testuiduser'
  #assert_success ignore error since we do not plan to insert password, user will be created

  run docker cp "${BATS_TEST_DIRNAME}/dataplaneapi.yaml" "${DOCKER_CONTAINER_NAME}:/home/testuiduser/dataplaneapi.yaml"
  assert_success

  run dpa_docker_exec 'chown testuiduser /home/testuiduser/dataplaneapi.yaml'
  assert_success

  run dpa_docker_exec 'cp /etc/haproxy/haproxy.cfg /home/testuiduser/haproxy.cfg'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /home/testuiduser/dataplaneapi.yaml"
  assert_success
  until dpa_curl GET "/info"; do
      sleep 0.1
  done
}

# teardown returns original configuration to dataplane
teardown() {
  run docker cp "${E2E_DIR}/fixtures/dataplaneapi.yaml" "${DOCKER_CONTAINER_NAME}:/usr/local/bin/dataplaneapi.yaml"
  assert_success

  run dpa_docker_exec 'kill -SIGUSR2 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /usr/local/bin/dataplaneapi.yaml"
  assert_success
  until dpa_curl GET "/info"; do
      sleep 0.1
  done
}

@test "set_uid: check if running user is testuiduser for dataplane" {
    PS=$(dpa_docker_exec "ps -eo ruser,user,comm | grep dataplaneapi")
    echo $PS
    [ "${PS}" = "testuidu testuidu dataplaneapi" ]
}
