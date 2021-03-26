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

[[ ! -d dist ]] && print_red "Error: no dist folder. Please run 'build.sh' before running 'test.package.sh'" && false

print_yellow "> Cleaning build folder"
rm -r build/bin build/build build/dist build/mllint.egg-info &> /dev/null || true

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
