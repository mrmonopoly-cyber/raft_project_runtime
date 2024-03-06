#!/bin/sh

my_ip=""
function get_ip(){
    my_ip=$(ip addr show | grep 192.168 | cut -d' ' -f 6 | cut -d'/' -f1);
}


work_dir=/root
my_ip_pos=$work_dir/my_ip
others_ip_pos=$work_dir/others_ip


while [[ -z $my_ip ]]; do
    get_ip
done

others_ip=$(arp -a | grep -v gateway | cut -d'(' -f2 | cut -d')' -f1)

touch $others_ip_pos
touch $my_ip_pos

echo "$my_ip" > $my_ip_pos
for IP in $others_ip; do
    echo "$IP" >> $others_ip_pos
done


##start main program of raft
$work_dir/raft/raft_main

exit 0


