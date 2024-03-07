#!/bin/sh

others_ip=/root/raft_project_runtime/others_ip
nmap 192.168.122.0/24 -sn | grep report | cut -d' ' -f5 | tail -n +2 > $others_ip
