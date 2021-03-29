#!/bin/sh
set -e # exit on first error
cd $(dirname $0)

# Prints $1 as yellow text
function print_yellow {
  local yellow='\033[0;33m'
  local nocolor='\033[0m'
  echo -e "${yellow}$1${nocolor}"
}

# Prints $1 as red text
function print_red {
  local red='\033[0;31m'
  local nocolor='\033[0m'
  echo -e "${red}$1${nocolor}"
}

print_yellow "> Cleaning build folder"
./clean.sh

print_yellow "> Building mllint..."
./build.sh --snapshot

print_yellow "> Copying ./dist to ./build/bin"
cp -r dist build/bin
print_yellow "> Copying ReadMe.md to ./build"
cp ReadMe.md build/


print_yellow "> Creating 'mllint' Python package..."
pip list | grep build > /dev/null || pip install build # install build package if not already installed
cd build && python -m build

echo
print_yellow "> Done! Find your package and wheels in  ./build/dist "
echo
