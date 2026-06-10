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


function resource_post() {
  local endpoint;  endpoint="$1"; shift
  local data;      data="$1"; shift
	local qs_params; qs_params="$1"
	get_version
	run dpa_curl POST "$endpoint?$qs_params&version=${VERSION}" "$data"
	assert_success
	dpa_curl_status_body '$output'
}

# Identical to resource_post() but posts the raw file body with the
# text/plain content type. --data-binary is used so newlines in the
# config file are preserved (curl -d strips them).
function resource_post_text() {
  local endpoint="$1" data="@${BATS_TEST_DIRNAME}/$2" qs_params="$3"
  resource_get "/services/haproxy/configuration/version"
	version=${BODY}
  run curl -m 10 -s -H 'content-type: text/plain' --user dataplaneapi:mypassword \
    "-XPOST" -w "\n%{http_code}" --data-binary "${data}" \
    "http://${LOCAL_IP_ADDRESS}:${E2E_PORT}${BASE_PATH}${endpoint}?$qs_params&version=$version"
	assert_success
	dpa_curl_status_body '$output'
}

function resource_post_no_data() {
	local endpoint="$1"; shift
	get_version
	run dpa_curl POST "$endpoint?version=${VERSION}"
	assert_success
	dpa_curl_status_body '$output'
}

function resource_put() {
  local endpoint;  endpoint="$1"; shift
  local data;      data="$1"; shift
	local qs_params; qs_params="$1"
   	get_version
	run dpa_curl PUT "$endpoint?$qs_params&version=${VERSION}" "$data"
	assert_success
	dpa_curl_status_body '$output'
}

function resource_put_no_data() {
	local endpoint="$1"; shift
	get_version
	run dpa_curl PUT "$endpoint?version=${VERSION}"
	assert_success
	dpa_curl_status_body '$output'
}

function resource_delete() {
	local endpoint;  endpoint="$1"; shift
	local qs_params; qs_params="$1"
	get_version
	run dpa_curl DELETE "$endpoint?$qs_params&version=${VERSION}"
	assert_success
	dpa_curl_status_body '$output'
}

function resource_get() {
  local endpoint;  endpoint="$1"; shift
  local qs_params; qs_params="$1"
  run dpa_curl GET "$endpoint?$qs_params"
	assert_success
	dpa_curl_status_body_safe '$output'
}
