#!/bin/sh

if [[ $# != 1 ]]; then
    echo "invalid number of argument, specify ONE branch"
    exit
fi

git config pager.branch false
exist=$(git branch -a | awk '/\yRaftDev\y/ && !/\yRaftDev\-[^[:space:]]*\y/' | head -1)

if [[ -n $exist ]]; then
    git merge $exist
fi

cd ./raft
./build.sh
