#!/bin/sh

first_user=$(whoami)
wiki_libvirt="https://wiki.archlinux.org/title/libvirt"
wiki_virt_manager="https://wiki.archlinux.org/title/Virt-manager"

if [[ $first_user == "root" ]]; then
    echo "you must not be root to execute this program"
    exit
fi


echo "installing the correct packages"
sudo pacman -Sy libvirt virt-manager iptables-nft dnsmasq virt-viewer dmidecode openbsd-netcat qemu-full archiso openssh 

echo "enabling the deamon for libvirtd"
sudo systemctl enable --now libvirtd

echo "adding libvirt group to user $first_user"
sudo usermod  -aG libvirt $first_user

echo "creating iso"
sudo mkarchiso -A raft_install.sh  -w work -o out -v -r ./iso_creation 
sudo mv ./out/* /var/lib/libvirt/images/raft_live_install.iso

echo "coping ssh key for remote node access"
cp ./raft_node_key ~/.ssh

echo "to open the gui to manage the virtual machine use virt-manager"
echo "for more infos check:"
echo $wiki_libvirt
echo $wiki_virt_manager


answer="p"
while [ $answer != "y" ] && [ $answer != "n" ]; do
    echo -n "do you want to open the two links to the wiki?:[y/n]"
    read answer
done

if [[ $answer == "y" ]]; then
    xdg-open $wiki_virt_manager
    xdg-open $wiki_libvirt
fi

exit
