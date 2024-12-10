_TRACES_PATH="/services/haproxy/configuration/traces"
_TRACE_ENTRIES_PATH="/services/haproxy/configuration/traces/entries"

function resource_delete_body() {
    local endpoint="$1" data="$2" qs_params="$3"
	get_version
	run dpa_curl DELETE "$endpoint?$qs_params&version=${VERSION}" "$data"
	assert_success
	dpa_curl_status_body '$output'
}
