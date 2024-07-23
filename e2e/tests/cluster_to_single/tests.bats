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
load '../../libs/debug'
load '../../libs/resource_client'
load 'utils/_helpers'


## This tests that we can switch from cluster mode to single mode
## using the cluster + user defined in cluster.json

setup() {
  # Stop dapi
  run dpa_docker_exec 'kill -s 12 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  # Set up everything for running in cluster mode
  run dpa_docker_exec 'mkdir -p /var/lib/dataplaneapi/storage/certs-cluster'
  assert_success

  run docker cp "${BATS_TEST_DIRNAME}/data/ca.crt" "${DOCKER_CONTAINER_NAME}:/var/lib/dataplaneapi/storage/certs-cluster/dataplane-famous_condor.crt"
  assert_success

  run docker cp "${BATS_TEST_DIRNAME}/data/ca.key" "${DOCKER_CONTAINER_NAME}:/var/lib/dataplaneapi/storage/certs-cluster/dataplane-famous_condor.key"
  assert_success

  run docker cp "${BATS_TEST_DIRNAME}/data/cluster.json" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplane/cluster.json"
  assert_success

  # Start dapi
  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
  assert_success

  # Wait for dapi to be ready with cluster user
  # using the cluster mode user in this call (see dpa_curl_clustermode in utils/_helpers.bash)
  until dpa_curl_clustermode GET "/info"; do
      sleep 0.1
  done

}

# teardown returns original configuration to dataplane
teardown() {
  # Stop dapi
  run dpa_docker_exec 'kill -s 12 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  run dpa_docker_exec 'rm /var/lib/dataplaneapi/storage/certs-cluster/dataplane-famous_condor.key'
  assert_success
  run dpa_docker_exec 'rm /var/lib/dataplaneapi/storage/certs-cluster/dataplane-famous_condor.crt'
  assert_success
  run dpa_docker_exec 'rm /etc/haproxy/dataplane/cluster.json'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
  assert_success

  until dpa_curl GET "/info"; do
      sleep 0.1
  done

}

@test "cluster to single: with cluster.json" {
    # DELETE cluster and check that we get back to single mode user
    # using the cluster mode user in this call (see dpa_curl_clustermode in utils/_helpers.bash)
    run dpa_curl_clustermode DELETE "/cluster"
    dpa_curl_status_body '$output'
    assert_equal "$SC" 204

    # Check that single mode user is now ok
    until dpa_curl GET "/info"; do
      sleep 0.1
    done
}
