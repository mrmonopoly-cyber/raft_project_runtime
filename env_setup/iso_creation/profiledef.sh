#!/usr/bin/env bash
# shellcheck disable=SC2034

iso_name="raft_live_install.iso"
iso_label=""
iso_publisher=""
iso_application="raft cluster system"
iso_version=""
install_dir="arch"
buildmodes=('iso')
bootmodes=('bios.syslinux.mbr' 'bios.syslinux.eltorito')
arch="x86_64"
pacman_conf="pacman.conf"
airootfs_image_type="erofs"
airootfs_image_tool_options=('-zlzma,109' -E 'ztailpacking,fragments,dedupe')
file_permissions=(
  ["/etc/shadow"]="0:0:400"
)
