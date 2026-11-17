#!/usr/bin/env bash
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

# haproxy_worker_pid echoes the PID of the HAProxy worker, as reported by the
# runtime info endpoint. A reload replaces the worker, so an unchanged PID
# proves that no reload happened.
#
# The reload command of the e2e setup (kill -s 12 1) returns before HAProxy has
# finished reloading, so it first waits (up to 10s) until the stats socket is
# served by the current worker listed by the master; it fails if that does not
# happen. Old workers are ignored: they may linger (peers) but no longer serve
# the stats socket.
#
# Example:
# pid=$(haproxy_worker_pid)
function haproxy_worker_pid() {
  local i procs pid
  for i in $(seq 1 100); do
    procs=$(dpa_docker_exec 'echo "show proc" | socat /var/lib/haproxy/master - 2>/dev/null')
    pid=$(dpa_curl GET "/services/haproxy/runtime/info" | head -n -1 | jq -r '.info.pid // empty' 2>/dev/null)
    # the current worker is the first one listed under "# workers"
    if [[ -n "$pid" ]] && sed -n '/^# workers/{n;p;q}' <<< "$procs" | grep -Eq "^${pid}[[:space:]]+worker"; then
      echo "$pid"
      return 0
    fi
    sleep 0.1
  done
  echo "HAProxy did not settle: show proc: ${procs}; runtime info pid: ${pid}" >&2
  return 1
}
