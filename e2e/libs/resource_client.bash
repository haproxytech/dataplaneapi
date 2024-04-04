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

function resource_put() {
  local endpoint;  endpoint="$1"; shift
  local data;      data="$1"; shift
	local qs_params; qs_params="$1"
   	get_version
	run dpa_curl PUT "$endpoint?$qs_params&version=${VERSION}" "$data"
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
	dpa_curl_status_body '$output'
}
