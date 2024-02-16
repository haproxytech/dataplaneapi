#!/usr/bin/env bash
#
# Copyright 2023 HAProxy Technologies
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

_RAW_BASE_PATH="/services/haproxy/configuration/raw"

# Identical to resource_post() but with the text/plain content type.
function resource_post_text() {
  local endpoint="$1" data="@${BATS_TEST_DIRNAME}/$2" qs_params="$3"
  resource_get "/services/haproxy/configuration/version"
	version=${BODY}
  run curl -m 10 -s -H 'content-type: text/plain' --user dataplaneapi:mypassword \
    "-XPOST" -w "\n%{http_code}" "-d${data}" \
    "http://${LOCAL_IP_ADDRESS}:${E2E_PORT}${BASE_PATH}${endpoint}?$qs_params&version=$version"
	assert_success
	dpa_curl_status_body '$output'
}
