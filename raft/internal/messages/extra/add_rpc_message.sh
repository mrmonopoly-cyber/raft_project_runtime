#!/bin/sh

extra_dir=./extra

if [[ $# -lt 1 ]]; 
then
    echo "invalid input, give one name"
    exit 1
fi

if [[ ! -d $extra_dir ]];
then
    echo "execute the script in the messages directory"
    exit 2
fi

new_rcp=$1

mkdir $new_rcp
rpc_file=$new_rcp.go
cp $extra_dir/dummy_rpc.go ./$new_rcp/$rpc_file
sed -i "s/NEW_RPC/$new_rcp/g" $new_rcp/$rpc_file
sed -i "s/extra/$new_rcp/g" $new_rcp/$rpc_file
