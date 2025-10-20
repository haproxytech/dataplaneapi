// Package exec provides a libdns provider which calls an external command.
//
// When called, the configured Command will be launched with the following
// environment variables:
//
// ACTION describes which action is the program expected to perform. Possible
// values are: get, append, set, delete.
//
// ZONE contains the DNS zone where the record is located.
//
// REC_NAME, REC_TTL, REC_TYPE, REC_DATA: respectively the record name, TTL in
// seconds, type (usually TXT), and value.
//
// When ran with the `get` action, the REC_* variables are not provided. Instead,
// the program is expected to scan the given zone and print the results as a JSON
// document on stdout. It must be parsable as a `[]libdns.RR`.
package exec
