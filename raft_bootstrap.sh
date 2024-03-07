#!/bin/sh


my_ip=""
function get_ip(){
    my_ip=$(ip addr show | grep 192.168 | cut -d' ' -f 6 | cut -d'/' -f1);
}


work_dir=/root
raft_dir=$work_dir/raft_project_runtime
my_ip_pos=$raft_dir/my_ip
raft_executable=raft_main

while [[ -z $my_ip ]]; do
    get_ip
done

touch $my_ip_pos

echo "$my_ip" > $my_ip_pos

##start main program of raft
$raft_dir/$raft_executable

exit 0


