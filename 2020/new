#!/usr/bin/env bash

set -e

DAY=${1?"Usage: $0 7"}

mkdir "$DAY"
cp template.py "$DAY"/main.py
touch "$DAY"/input.txt

URL_PREFIX='https://adventofcode.com/2020/day'
SUBLIME='/c/Program Files/Sublime Text 3/sublime_text.exe'
CHROME='/c/Program Files (x86)/Google/Chrome/Application/chrome.exe'
"$SUBLIME" -n . "$DAY"/main.py "$DAY"/input.txt
"$CHROME" "${URL_PREFIX}/${DAY}/input" "${URL_PREFIX}/${DAY}"
