#!/bin/sh

others_ip=/root/mount/raft/others_ip
nmap 192.168.122.0/24 -sn | grep report | cut -d' ' -f5 | tail -n +2 | head -n -1 > $others_ip
