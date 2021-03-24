#!/usr/bin/env bash
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

# run_only run one test by its $BATS_TEST_NUMBER or $BATS_TEST_DESCRIPTION
# and skip running other tests.
# This function should be used as first statement inside BATS `setup` func
function run_only() {
    if [ -n $TESTNUMBER ] && [ "$BATS_TEST_NUMBER" -ne $TESTNUMBER ]; then
        skip
    fi

    if [ -n "$TESTDESCRIPTION" ] && [ "$BATS_TEST_DESCRIPTION" != "$TESTDESCRIPTION" ]; then
        skip
    fi
}
