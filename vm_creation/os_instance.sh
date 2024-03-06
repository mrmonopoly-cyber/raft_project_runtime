#partition creation
echo "making partitions"
echo "label: dos" | sfdisk /dev/vda
echo ";" | sfdisk /dev/vda
sfdisk -A /dev/vda 1

echo "making filesystem"
mkfs.ext4 /dev/vda1

echo "mounting partitions"
mount /dev/vda1 /mnt

echo "updating keys"
pacman-key --init
pacman-key --populate
pacman -Sy

echo "mounting partitions"
mount /dev/vda1 /mnt
echo "installing packages"
pacstrap -K /mnt base linux linux-firmware dhcp dhcpcd iw grub vim bash go openssh

if [[ -e /mnt/boot ]]; then
    echo "creating boot dir"
    mkdir /mnt/boot
fi

echo "setting hostname"
echo "node_raft" > /mnt/etc/hostname

mount --bind /proc /mnt/proc
mount --bind /sys /mnt/sys
mount --bind /dev /mnt/dev


echo "setting up network"
chroot /mnt/ systemctl enable dhcpcd

echo "setting ssh"
chroot /mnt systemctl enable sshd
cp -r /etc/ssh /mnt/etc
cp -r /root/.ssh /mnt/root

echo "setting up bootloader"
chroot /mnt grub-install --target=i386-pc /dev/vda
chroot /mnt sed -i "s/GRUB_TIMEOUT=5/GRUB_TIMEOUT=0/" /etc/default/grub
chroot /mnt grub-mkconfig -o /boot/grub/grub.cfg

echo "removing root passwd"
chroot /mnt passwd -d root

echo "setting up time to Rome"
chroot /mnt ln -sf /usr/share/zoneinfo/Europe/Rome /etc/localtime

#setup raft daemon bootstrap
echo "setting raft_node daemon"
cp /root/raft_daemon.service /usr/lib/systemd/system/
chroot /mnt systemctl enable raft_daemon.service

reboot

