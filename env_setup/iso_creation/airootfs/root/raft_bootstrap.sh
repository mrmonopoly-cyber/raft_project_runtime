#!/bin/sh

work_dir=/root
mount_p=$work_dir/mount
raft_dir=$mount_p/raft
my_ip_pos=$raft_dir/my_ip
raft_executable=raft
raft_branch=raft_executables 
raft_url=https://github.com/mrmonopoly-cyber/raft_project_runtime.git 
my_ip=""

while [[ -z $my_ip ]]; do
    my_ip=$(ip addr show | grep 192.168 | cut -d' ' -f 6 | cut -d'/' -f1);
done

if [[ ! -b $raft_dir ]];
then
    rm -rf $raft_dir
fi

if [[ ! -b $raft_dir ]]; then
    mkdir $raft_dir
fi

is_mounted=$(mount | grep "/dev/vda")
if [[ -z $is_mounted ]]
then
    mkfs.ext4 /dev/vda
    mount /dev/vda $mount_p
fi


git clone --branch $raft_branch --depth=1 $raft_url $raft_dir
echo "repo cloned"

touch $my_ip_pos
echo "$my_ip" > $my_ip_pos
echo "setted ips"

##start main program of raft

echo "started execution of raft"
$raft_dir/raft/bin/$raft_executable

echo "exiting"
exit 1
