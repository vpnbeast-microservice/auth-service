# Auth Service
[![CI](https://github.com/vpnbeast/auth-service/workflows/CI/badge.svg?event=push)](https://github.com/vpnbeast/auth-service/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/vpnbeast/auth-service)](https://hub.docker.com/r/vpnbeast/auth-service/)
[![Go Report Card](https://goreportcard.com/badge/github.com/vpnbeast/auth-service)](https://goreportcard.com/report/github.com/vpnbeast/auth-service)
[![codecov](https://codecov.io/gh/vpnbeast/auth-service/branch/master/graph/badge.svg)](https://codecov.io/gh/vpnbeast/auth-service)
[![Go version](https://img.shields.io/github/go-mod/go-version/vpnbeast/auth-service)](https://github.com/vpnbeast/auth-service)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This API gets the `/users/authenticate` requests from gateway with username and password, validates the request, generates access and refresh token and then
return the response which contains detailed information about the user.

## Prerequisites
auth-service requires [vpnbeast/config-service](https://github.com/vpnbeast/config-service) to fetch configuration. Configurations
are stored at [vpnbeast/properties](https://github.com/vpnbeast/properties).

## Configuration
This project fetches the configuration from [config-service](https://github.com/vpnbeast/config-service).
But you can still override them with environment variables:
```
SERVER_PORT
METRICS_PORT
METRICS_ENDPOINT
WRITE_TIMEOUT_SECONDS
READ_TIMEOUT_SECONDS
ISSUER
PRIVATE_KEY
PUBLIC_KEY
ACCESS_TOKEN_VALID_IN_MINUTES
REFRESH_TOKEN_VALID_IN_MINUTES
ENCRYPTION_SERVICE_URL
DB_URL
DB_DRIVER
HEALTH_PORT
HEALTH_ENDPOINT
DB_MAX_OPEN_CONN
DB_MAX_IDLE_CONN
DB_CONN_MAX_LIFETIME_MIN
HEALTHCHECK_MAX_TIMEOUT_MIN
```

## Development
This project requires below tools while developing:
- [Golang 1.16](https://golang.org/doc/go1.16)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)

## License
Apache License 2.0
