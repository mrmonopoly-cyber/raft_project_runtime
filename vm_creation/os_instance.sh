#partition creation
echo "label: dos" | sfdisk /dev/vda
echo ";" | sfdisk /dev/vda
sfdisk -A /dev/vda 1

mount /dev/vda1 /mnt

pacman -Sy

pacstrap -K /mnt base linux linux-firmware dhcp dhcpcd iw grub vim bash go

mkdir /mnt/boot

echo "node_raft" > /mnt/etc/hostname

mount --bind /proc /mnt/proc
mount --bind /sys /mnt/sys
mount --bind /dev /mnt/dev

chroot /mnt/ systemctl enable dhcpcd


chroot /mnt grub-install /dev/vda
chroot /mnt grub-mkconfig -o /boot/grub/grub.cfg

chroot /mnt passwd -d root
chroot /mnt ln -sf /usr/share/zoneinfo/Europe/Rome /etc/localtime

chroot /mnt sed -i "s/GRUB_TIMEOUT=5/GRUB_TIMEOUT=0/" /etc/default/grub
