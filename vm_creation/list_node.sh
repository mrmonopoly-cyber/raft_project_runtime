#!/bin/sh


# virsh --connect=qemu:///system domifaddr --source arp --domain <NODE_NUM>


function node_info (){
    virsh --connect=qemu:///system domifaddr --source arp --domain $1
}

node_num=0
if [[ -e ./.work/node_count ]]; then
    node_num=$(cat ./.work/node_count)
fi
node_name=""
public_ip=""

for ((i = 2; i < $node_num; i++)); do
    node_name=VM_RAFT_$i
    public_ip=$(node_info $node_name | grep 200.168.122)

    echo $node_name $public_ip
done
