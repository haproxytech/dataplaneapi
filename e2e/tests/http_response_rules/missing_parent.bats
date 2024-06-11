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
load '../../libs/get_json_path'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "http_response_rules: Fail creating a HTTP Response rule when frontend doesn't exist" {
  resource_post "$_RES_RULES_BASE_PATH/0" "data/post.json" "parent_type=frontend&parent_name=ghost&force_reload=true"
	assert_equal "$SC" 400
}

@test "http_response_rules: Fail creating a HTTP Response rule when backend doesn't exist" {
  resource_post "$_RES_RULES_BASE_PATH/0" "data/post.json" "parent_type=backend&parent_name=ghost&force_reload=true"
	assert_equal "$SC" 400
}
