# API Template

This repo is intended to be a walking skeleton of a simple API that can be used as a starting point when designing a new HTTP API. 

The application tries to adhere to modern application patterns including:
- Graceful shutdown of server
- Dynamic Configuration with environment variables (see Sourcing Environment Variables)
- Static analysis with [golangci-lint](https://github.com/golangci/golangci-lint)
- Structured logging with [Zero Logs](https://github.com/rs/zerolog)
- Request logging and panic recovery middleware


## Sourcing Environment Variables

This application is designed to source environment variables with the help of [direnv](https://direnv.net/).

Steps to install:

To install binary builds you can run this bash installer:
```sh
curl -sfL https://direnv.net/install.sh | bash
```

Fetch the binary, `chmod +x direnv` and put it somewhere in your PATH.

Navigate to application directory and run

```sh
sudo direnv allow
```

Reload the terminal and you should see the following console statement: `direnv: loading ~/go/src/github.com/benjaminslavin1/api-template/.envrc` and the sourced variables 

## Development

The Makefile is the entrypoint for common development tasks.

These commands MAY depend on the .envrc.

There are currently 2 default targets:

`make test` — runs `go test` : runs all of the app's tests
`make lint` — runs `golangci-lint run` : runs all of the app's linting

## Linting with Drone

TODO