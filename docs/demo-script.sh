#!/bin/bash
# Designed to be used with svg-term in create-demo-svg.sh
# Example Usage: svg-term --command "bash demo-script.sh $folder" --out "filename.svg"
# Inspiration taken from how 'fd' creates its screencast demo: https://github.com/sharkdp/fd/blob/master/doc/screencast.sh

set -e
set -u

# Folder where mllint should be executed
FOLDER="$1"

FG_BLACK_BG_BLUE="\033[30;44m"
FG_BLACK_BG_YELLOW="\033[30;43m"
COLOR_RESET="\033[0m"
# Terminal prompt to show during the demo
PROMPT="$FG_BLACK_BG_BLUE ~/tudelft/thesis/mllint-example-projects ▶$FG_BLACK_BG_YELLOW  main ▶$COLOR_RESET"

function prompt {
  printf "%b " "$PROMPT"
}

# This function splits its first argument into individual characters and prints each of them with a slight (80 ms) delay, to simulate someone typing.
function type {
  arg=$1
  for ((i=0;i<${#arg};i++)); do
    echo -n "${arg:i:1}"
    sleep 0.08
  done
}

# This function takes whatever arguments it is given, types them out using the `type` function above, sleeps for a little bit before 'pressing enter', executing the entire command and printing a new terminal prompt.
function enter {
  input="$@"
  wait_time=0.6

  type "$input"
  sleep "$wait_time" 
  echo
  eval "$input"
  echo
  prompt
}

# This function performs the actual demo of mllint
function demo {
  cd "$FOLDER"
  
  # silently remove report from previous run
  [[ -f report.md ]] && rm report.md

  prompt
  sleep 0.5
  
  enter tree
  sleep 1.5
  
  enter poetry run mllint -o report.md
  sleep 1.5

  enter poetry run mllint render report.md
  sleep 3

  # this final echo is there so that the SVG will respect the previous `sleep`
  echo
}

demo
