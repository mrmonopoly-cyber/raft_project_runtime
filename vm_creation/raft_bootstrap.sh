#!/bin/sh

my_ip=$(ip addr show | grep 192.168 | cut -d' ' -f 6 | cut -d'/' -f1)
others_ip=$(arp -a | grep -v gateway | cut -d'(' -f2 | cut -d')' -f1)

echo "$my_ip" > my_ip
for IP in $others_ip; do
    echo "$IP" >> others_ip
done


##start main program of raft
# go [program]
