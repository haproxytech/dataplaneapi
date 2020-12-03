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

# version return the current HAProxy configuration file version, useful to
# avoid keeping track of it at each POST/PUT call.
#
# Any argument is required.
#
# Example:
# version
# >>> 10
function version() {
  read -r SC RES < <(auth_curl GET "/v2/services/haproxy/configuration/global")
  V="$(RES=${RES} jq -n 'env.RES | fromjson | ._version')"
  echo "$V"
}
