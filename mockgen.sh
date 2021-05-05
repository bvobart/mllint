#!/bin/bash
# Uses the `mockgen` tool to generate mocks implementations of `mllint`s main interfaces.
# If you haven't got `mockgen`, use `go install github.com/golang/mock/mockgen` to install it :)

cd $(dirname $0)

mockgen -source=api/linter.go > api/mock_api/linter.go