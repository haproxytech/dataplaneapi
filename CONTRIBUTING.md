# ![HAProxy](assets/images/haproxy-weblogo-210x49.png "HAProxy")

## Contributing guide

This document is a short guide for contributing to this project.

## API Specification - Development guide

### Dataplane API generation

Generation happens in two steps. Both run in CI (the `diff` job), so run both locally and commit the result whenever the specification changes:

```
make specification
make generate-native
```

`make specification` is the main step: it installs the pinned [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) into `bin/` and runs `cmd/generate/specification`, which combines the OpenAPI 3 spec files under [`specification/`](./specification) and regenerates the chi-based server code.

`make generate-native` regenerates only `embedded_spec.go` from the go-swagger spec in client-native. It needs the `swagger` binary on your `PATH`; `make generate` does the same inside Docker if you don't want to install it. Use `make generate-native` (not the Docker variant) when working against a local client-native checkout (a `replace github.com/haproxytech/client-native/v6 => ../client-native` directive in `go.mod`).

Generated files are marked with a `// Code generated ... DO NOT EDIT.` header and must not be edited by hand — change the specification and regenerate instead.

### Parent-typed resources

Some resources live under multiple parent types (e.g. `acl` belongs to backends, frontends, defaults, ...). Their parent expansion is handled as part of `make specification`: `cmd/generate/specification` template-expands each child's spec using the parent definitions in [cmd/generate/parents/parents.go](./cmd/generate/parents/parents.go).

To add or change the parents of a child resource, update the relevant `case` in `cmd/generate/parents/parents.go` and re-run `make specification`. There is no separate generation step or CI job for parents — it is covered by the single `make specification` run.


## Commit Messages and General Style

For commit messages and general style please follow the haproxy project's [CONTRIBUTING guide](https://github.com/haproxy/haproxy/blob/master/CONTRIBUTING) and use that where applicable.
