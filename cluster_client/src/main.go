package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"libvirt.org/go/libvirt"
	// "raft/client/protobuf/protobuf"
)

//virsh --connect=qemu:///system list
//virsh --connect=qemu:///system domifaddr --source arp --domain VM_RAFT_3
//ip range of cluster 192.168.122.2 192.168.122.255

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

    println(ip[0].Addrs[0].Addr)

    var ipAddr string = ip[0].Addrs[0].Addr
    var conn net.Conn
    var isLeader string

    if strings.Contains(ipAddr,"192.168.122"){
        ipAddr = ip[0].Addrs[0].Addr
    }
    
    for {
        conn,err = net.Dial("tcp",ipAddr+":8080")
        if err != nil {
            log.Panicf("error dialing conection %v\n",err)
        }
        var reader = bufio.NewReader(conn)
        isLeader, err = reader.ReadString('\n')

        log.Println("new leader :", isLeader)
        if isLeader != "ok" {
            break
        }
    }
    
    println(isLeader)

    // var newMex = protobuf.ClientReq{
    //     Op: protobuf.Operation_READ,
    //     Payload: []byte("pippo.txt"),
    // }

    conn.Close()


}
