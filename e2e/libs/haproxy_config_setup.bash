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
  assert_equal 0 1
  run dpa_curl_text_plain POST "/services/haproxy/configuration/raw?skip_version=true" "@${BATS_TEST_DIRNAME}/haproxy.cfg"
  assert_success

  dpa_curl_status_body '$output'
  assert_equal $SC 202
  sleep 0.002

  run dpa_curl GET "/services/haproxy/configuration/global"
  assert_success

  dpa_curl_status_body '$output'
  V="$(RES=${BODY} jq -n 'env.RES | fromjson | ._version')"
  while [ "$V" = "42" ]
  do
    sleep 0.001
    dpa_curl GET "/services/haproxy/configuration/global"
    assert_success

    dpa_curl_status_body '$output'
    V="$(RES=${BODY} jq -n 'env.RES | fromjson | ._version')"
  done
  exit 0
}

# teardown returns original configuration to dataplane
teardown() {
  dpa_curl_text_plain POST "/services/haproxy/configuration/raw?skip_version=true" "@${E2E_DIR}/fixtures/haproxy.cfg"
  assert_success

  assert_equal $SC 202

  sleep 0.002
  dpa_curl GET "/services/haproxy/configuration/global"
  assert_success

  dpa_curl_status_body '$output'
  V="$(RES=${BODY} jq -n 'env.RES | fromjson | ._version')"
  while [ "$V" != "42" ]
  do
    sleep 0.001
    dpa_curl GET "/services/haproxy/configuration/global"
    assert_success

    dpa_curl_status_body '$output'
    V="$(RES=${BODY} jq -n 'env.RES | fromjson | ._version')"
  done
}
