#!/usr/bin/env bats
#
# Copyright 2026 HAProxy Technologies
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
load '../../libs/haproxy_worker'
load '../../libs/resource_client'
load '../../libs/version'

load 'utils/_helpers'

@test "runtime_apply: Commit a transaction changing weights in two backends without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  new_transaction
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv1?transaction_id=$TRANSACTION_ID" data/srv1_weight.json
  assert_equal "$SC" 202
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?transaction_id=$TRANSACTION_ID" data/srv2_weight.json
  assert_equal "$SC" 202
  dpa_request PUT "$_BACKEND_BASE_PATH/be_defsrv/servers/srv1?transaction_id=$TRANSACTION_ID" data/defsrv_srv1_weight.json
  assert_equal "$SC" 202

  dpa_request PUT "$_TRANSACTIONS_BASE_PATH/$TRANSACTION_ID"
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  assert_equal "$(runtime_server_field be_plain srv1 weight)" 11
  assert_equal "$(runtime_server_field be_plain srv2 weight)" 25
  assert_equal "$(runtime_server_field be_defsrv srv1 weight)" 8
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Commit a transaction adding and deleting servers without a reload" {
  local pid; pid=$(haproxy_worker_pid)

  new_transaction
  dpa_request POST "$_BACKEND_BASE_PATH/be_plain/servers?transaction_id=$TRANSACTION_ID" data/new_server.json
  assert_equal "$SC" 202
  dpa_request DELETE "$_BACKEND_BASE_PATH/be_plain/servers/srv1?transaction_id=$TRANSACTION_ID"
  assert_equal "$SC" 202

  dpa_request PUT "$_TRANSACTIONS_BASE_PATH/$TRANSACTION_ID"
  assert_equal "$SC" 200
  assert_equal "$RELOAD_ID" ""

  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv_new"
  assert_equal "$SC" 200
  assert_equal "$(get_json_path "$BODY" '.address')" "10.0.9.9"
  assert_equal "$(get_json_path "$BODY" '.port')" 8090
  resource_get "$_RUNTIME_BACKEND_BASE_PATH/be_plain/servers/srv1"
  assert_equal "$SC" 404
  assert_equal "$(haproxy_worker_pid)" "$pid"
}

@test "runtime_apply: Commit a transaction with a change not applicable at runtime reloads" {
  local pid; pid=$(haproxy_worker_pid)

  new_transaction
  dpa_request PUT "$_BACKEND_BASE_PATH/be_plain/servers/srv2?transaction_id=$TRANSACTION_ID" data/srv2_weight.json
  assert_equal "$SC" 202
  dpa_request POST "$_BACKEND_BASE_PATH/be_plain/acls/0?transaction_id=$TRANSACTION_ID" data/acl.json
  assert_equal "$SC" 202

  dpa_request PUT "$_TRANSACTIONS_BASE_PATH/$TRANSACTION_ID"
  assert_equal "$SC" 202
  assert [ -n "$RELOAD_ID" ]

  # the weight reaches HAProxy through the reload
  wait_reload_succeeded "$RELOAD_ID"
  wait_new_worker "$pid"
  assert_equal "$(runtime_server_field be_plain srv2 weight)" 25
}
