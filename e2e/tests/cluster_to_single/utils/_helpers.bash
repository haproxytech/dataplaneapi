#!/usr/bin/env bash
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

_ACL_BASE_PATH="/services/haproxy/configuration/acls"


function dpa_curl_clustermode() {
  verb=$1; shift
  endpoint=$1; shift
  data="@${BATS_TEST_DIRNAME}/$1";
  if [ -z "$1" ]; then
    data="/dev/null"
  fi
  curl -k -m 10 -s -H 'content-type: application/json' --user dpapi-c-vU9DIiJH:WLKrmlnvOxyHSAmtSi0xRue3 "-X${verb}" -w "\n%{http_code}" "-d${data}" "https://${LOCAL_IP_ADDRESS}:${E2E_PORT}${BASE_PATH}${endpoint}"
}
