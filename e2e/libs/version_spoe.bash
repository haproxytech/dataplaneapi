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

# auth_curl is going to return the response status code along with the body
# these values can be easily read as following
#
# read -r SC BODY < <(auth_curl GET /v3/services/haproxy/runtime/info)
# echo "Status Code: ${SC}"
# echo "Body: ${BODY}"
#
# Arguments:
# - HTTP verb
# - original URL
# - HTTP POST data
function deprecated_auth_curl() {
  verb=$1; shift
  endpoint=$1; shift
  data=${1:-"/dev/null"}
  response=$(curl -m 10 -s -H 'content-type: application/json' --user dataplaneapi:mypassword "-X${verb}" -w "\n%{http_code}" "-d${data}" "http://${LOCAL_IP_ADDRESS}:${E2E_PORT}${endpoint}")
  status_code=$(tail -n1 <<< "$response")
  response=$(sed '$ d' <<< "$response")
  echo "$status_code $response"
}

function get_version() {
  resource_get "/services/haproxy/spoe/version" "spoe=$SPOE_FILE"
  eval VERSION="'$BODY'"
}
