#!/bin/bash
# Creates an SVG of a run of mllint.
# Assumes you're using it on one of the mllint-example-projects with poetry and all dependencies already installed.

# Depends on:
# - svg-term: https://github.com/marionebl/svg-term-cli
# - asciinema: https://asciinema.org/
# - bash (don't run it on Windows)
# 
# TODO: there are some tiny artefacts in the completed SVG relating to issues with Powerline fonts not being available.
# I tried this: https://gist.github.com/danielfullmer/5e29d1e9534dded5c183
# using this: https://graphicdesign.stackexchange.com/questions/10733/how-do-i-use-a-custom-font-in-an-svg-image-on-my-site
# but I wasn't able to get it fixed. Have any idea to fix it? Feel free to tackle it! :)

# exit on first error
set -e

# File location of this script.
script_dir="$(dirname $0)"

# The folder to run mllint on and filename for the asciicast and finished .svg
folder="$1"
filename="$2"
[[ -z "$folder" ]] && folder="."
[[ -z "$filename" ]] && filename="mllint-run-$(date +%F.%T).cast"

echo "> Recording and creating SVG in '$folder'..."
svg-term --command "bash $script_dir/demo-script.sh $folder" --term konsole --profile "$script_dir/demo-colors.colorscheme" --out "$filename.svg"

echo "> Making final adjustments..."
# svg-term apparently has some trouble selecting the right colour from the color scheme for text with a background colour (i.e. the terminal prompt), so we manually fix that here.
sed -i 's/fill:#71bef2/fill:#df00fe/g' "$filename.svg"
sed -i 's/fill:#dbab79/fill:#e6b822/g' "$filename.svg"

echo "> Done: $filename.svg"
