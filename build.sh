#!/bin/sh
set -e # exit on first error
which goreleaser > /dev/null
goreleaser build --rm-dist $@
