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
load '../../libs/cleanup'

CONSUL_CONTAINER_NAME=consul-server
CONSUL_NETWORK=my-consul-net

setup() {
  # Stop dapi
  run dpa_docker_exec 'kill -s 12 1'
  assert_success

  run dpa_docker_exec 'pkill -9 dataplaneapi'
  assert_success

  ## consul
  run docker network create ${CONSUL_NETWORK}
  assert_success

  # start consul
  docker run -d --name ${CONSUL_CONTAINER_NAME} --net=${CONSUL_NETWORK} -p 8500:8500 -p 8600:8600/udp hashicorp/consul agent -server -ui -bootstrap-expect=1 -client=0.0.0.0
  assert_success

  CONSUL_ADDR=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONSUL_CONTAINER_NAME})
  debug "Local Consul IP: ${CONSUL_ADDR}"

  run docker network connect ${CONSUL_NETWORK} ${DOCKER_CONTAINER_NAME}
  assert_success

  # register services
  run docker cp "${BATS_TEST_DIRNAME}/data/sd_3_2.json" "${CONSUL_CONTAINER_NAME}:./"
  assert_success
  run docker exec -d ${CONSUL_CONTAINER_NAME} /bin/sh -c "consul services register ./sd_3_2.json"

  ####  dapi ####
  CONSUL_CONF=$(mktemp)
  RES=$(cat ${BATS_TEST_DIRNAME}/data/consul.json | sed "s/ADDRESS/${CONSUL_ADDR}/g")
  echo $RES > ${CONSUL_CONF}

  run dpa_docker_exec 'mkdir -p /var/lib/dataplaneapi/storage/service_discovery'
  assert_success

  #run docker cp "${BATS_TEST_DIRNAME}/data/consul.json" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplane/service_discovery/consul.json"
  run docker cp "${CONSUL_CONF}" "${DOCKER_CONTAINER_NAME}:/etc/haproxy/dataplane/service_discovery/consul.json"
  assert_success

  # Start dapi
  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
  assert_success

  # Wait for dapi to be ready with single user
  until dpa_curl GET "/info"; do
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

  cleanup ${CONSUL_CONTAINER_NAME}
  run docker rm ${CONSUL_CONTAINER_NAME}
  assert_success

  run docker network disconnect ${CONSUL_NETWORK} ${DOCKER_CONTAINER_NAME}

  run docker network rm ${CONSUL_NETWORK}
  assert_success

  run dpa_docker_exec 'rm /etc/haproxy/dataplane/service_discovery/consul.json'
  assert_success

  run docker exec -d ${DOCKER_CONTAINER_NAME} /bin/sh -c "CI_DATAPLANE_RELOAD_DELAY_OVERRIDE=1 dataplaneapi -f /etc/haproxy/dataplaneapi.yaml"
  assert_success

  until dpa_curl GET "/info"; do
      sleep 0.1
  done

}

@test "check consul" {

    backend_nb=0
    until [ "$backend_nb" -eq 2 ]; do
      backend_nb=$(dpa_docker_exec 'cat /etc/haproxy/haproxy.cfg' | grep backend | wc -l)
      sleep 1s
    done
    # check first SRV, if this one is found, all should be there
    srv1=$(dpa_docker_exec 'cat /etc/haproxy/haproxy.cfg' | grep 182.242.119.47 | wc -l)
    assert_equal "$srv1" 1

}
