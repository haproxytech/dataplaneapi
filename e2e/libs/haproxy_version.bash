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

    IFS='.' read -ra version_parts <<< "$HAPROXY_VERSION"
    haproxy_major="${version_parts[0]}"
    haproxy_minor="${version_parts[1]}"

    IFS='.' read -ra version_parts_target <<< "$target"
    target_major="${version_parts_target[0]}"
    target_minor="${version_parts_target[1]}"

    if [[ "$haproxy_major" -eq "$target_major" && "$haproxy_minor" -ge "$target_minor" || "$haproxy_major" -gt "$target_major" ]] ; then
        return 0
    else
        return 1
    fi
}
