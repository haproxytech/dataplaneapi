# Dataplane API E2E test suite

## Requirements

- [Docker](https://www.docker.com/)
- [BATS](https://github.com/sstephenson/bats)
- [cURL](https://github.com/curl/curl)
- [jq](https://github.com/stedolan/jq)

## Run the test suite

You can use the recipe `e2e` available in the root `Makefile`:
this will compile the latest build according to local changes, injecting
the compiled binary into a HAProxy container instance.

```bash
$ make e2e
# mkdir -p /path/to/haproxytech/dataplaneapi/build
# CGO_ENABLED=0 go build -ldflags REDACTED
# ./e2e/run.bash
# >>> Provisioning the e2e environment
# >>> Starting test suite
# 1..81
# REDACTED
# >>> Stopping dataplaneapi-e2e docker container
```

In any case, at the end of execution, a clean-up will occur, removing the
Docker container.

### Run specific tests
When adding new tests, you may not want to run all tests at once.
To run tests in specific directory, you can use the following command:
```bash
TESTNAME="dir_name" make e2e
```
### Run only one test in specific dir
To run only one test by its description:
```bash
TESTNAME="dir_name" TESTDESCRIPTION="Test description"  make e2e
```

Alternatively, you can run it by its BATS number ($BATS_TEST_NUMBER is used):
```bash
TESTNAME="dir_name" TESTNUMBER=2  make e2e
```

Prerequisite to run test like this is to use BATS convenient `setup` function
in tests and invoke `run_only` function inside of that function

### Parameters

#### Host port

The E2E container will map the host port `8042`: in case it's already
allocated, you can prepend the environment variable `E2E_PORT` to specify a
different one.

```bash
$ E2E_PORT=8081 make e2e
```

#### HAProxy version

By default, test suite is running against HAProxy 2.3 release: this can be
configured using the environment variable `HAPROXY_VERSION`.

```bash
$ HAPROXY_VERSION=2.2 make e2e
```

#### Base Docker image

The suite is running on `haproxytech/haproxy-alpine` Docker image: this
can be overridden using the environment variable `DOCKER_BASE_IMAGE`.

```bash
$ DOCKER_BASE_IMAGE=registry.tld/repository/image make e2e
```

## Debug a failed test

You can redirect output as any bash script using the redirection `&> 3`.

## Writing test cases

We require grouping tests per feature or API group, documenting the expected
behavior and result as normal code would.

Each test should be self-contained and without any external dependency or
pre-condition: if you need this, use the `setup` and `teardown` functions.

```bash
#!/usr/bin/env bats

setup() {
    # executed before each test
    echo "setup" >&3
}

teardown() {
    # executed after each test
    echo "teardown" >&3
}

@test "test_name" {
    [ true = true]
}
```

If you need some assets as a request payload, put these fixtures in the same
test folder in order to load it locally: try to avoid inline declaration.
`/path/to/post/endpoint` is without base path (currently `/v2`)

```bash
@test "Add a new TCP Request Rule to backend" {
	read -r SC _ < <(dataplaneapi POST "/path/to/post/endpoint" payload.json)
	[ "${SC}" = 201 ]
}
```

## Libraries

Some utilities have been developed to make test expectation and execution
smoother.

- [`dataplaneapi`](./libs/dataplaneapi.bash)
- [`auth_curl`](./libs/auth_curl.bash#)
- [`get_json_path`](./libs/get_json_path.bash)
- [`version`](./libs/version.bash)
- [`haproxy_config_setup`](./libs/haproxy_config_setup.bash)

Each library can be loaded using the relative path according to test file
location, as following:

```bash
`load '../../libs/${LIB_NAME}'
```

> The placeholder `${LIB_NAME}` doesn't need to contain the file extension.
