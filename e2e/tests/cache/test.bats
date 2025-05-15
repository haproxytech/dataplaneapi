#!/usr/bin/env bats
#
# Copyright 2022 HAProxy Technologies
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
load '../../libs/get_json_path'
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "cache: All Cache Tests" {
  resource_post "$_CACHE_BASE_PATH" "data/cache.json"
  assert_equal "$SC" 202

  resource_get "$_CACHE_BASE_PATH/cache_created"
  assert_equal "$SC" 200

  assert_equal "cache_created" "$(get_json_path "$BODY" ".name")"
  assert_equal 1 "$(get_json_path "$BODY" ".max_object_size")"
  assert_equal 1000 "$(get_json_path "$BODY" ".total_max_size")"

  resource_post "$_CACHE_BASE_PATH" "data/cache_same_name.json"
  assert_equal "$SC" 409

  resource_post "$_CACHE_BASE_PATH" "data/cache_unvalid.json"
  assert_equal "$SC" 400


  resource_delete "$_CACHE_BASE_PATH/test_cache2"
  assert_equal "$SC" 202

  resource_get "$_CACHE_BASE_PATH/test_cache2"
  assert_equal "$SC" 404

  resource_delete "$_CACHE_BASE_PATH/i_am_not_here"
  assert_equal "$SC" 404

  resource_get "$_CACHE_BASE_PATH/test_cache"
  assert_equal "$SC" 200
  assert_equal "test_cache" "$(get_json_path "$BODY" ".name")"
  assert_equal 60 "$(get_json_path "$BODY" ".max_age")"
  assert_equal 8 "$(get_json_path "$BODY" ".max_object_size")"
  assert_equal 1024 "$(get_json_path "$BODY" ".total_max_size")"


  resource_get "$_CACHE_BASE_PATH/i_am_not_here"
  assert_equal "$SC" 404

  resource_get "$_CACHE_BASE_PATH"
  assert_equal "$SC" 200
  assert_equal 2 "$(get_json_path "$BODY" ". | length")"

  resource_put "$_CACHE_BASE_PATH/test_cache" "data/cache_same_name.json"
  assert_equal "$SC" 202

  resource_get "$_CACHE_BASE_PATH/test_cache"
  assert_equal "$SC" 200
  assert_equal "test_cache" "$(get_json_path "$BODY" ".name")"
  assert_equal 1 "$(get_json_path "$BODY" ".max_object_size")"
  assert_equal 1000 "$(get_json_path "$BODY" ".total_max_size")"

  resource_put "$_CACHE_BASE_PATH/i_am_not_here" "data/cache_same_name.json"
  assert_equal "$SC" 404

  resource_put "$_CACHE_BASE_PATH/test_cache" "data/cache_unvalid.json"
  assert_equal "$SC" 400

}
