#!/usr/bin/env bash
set -o nounset


#################### Dependency Checking
dep_go=1
go version >> /dev/null 2>&1
if [ $? -ne 0 ]; then
    dep_go=0
fi

dep_rg=1
rg --version >> /dev/null 2>&1
if [ $? -ne 0 ]; then
    dep_rg=0
fi

dep_pandoc=1
pandoc -v >> /dev/null 2>&1
if [ $? -ne 0 ]; then
    dep_pandoc=0
fi


if [ $dep_go -ne 1 ] || [ $dep_rg -ne 1 ] || [ $dep_pandoc -ne 1 ]; then
    deps='Error: You need to install [ '
    if [ $dep_go -ne 1 ]; then
        deps="$deps Go Compiler, "
    fi
    if [ $dep_rg -ne 1 ]; then
        deps="$deps RipGrep, "
    fi
    if [ $dep_pandoc -ne 1 ]; then
        deps="$deps Pandoc, "
    fi
    deps="${deps::${#deps}-2} ]"
    echo $deps
    exit 1
fi


echo "You can use run.sh script to generate documentation"