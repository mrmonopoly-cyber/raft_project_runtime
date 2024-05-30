#!/bin/sh

work=/usr/share/raft
others_ip=${work}/others_ip
others_ip_temp=${work}/others_ip.temp

ip_assigned=$(ip a | grep "10.0.0")

if -z $ip_assigned ; then
    sleep 5
fi

nmap 10.0.0.0/24 -sn | grep report | cut -d' ' -f5 | tail -n +2 | head -n -1 > $others_ip_temp
mv $others_ip_temp $others_ip
