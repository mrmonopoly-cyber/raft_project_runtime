#!/bin/sh

sources_dir=sources

work_dir=.work
node_count_file=$(pwd)/$work_dir/node_count
source_xml_vm=$(pwd)/$sources_dir/vm.xml
source_xml_pool=$(pwd)/$sources_dir/new_pool.xml
temp_xml_vm=new_node.xml
temp_xml_pool=new_pool.xml
subs_expr_disk=PATH_DISK
subs_expr_pool=PATH_MY_DIR
subs_expr_vm_name=RAFT_NODE_NAME
disk_prefix=disk_raft_node_
pwd=$(pwd)

if [[ ! -e $work_dir ]]; then
    mkdir .work
fi

cd $work_dir

if [[ ! -e $node_count_file ]]; then
    echo "2" > $node_count_file
fi

RAFT_NODE=$(cat $node_count_file)

echo $RAFT_NODE
next_node=$[$RAFT_NODE+1]

echo "$next_node" > $node_count_file

cp $source_xml_vm $temp_xml_vm
cp $source_xml_pool $temp_xml_pool

pwd_conv=$(echo $pwd\/.work | sed 's/\//\\\//g')

virsh --connect=qemu:///system vol-create-as default $disk_prefix$RAFT_NODE.qcow2  10G --format qcow2


pwd_conv=$(echo /var/lib/libvirt/images/${disk_prefix}${RAFT_NODE}.qcow2 | sed 's/\//\\\//g')

sed -i "s/$subs_expr_disk/$pwd_conv/" $temp_xml_vm
sed -i "s/$subs_expr_vm_name/VM_RAFT_$RAFT_NODE/" $temp_xml_vm

virsh --connect=qemu:///system create --file $temp_xml_vm


cd ..
