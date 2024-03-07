#!/bin/sh

repo_raft=https://github.com/mrmonopoly-cyber/raft_project_runtime.git
branch=raft_executables

#setting mirrors
echo "updating mirrors"
systemctl start reflector

#partition creation
echo "making partitions"
echo "label: dos" | sfdisk /dev/vda
echo ";" | sfdisk /dev/vda
sfdisk -A /dev/vda 1

echo "formatting partitions"
mkfs.ext4 /dev/vda1

echo "mounting partitions"
mount /dev/vda1 /mnt

echo "updating keys"
pacman-key --init
pacman-key --populate
pacman -Sy

echo "installing packages"
pacstrap -K /mnt base dhcp dhcpcd syslinux bash openssh net-tools 
echo "installing kernel"
pacstrap -U /mnt /root/linux-6.7.8.arch1-1-x86_64.pkg.tar.zst /root/linux-firmware-20240220.97b693d2-1-any.pkg.tar.zst

echo "setting hostname"
echo "node_raft" > /mnt/etc/hostname

echo "mount support partitions(proc,sys,dev)"
mount --bind /proc /mnt/proc
mount --bind /sys /mnt/sys
mount --bind /dev /mnt/dev

echo "setting up network"
chroot /mnt/ systemctl enable dhcpcd

echo "setting ssh"
chroot /mnt/ systemctl enable sshd
cp -r /etc/ssh/ /mnt/etc
cp -r /root/.ssh /mnt/root

echo "setting up bootloader"
chroot /mnt syslinux-install_update -i -a -m
sed -i "s/sda3/vda1/g" /mnt/boot/syslinux/syslinux.cfg
sed -i "s/TIMEOUT 50/TIMEOUT 0/" /mnt/boot/syslinux/syslinux.cfg

echo "removing root passwd"
chroot /mnt passwd -d root

echo "setting up time to Rome"
chroot /mnt ln -sf /usr/share/zoneinfo/Europe/Rome /etc/localtime

#setup raft daemon bootstrap
echo "setting raft_node daemon"
cp /root/raft_daemon.service /mnt/usr/lib/systemd/system/
chroot /mnt systemctl enable raft_daemon.service

echo "copying bootstrap procedure"
cp /root/raft_bootstrap.sh /mnt/root
chroot /mnt/ chmod +x /root/raft_bootstrap.sh

echo "creating dir for raft program"
git clone --depth=1 $repo_raft -b $branch

#clean up
echo "umount all"
umount -R /mnt

echo "done"
exit 0


