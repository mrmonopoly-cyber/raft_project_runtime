#!/bin/sh

sudo rm /var/lib/libvirt/images/disk_raft_node_*
virsh --connect=qemu:///system destroy VM_RAFT_*
