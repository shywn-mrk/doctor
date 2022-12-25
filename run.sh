#!/usr/bin/env bash
set -o nounset
set -o errexit

marketplaceDir="$1/vendor/digikala/supernova-digikala-marketplace"
if [ ! -d $marketplaceDir ]; then
    echo "Error: wrong path to supernova-env-dev"
fi

doctorDir="$marketplaceDir/docs/doctor"
if [ ! -d $doctorDir ]; then
    mkdir $doctorDir
fi

doctorFile="$doctorDir/index.html"

go build -o doctor ./*.go

rg -oUiIN --color never --crlf -E utf8 --no-heading --multiline-dotall --trim '/\*+\s*(@doctor.*?)\*/' --glob-case-insensitive -g '*.php' -r '$1' $1 | \
    rg -oiIN --color never --crlf -E utf8 --no-heading '(\*\s?)?(.*)' -r '$2' | \
    ./doctor | \
    pandoc --from gfm --to html --standalone --metadata 'title="Seller Docs"' > $doctorFile

fileLink="file://$doctorFile"

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     xdg-open $fileLink;;
    Darwin*)    open $fileLink;;
    *)          echo "Open in your browser: $fileLink"
esac
