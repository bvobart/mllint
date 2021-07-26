#!/bin/bash

cd $(dirname $0)
cd ../content/docs

cd categories
find . ! -name '_index.md' -type f -exec rm -f {} +

cd ../rules
find . ! -name '_index.md' -type f -exec rm -f {} +
