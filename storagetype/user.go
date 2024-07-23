// Copyright 2021 HAProxy Technologies
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

package storagetype

import (
	"strings"
)

const (
	DapiClusterUserPrefix = "dpapi-c-"
)

type (
	Users []User
	User  struct {
		Name     string  `json:"name" yaml:"name"`
		Insecure *bool   `json:"insecure,omitempty" yaml:"insecure,omitempty"`
		Password *string `json:"password,omitempty" yaml:"password,omitempty"`
	}
)

func (u User) IsClusterUser() bool {
	return strings.HasPrefix(u.Name, DapiClusterUserPrefix)
}
