#!/usr/bin/env bats
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

load '../../libs/dataplaneapi'
load '../../libs/haproxy_config_setup'

@test "x_issue_132: https://github.com/haproxytech/dataplaneapi/issues/132" {
    run dpa_curl POST "/services/haproxy/configuration/servers?backend=bug_132&force_reload=true&version=2" add_server.json
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201
}
