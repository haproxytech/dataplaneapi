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
load '../../libs/haproxy_config_setup'
load '../../libs/version'

load 'utils/_helpers'

@test "transactions: Outdated transactions cannot be committed" {
	resource_post "$_TRANSACTIONS_BASE_PATH" ""
	assert_equal "$SC" 201
	local first_id; first_id=$(get_json_path "${BODY}" ".id")

	resource_post "$_TRANSACTIONS_BASE_PATH" ""
	assert_equal "$SC" 201
	local second_id; second_id=$(get_json_path "${BODY}" ".id")

	# commit first one, must succeed
	resource_put "$_TRANSACTIONS_BASE_PATH/$first_id" ""
 	assert_equal "$SC" 202

	# try to commit second one, must be outdated
	resource_put "$_TRANSACTIONS_BASE_PATH/$second_id" ""
	assert_equal "$SC" 406
}
