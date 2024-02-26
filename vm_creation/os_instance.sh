#partition creation
echo "label: dos" | sfdisk /dev/vda
echo ";" | sfdisk /dev/vda

mount /dev/vda1 /mnt

pacman -Sy

pacstrap -K /mnt base linux linux-firmware dhcp dhcpcd iw grub vim bash

mkdir /mnt/boot

echo "node_raft" > /mnt/etc/hostname

mount --bind /proc /mnt/proc
mount --bind /sys /mnt/sys
mount --bind /dev /mnt/dev

chroot /mnt/ systemctl enable dhcpcd


chroot /mnt grub-install /dev/vda
chroot /mnt grub-mkconfig -o /boot/grub/grub.cfg
