#!/usr/bin/env bash
set -o nounset
set -o errexit

go build -o adapter ./*.go

rg -oUiIN --color never --crlf -E utf8 --no-heading --multiline-dotall --trim '/\*+\s*(@doctor.*?)\*/' --glob-case-insensitive -g '*.php' -r '$1' $1 | \
    rg -oiIN --color never --crlf -E utf8 --no-heading '(\*\s?)?(.*)' -r '$2' | \
    ./adapter
