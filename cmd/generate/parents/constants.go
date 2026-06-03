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

const (
	ServerChildType                = "server"
	HTTPAfterResponseRuleChildType = "http_after_response_rule"
	HTTPCheckChildType             = "http_check"
	HTTPErrorRuleChildType         = "http_error_rule"
	HTTPRequestRuleChildType       = "http_request_rule"
	HTTPResponseRuleChildType      = "http_response_rule"
	TCPCheckChildType              = "tcp_check"
	TCPRequestRuleChildType        = "tcp_request_rule"
	TCPResponseRuleChildType       = "tcp_response_rule"
	QUICInitialRuleType            = "quic_initial_rule"
	ACLChildType                   = "acl"
	BindChildType                  = "bind"
	FilterChildType                = "filter"
	LogTargetChildType             = "log_target"
	SSLFrontUseChildType           = "ssl_front_use"
)

type CnParentType string

const (
	BackendParentType    CnParentType = "backend"
	FrontendParentType   CnParentType = "frontend"
	DefaultsParentType   CnParentType = "defaults"
	LogForwardParentType CnParentType = "log_forward"
	PeerParentType       CnParentType = "peers"
	RingParentType       CnParentType = "ring"
	GlobalParentType     CnParentType = "global"
	FCGIAppParentType    CnParentType = "fcgi-app"
	ResolverParentType   CnParentType = "resolvers"
)
