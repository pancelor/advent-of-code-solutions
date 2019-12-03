#!/usr/bin/env bash

set -e

DAY=${1?"Usage: $0 7"}

mkdir "$DAY"
cp main.go "$DAY"
touch "$DAY"/input.txt

URL_PREFIX='https://adventofcode.com/2019/day'
SUBLIME='/c/Program Files/Sublime Text 3/sublime_text.exe'
CHROME='/c/Program Files (x86)/Google/Chrome/Application/chrome.exe'
"$SUBLIME" -n "$DAY" "$DAY"/main.go "$DAY"/input.txt
"$CHROME" "${URL_PREFIX}/${DAY}/input" "${URL_PREFIX}/${DAY}"