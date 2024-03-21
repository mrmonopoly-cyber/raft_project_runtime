#!/bin/sh

if [[ $# != 1 ]]; then
    echo "invalid number of argument, specify ONE branch"
    echo "this script allow the user to specify which branch of raft will be executed by the
    vms, to do that you simply need to run this script passing as argument the name 
    of the branch you want to use"
    exit
fi


git config pager.branch false
branch=$1 
clean=$(echo $branch | sed  "s/-/\\\-/")
exist=$(git branch -a | awk "/\y$clean\y/ && !/\yRaftDev\[^[:space:]]*\y/"| head -1)
repo=$(git remote get-url origin)

if [[ -n $exist ]]; then
    echo "cleaning"
    rm -rf ./raft
    echo "cloning"
    git clone --depth=1 --branch $branch $repo ./raft/
    rm ./raft/.git
    cd ./raft/raft
    echo "building"
    ./build.sh
    cd ../..
    git add .
    git commit -m "updated run code"
    git push
fi
echo "done"

