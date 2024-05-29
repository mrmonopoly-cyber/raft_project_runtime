#!/bin/sh

first_user=$(whoami)
wiki_libvirt="https://wiki.archlinux.org/title/libvirt"
wiki_virt_manager="https://wiki.archlinux.org/title/Virt-manager"

function separation() {
    echo "--------------------------------------------------------------------------------"
}

function setup_network() {
    echo "cleaning possible legacy: $2"
    virsh --connect=qemu:///system net-destroy --network $2
    virsh --connect=qemu:///system net-undefine --network $2
    echo "setting up network: $2"
    virsh --connect=qemu:///system net-define --file $1
    virsh --connect=qemu:///system net-autostart --network $2
    virsh --connect=qemu:///system net-start --network $2
}

if [[ $first_user == "root" ]]; then
    echo "you must not be root to execute this program"
    exit
fi


echo "installing the correct packages"
sudo pacman -Sy go libvirt virt-manager iptables-nft dnsmasq virt-viewer dmidecode openbsd-netcat qemu-full archiso openssh --noconfirm
separation


echo "enabling the deamon for libvirtd"
sudo systemctl enable --now libvirtd
separation

echo "adding libvirt group to user $first_user"
sudo usermod  -aG libvirt $first_user
separation

setup_network ./raft_net_private.xml raft_network_private
setup_network ./raft_net_public.xml raft_network_public
separation

echo "creating iso"
sudo mkarchiso -A raft_install.sh  -w work -o out -v -r ./iso_creation 
sudo mv ./out/raft_live_install.iso--x86_64.iso /var/lib/libvirt/images/raft_live_install.iso
separation

echo "coping ssh key for remote node access"
cp ./raft_node_key ~/.ssh
separation

exit
