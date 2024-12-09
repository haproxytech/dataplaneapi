// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	cnparents "github.com/haproxytech/client-native/v6/configuration/parents"
)

func operations(childType string) []string {
	switch childType {
	case cnparents.ServerChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
		}
	case cnparents.HTTPAfterResponseRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.HTTPCheckChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.HTTPErrorRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.HTTPRequestRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.HTTPResponseRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.TCPCheckChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.TCPRequestRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.TCPResponseRuleChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.QUICInitialRuleType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.ACLChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.BindChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
		}
	case cnparents.FilterChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	case cnparents.LogTargetChildType:
		return []string{
			"Create",
			"Get",
			"GetAll",
			"Delete",
			"Replace",
			"ReplaceAll",
		}
	}
	return []string{}
}
