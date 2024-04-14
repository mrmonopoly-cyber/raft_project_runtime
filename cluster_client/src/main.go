package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"raft/client/protobuf/protobuf"
	"strings"

	"google.golang.org/protobuf/proto"
	"libvirt.org/go/libvirt"
)

//virsh --connect=qemu:///system list
//virsh --connect=qemu:///system domifaddr --source arp --domain VM_RAFT_3
//ip range of cluster 192.168.122.2 192.168.122.255

func  Recv(conn net.Conn) ([]byte, error) {

    buffer := &bytes.Buffer{}

	// Create a temporary buffer to store incoming data
	tmp := make([]byte, 1024) // Initial buffer size


    if conn == nil {
        return nil, errors.New("connection not instantiated")
    }
    //log.Println("want to read")
    //log.Printf("start reading from %v\n", this.GetIp())

    var bytesRead int = len(tmp)
    var errConn error
    var errSavi error
    for bytesRead == len(tmp){
		// Read data from the connection
		bytesRead, errConn = conn.Read(tmp)

        // Write the read data into the buffer
        _, errSavi = buffer.Write(tmp[:bytesRead])
        
        // check error saving
        if errSavi != nil {
            return nil, errSavi
        }

		if errConn != nil {
			if errConn != io.EOF {
				// Handle other errConnors
				return nil, errConn
			}

            if errConn == io.EOF {
                return nil, errConn
            }
			break
		}
	}

    //log.Printf("end reading from %v : %v\n", this.GetIp(), buffer)
    
    //log.Println("found no error, received message: ", buffer)
	return buffer.Bytes(), nil 

}


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
    var mex []byte
    var leaderIp protobuf.LeaderIp

    if strings.Contains(ipAddr,"192.168.122"){
        ipAddr = ip[0].Addrs[0].Addr
    }

    // ipAddr = "192.168.122.87"
    
    for {
        conn,err = net.Dial("tcp",ipAddr+":8080")
        if err != nil {
            log.Panicf("error dialing conection %v\n",err)
        }
        mex, err = Recv(conn)
        err = proto.Unmarshal(mex,&leaderIp)
        if err != nil {
            panic("failed decoding ip of leader")
        }

        log.Println("new leader :", leaderIp.Ip)
        if strings.Contains(leaderIp.Ip,"ok") {
            break
        }
        conn.Close()
        ipAddr = leaderIp.Ip
        err = nil
    }

    println(ipAddr)
    // conn,err= net.Dial("tpc",isLeader+":8080")
    // if err != nil {
    //     panic("error dialing connection with leader")
    // }
    var req = protobuf.ClientReq{
        Op: protobuf.Operation_READ,
    }

    mex, err = proto.Marshal(&req)
    if err != nil {
        panic("error encoding")
    }

    println("sending data")
    _,err = conn.Write(mex)

    if err != nil {
        panic("error sending data")
    }


    



}
