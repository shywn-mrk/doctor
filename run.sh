#!/usr/bin/env bash
set -o nounset
set -o errexit


go build -o adapter ./*.go

ag -i --php --only-matching --nofilename  --nocolor --nobreak '/\*+?\s*?\@doctor(.*?\n)*?.*?\*/' $1 | ./adapter