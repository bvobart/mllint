#!/bin/sh
CGO_ENABLED=0 go test -cover -coverprofile=coverage.txt -covermode=atomic ./...
