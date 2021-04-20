#!/usr/bin/env bats
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

load '../../libs/dataplaneapi'
load "../../libs/get_json_path"
load '../../libs/haproxy_config_setup'

@test "acl_runtime: Return ACL file entries list" {
    run dpa_curl GET "/services/haproxy/runtime/acl_file_entries?acl_id=0"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    assert_equal "$(get_json_path "${BODY}" " .[0].value" )" "/static"
    assert_equal "$(get_json_path "${BODY}" " .[1].value" )" "/images"
    assert_equal "$(get_json_path "${BODY}" " .[2].value" )" "/javascript"
    assert_equal "$(get_json_path "${BODY}" " .[3].value" )" "/stylesheets"
}


@test "acl_runtime: Return ACL file entries by their ID" {
    run dpa_curl GET "/services/haproxy/runtime/acl_file_entries?acl_id=0"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

    local LIST=$BODY

    for index in {0..3}; do
    	id="$(get_json_path "${LIST}" " .[${index}].id" )"
    	value="$(get_json_path "${LIST}" " .[${index}].value" )"

    	run dpa_curl GET "/services/haproxy/runtime/acl_file_entries/${id}?acl_id=0"
		assert_success

		dpa_curl_status_body '$output'
		assert_equal $SC 200

	 	assert_equal "$(get_json_path "${BODY}" " .value" )" $value
	done
}

@test "acl_runtime: Add an ACL file entry" {
    run dpa_curl POST "/services/haproxy/runtime/acl_file_entries?acl_id=0" "/data/post.json"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 201

	assert_equal "$(get_json_path "${BODY}" " .value" )" "/js"
}

@test "acl_runtime: Delete an ACL file entry by its ID" {
	# checking items and retrieving first ID
	run dpa_curl GET "/services/haproxy/runtime/acl_file_entries?acl_id=1"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

	assert_equal "$(get_json_path "${BODY}" " . | length" )" "5"

	local id
	id="$(get_json_path "${BODY}" ".[0].id")"

	# deleting the entry file by its ID
    run dpa_curl DELETE "/services/haproxy/runtime/acl_file_entries/${id}?acl_id=1"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 204

	# checking the file entry has been deleted counting the items
	run dpa_curl GET "/services/haproxy/runtime/acl_file_entries?acl_id=1"
    assert_success

    dpa_curl_status_body '$output'
    assert_equal $SC 200

	assert_equal "$(get_json_path "${BODY}" " . | length" )" "4"
}
