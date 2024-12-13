// This file is safe to edit. Once it exists it will not be overwritten

// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package dataplaneapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

func TestConvOpenAPIV2ToV3(t *testing.T) {
	v2JSONString := string(SwaggerJSON)
	v2JSONString = strings.ReplaceAll(v2JSONString, "#/definitions", "#/components/schemas")
	curatedV2 := json.RawMessage([]byte(v2JSONString))

	var v2 openapi2.T
	err := v2.UnmarshalJSON(curatedV2)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = openapi2conv.ToV3(&v2)
	if err != nil {
		t.Error(err)
		return
	}
}
