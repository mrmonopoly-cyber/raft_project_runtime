#!/bin/sh

first_user=$(whoami)
wiki_libvirt="https://wiki.archlinux.org/title/libvirt"
wiki_virt_manager="https://wiki.archlinux.org/title/Virt-manager"

if [[ $first_user == "root" ]]; then
    echo "you must not be root to execute this program"
    exit
fi


echo "installing the correct packages"
sudo pacman -Sy libvirt virt-manager iptables-nft dnsmasq virt-viewer dmidecode openbsd-netcat

echo "enabling the deamon for libvirtd"
sudo systemctl enable --now libvirtd

echo "adding libvirt group to user $first_user"
sudo usermod  -aG libvirt $first_user

echo "getting and coping the install iso in the dir /var/lib/libvirt/images/"
wget https://download1526.mediafire.com/3f24vvgj9zwg4Gc18H9p4VxkEJjA4WeXsID9RlXQoPTP3B7weJfEdxpwyl9u-XKfhirBPkPfrOdNqSsW5HyrtxbHlruvivcd70EkEvHdOA35cpHviRP772jjuBCsPfg7v1ewiXHEERIo8yMvn5zdeqBcQ12sdnEX-i-AKS2MNoP_Zcw/m7oyfnz9q50f0ct/raft_live_install.iso
sudo cp ./raft_live_install.iso /var/lib/libvirt/images/

echo "installing ssh"
sudo pacman -Sy openssh

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
