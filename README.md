## Auth Service

[![CI](https://github.com/vpnbeast/auth-service/workflows/CI/badge.svg?event=push)](https://github.com/vpnbeast/auth-service/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/vpnbeast/auth-service)](https://hub.docker.com/r/vpnbeast/auth-service/)
[![Go Report Card](https://goreportcard.com/badge/github.com/vpnbeast/auth-service)](https://goreportcard.com/report/github.com/vpnbeast/auth-service)

This API gets the `/users/authenticate` requests from gateway with username and password, validates the request, generates access and refresh token and then
return the response which contains detailed information about the user.

### Configuration
List of supported environment variables:
- TBD

### Development
This project requires below tools while developing:
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)
