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
  read -r SC BODY < <(dataplaneapi_text_plain POST "/services/haproxy/configuration/raw?skip_version=true" "@${BATS_TEST_DIRNAME}/haproxy.cfg")
	#echo "Status Code: ${SC}"
	#echo "Body: ${BODY}"
	[ "${SC}" = 202 ]
  sleep 0.002
  read -r SC RES < <(dataplaneapi GET "/services/haproxy/configuration/global")
  V="$(RES=${RES} jq -n 'env.RES | fromjson | ._version')"
  while [ "$V" = "42" ]
  do
    sleep 0.001
    read -r SC RES < <(dataplaneapi GET "/services/haproxy/configuration/global")
    V="$(RES=${RES} jq -n 'env.RES | fromjson | ._version')"
  done
}

# teardown returns original configuration to dataplane
teardown() {
  read -r SC _ < <(dataplaneapi_text_plain POST "/services/haproxy/configuration/raw?skip_version=true" "@${E2E_DIR}/fixtures/haproxy.cfg")
  [ "${SC}" = 202 ]
  sleep 0.002
  read -r SC RES < <(dataplaneapi GET "/services/haproxy/configuration/global")
  V="$(RES=${RES} jq -n 'env.RES | fromjson | ._version')"
  while [ "$V" != "42" ]
  do
    sleep 0.001
    read -r SC RES < <(dataplaneapi GET "/services/haproxy/configuration/global")
    V="$(RES=${RES} jq -n 'env.RES | fromjson | ._version')"
  done
}
