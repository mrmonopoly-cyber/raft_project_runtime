#!/bin/sh

vms=$(virsh --connect=qemu:///system list | grep VM_RAFT | cut -d ' ' -f 6)
volumes=$(virsh --connect=qemu:///system vol-list --pool default| grep disk_raft_node)

for VM in $vms; do
    virsh --connect=qemu:///system destroy $VM
done

for VOL in $volumes; do
    virsh --connect=qemu:///system vol-delete --pool default $VOL
done

rm -rf .work
