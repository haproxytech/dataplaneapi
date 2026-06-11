# Upgrade notes

## Migration from go-swagger to chi / oapi-codegen

The API server was rewritten from go-swagger generated code to a chi router with
oapi-codegen generated handlers (OpenAPI 3). Endpoint paths, methods, parameters
and success responses are unchanged. The following client-visible behavior
changes were made intentionally:

- **Error body `code` now mirrors the HTTP status.** Validation failures that
  previously carried go-swagger's internal codes in the body (e.g.
  `{"code": 606, ...}` with HTTP 422) now return the HTTP status in the body as
  well (`{"code": 422, ...}`). Clients matching on the old 6xx body codes must
  switch to the HTTP status or the new body codes.

- **Raw configuration endpoints are text-only.**
  `POST /services/haproxy/configuration/raw` accepts only `text/plain` bodies;
  the previous (accidental) go-swagger behavior of also accepting a JSON-encoded
  string body is gone. Error responses from these endpoints are now
  `application/json` instead of `text/plain` (the JSON content is the same).

- **Unmatched paths require authentication.** Basic authentication is enforced
  in front of the router, so a request to a non-existent path without valid
  credentials returns `401` where the old server returned `404`. With valid
  credentials the response is still `404`.

- **404 / 405 / 406 / 415 / 422 response bodies keep the go-swagger JSON shape**
  (`{"code": <status>, "message": "..."}`); only the wording of some messages
  changed slightly.

- **SSL certificate DELETE can now return `204` on the default path.**
  `DELETE /services/haproxy/storage/ssl_certificates/{name}` without
  `skip_reload` or `force_reload` attempts to remove the certificate through the
  HAProxy runtime API first and returns `204 No Content` on success, falling
  back to `202 Accepted` with a `Reload-ID` header otherwise. The old server had
  a bug that made the runtime branch unreachable, so this path always returned
  `202`. Clients should treat both `204` and `202` as success.

- **Request bodies are capped at 1 GiB by default.** The old server read
  request bodies without any size limit. Bodies exceeding the limit are
  rejected with `413 Request Entity Too Large`. The limit is configurable via
  `--max-body-size` (or `max_body_size` in the configuration file); `0`
  disables it.

- **Replacing acme providers, log profiles, crt-stores, mailers sections and
  http-errors sections now triggers a reload.** `PUT` on
  `/services/haproxy/configuration/acme/{name}`, `.../log_profiles/{name}`,
  `.../crt_stores/{name}`, `.../mailers_section/{name}` and
  `.../http_errors_sections/{name}` previously wrote the change to the
  configuration file but never scheduled a reload, so the edit did not take
  effect until an unrelated change reloaded HAProxy, and the response was
  always `200`. These endpoints now follow the standard pattern: outside a
  transaction they return `202 Accepted` with a `Reload-ID` header (or `200`
  with `force_reload=true`), and inside a transaction they return `202`.

### Chi router additions for go-swagger compatibility

Where chi's defaults differ from the old server, the router carries explicit
code to keep the wire behavior unchanged:

- **JSON bodies for router-level errors.** chi answers unmatched paths and
  disallowed methods with plain-text bodies (`404 page not found`); custom
  handlers restore the JSON shape the go-swagger server returned
  (`{"code": 404, "message": "path ... was not found"}`). They are registered
  on both the root router and the mounted `/v3` sub-router.

- **`Allow` header on 405 responses.** RFC 9110 requires a 405 response to
  list the methods the resource supports, and the go-swagger server did so.
  chi does not expose the matched route's allowed methods to a custom 405
  handler, so the handler rebuilds the list by probing the router with each
  standard method for the requested path.

- **Trailing slashes are ignored.** Requests are routed as if a trailing
  slash were absent (chi's `StripSlashes` middleware), so `/v3/.../maps` and
  `/v3/.../maps/` hit the same handler.

- **`Reload-ID` header spelling.** The specification spells the reload header
  `Reload-ID`, while the server sends the canonical `Reload-Id` form — exactly
  as the old go-swagger server did, since Go's `http.Header` canonicalizes
  header names on write. Header names are case-insensitive, so clients should
  treat the two spellings as equal.
