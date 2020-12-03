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

# get_json_path is a tiny wrapper to return content of a JSON from a variable.
# In case of missing
#
# Arguments:
# 1. the raw JSON
# 2. the JSON path according to jq syntax
#
# Example:
# get_json_path '[{"sites": [{"name": "foo"}, {"name": "bar"}]}]' '.[0].sites[1].name'
# >>> bar
function get_json_path() {
  local body path
  body=${1}
  path=${2}

  echo "$(JSON=${body} jq -rn "env.JSON | fromjson | ${path}")"
}
