#!/bin/sh

work_dir=/root
mount_p=$work_dir/mount
raft_dir=$mount_p/raft
raft_executable=raft
raft_branch=raft_executables 
raft_url=https://github.com/mrmonopoly-cyber/raft_project_runtime.git 
my_ip_pos=$raft_dir/my_ip
others_ip_pos_work=/usr/share/raft/others_ip
others_ip_pos=$raft_dir/others_ip
my_ip_private=""
my_ip_public=""

while [[ -z $my_ip_private ]]; do
    my_ip_private=$(ip addr show | grep 10.0.0 | cut -d' ' -f 6 | cut -d'/' -f1);
done

while [[ -z $my_ip_public ]]; do
    my_ip_public=$(ip addr show | grep 200.168.122 | cut -d' ' -f 6 | cut -d'/' -f1);
done

if [[ -e $raft_dir ]];
then
    rm -rf $raft_dir
fi

if [[ ! -e $raft_dir ]]; then
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
echo $my_ip_private > $my_ip_pos
echo $my_ip_public >> $my_ip_pos
echo "setted ips"

while [[ $(wc -l $others_ip_pos_work | cut -d' ' -f 1) -lt 1 ]]; do
    sleep 1
done
cp $others_ip_pos_work $others_ip_pos
##start main program of raft

echo "started execution of raft"
$raft_dir/raft/raft/bin/raft

echo "exiting"
exit 0
