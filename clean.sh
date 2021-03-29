#!/bin/sh
cd $(dirname $0)
rm -r dist/ build/bin build/build build/dist build/mllint.egg-info || true
