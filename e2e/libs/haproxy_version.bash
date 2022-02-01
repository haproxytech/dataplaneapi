#!/usr/bin/env bash
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

# haproxy_version curls the /runtime/info endpoint to get the current
# haproxy version and returns it as a string. haproxy must be setup with a
# valid stats socket for this to work, and it uses deprecated_auth_curl
# from version.bash script.
#
# Example:
# haproxy_version
# >>> 2.2.18-c6e1dfa
function haproxy_version() {
    read -r SC RES < <(deprecated_auth_curl GET "/v2/services/haproxy/runtime/info")
    V="$(get_json_path "$RES" ".[0].info.version")"
    debug $RES
    debug $V
    echo "$V"
}

# haproxy_version_ge returns 1 if the haproxy version is greater or
# equal to the given version, otherwise 0. Given version must be in format
# MAJOR.MINOR without the patch and commit numbers (e.g. 1.9, 2.2, etc.).
#
# Arguments:
# 1. the target MAJOR.MINOR version
#
# Example:
# haproxy_version_ge "2.1"
# >>> 1
function haproxy_version_ge() {
    target=$1; shift
    v=$(haproxy_version)
    major_minor="${v%.*}"
    numerical_v="${major_minor//.}"
    numerical_target="${target//.}"

    if [ "$numerical_v" -ge "$numerical_target" ]
    then
        return 0
    else
        return 1
    fi
}
