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
load '../../libs/haproxy_config_setup'
load '../../libs/resource_client'
load '../../libs/version'
load '../../libs/haproxy_version'

load 'utils/_helpers'

@test "global: Return a global configuration" {
  resource_get "$_GLOBAL_BASE_PATH"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.chroot')" "/var/lib/haproxy"
  assert_equal "$(get_json_path "$BODY" '.user')" "haproxy"
  assert_equal "$(get_json_path "$BODY" '.group')" "haproxy"
  assert_equal "$(get_json_path "$BODY" '.performance_options.maxconn')" "4000"
  assert_equal "$(get_json_path "$BODY" '.pidfile')" "/var/run/haproxy.pid"
  assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].address')" "/var/lib/haproxy/stats"
  assert_equal "$(get_json_path "$BODY" '.runtime_apis[0].level')" "admin"
  if haproxy_version_ge "2.8"; then
      assert_equal "$(get_json_path "$BODY" '.tune_ssl_options.ocsp_update_min_delay')" "10"
      assert_equal "$(get_json_path "$BODY" '.tune_options.stick_counters')" "50"
      assert_equal "$(get_json_path "$BODY" '.tune_options.h2_be_initial_window_size')" "10"
      assert_equal "$(get_json_path "$BODY" '.tune_options.h2_be_max_concurrent_streams')" "11"
      assert_equal "$(get_json_path "$BODY" '.tune_options.h2_fe_initial_window_size')" "12"
      assert_equal "$(get_json_path "$BODY" '.tune_options.h2_fe_max_concurrent_streams')" "13"
  fi
}
