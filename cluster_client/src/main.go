package main

import (
	"net"

	"libvirt.org/go/libvirt"
)

//virsh --connect=qemu:///system list
//virsh --connect=qemu:///system domifaddr --source arp --domain VM_RAFT_3

func main()  {
    var connHyper, err =libvirt.NewConnect("qemu:///system")
    if err != nil {
        panic("failed connection to hypervisor")
    }

    var doms []libvirt.Domain
    doms,err = connHyper.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
    if err != nil {
        panic("failed retrieve domains")
    }

    var ip []libvirt.DomainInterface
    ip, err = doms[0].ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)

    var ipAddr string = ip[0].Addrs[0].Addr

    // var conn net.Conn
    _,err = net.Dial("tcp",ipAddr+":8080")

    if err != nil {
        panic("error dialing connection")
    }


}
