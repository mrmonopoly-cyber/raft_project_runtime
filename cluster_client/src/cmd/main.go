package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"raft/client/src/internal/rpcs/client_request"
	"raft/client/src/internal/user_cli"
	"strings"

	"google.golang.org/protobuf/proto"
	"libvirt.org/go/libvirt"
)

//virsh --connect=qemu:///system list
//virsh --connect=qemu:///system domifaddr --source arp --domain VM_RAFT_3
//ip range of cluster 192.168.122.2 192.168.122.255

func main()  {
    var conn net.Conn
    var mex []byte
    var ipAddr string

    /*
    TODO:   
        implement different operation type with an interaction by the user,
        also implement add data in the payload field, when needed.
        OPERATION TO ADD:
            1- CREATE create ad file with name specified in the payload
            2- READ (add the name of the file to read in the payload)
  3- WRITE (write data in the file, if it does not exist the file will be created)
            4- RENAME (rename a file)
            5- DELETE (delete a file from the cluster)
    */
    

    var cli = usercli.NewCli()
    parameters, err := cli.Start()
    if err != nil {
      log.Panicln("Error in CLI: ", err)
    }

    if parameters["operation"] == "Create" || parameters["operation"] == "Write" {
      content, err1 := os.ReadFile(parameters["addParam"])
      if err1 != nil {
        log.Printf("Error in reading path %s : %s", parameters["addParam"], err1)
      }
      parameters["addParam"] = string(content)
    }

    log.Println(parameters)
    clientReq := clientrequest.NewClientRequest(parameters)
    mex, err2 := clientReq.Encode()
    
    if err2 != nil {
        panic("error encoding")
    }

    ipAddr = GetCluterNodeIp()
    conn = ConnectToLeader(ipAddr)
    SendCluster(conn, mex)
}

func EncodeMessage(req *protobuf.ClientReq) []byte{
    var mex []byte
    var err error

    mex, err = proto.Marshal(req)
    if err != nil {
        panic("error encoding")
    }
    return mex
}

func SendCluster(conn net.Conn, mex []byte){
    var err error

    _,err = conn.Write(mex)
    if err != nil {
        panic("error sending data")
    }
}

func ConnectToLeader(ipAddr string) net.Conn{
    var conn net.Conn
    var err error
    var mex []byte
    var leaderIp protobuf.LeaderIp

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

    return conn
}

func GetCluterNodeIp() string{
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

    if !strings.Contains(ipAddr,"192.168.122"){
        ipAddr = ip[1].Addrs[0].Addr
    }

    return ipAddr
}

func  Recv(conn net.Conn) ([]byte, error) {

    buffer := &bytes.Buffer{}

	// Create a temporary buffer to store incoming data
	var tmp []byte = make([]byte, 1024) // Initial buffer size
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
	return buffer.Bytes(), nil 

}
