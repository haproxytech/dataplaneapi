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

package parents

type Parent struct {
	PathParentType string
	ParentType     string
}

func Parents(childType string) []Parent {
	switch childType {
	case ServerChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "peers", ParentType: "Peer"},
			{PathParentType: "rings", ParentType: "Ring"},
		}
	case HTTPAfterResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case HTTPCheckChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
			{PathParentType: "healthchecks", ParentType: "Healthcheck"},
		}
	case HTTPErrorRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case HTTPRequestRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case HTTPResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case TCPCheckChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
			{PathParentType: "healthchecks", ParentType: "Healthcheck"},
		}
	case TCPRequestRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case TCPResponseRuleChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case QUICInitialRuleType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case ACLChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "fcgi_apps", ParentType: "FCGIApp"},
			{PathParentType: "defaults", ParentType: "Defaults"},
		}
	case BindChildType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "log_forwards", ParentType: "LogForward"},
			{PathParentType: "peers", ParentType: "Peer"},
		}
	case FilterChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
		}
	case LogTargetChildType:
		return []Parent{
			{PathParentType: "backends", ParentType: "Backend"},
			{PathParentType: "frontends", ParentType: "Frontend"},
			{PathParentType: "defaults", ParentType: "Defaults"},
			{PathParentType: "peers", ParentType: "Peer"},
			{PathParentType: "log_forwards", ParentType: "LogForward"},
		}
	case SSLFrontUseChildType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend"},
		}
	case ForceBeSwitchChildType:
		return []Parent{
			{PathParentType: "frontends", ParentType: "Frontend"},
		}
	}
	return nil
}
