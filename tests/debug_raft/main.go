package main

import (
	"debug_raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

func main()  {

    var conn net.Conn
    var err error
    var leader = strings.Split(os.Args[2], ":")[0]

    conn,err = net.Dial("tcp",os.Args[1])
    if err != nil{
        log.Panicln(err)
    }

    // newConf(conn,leader)
    // addConf(conn,leader)
    remConf(conn,leader)
}

func remConf(conn net.Conn, leader string){
    var addConfReq protobuf.ChangeConfReq = protobuf.ChangeConfReq{
        Op: protobuf.AdminOp_CHANGE_CONF_REM,
        Conf: &protobuf.ClusterConf{
            // Conf: []string{"10.0.0.20"},
            Conf: []string{"10.0.0.29"},
            Leader: &leader,
            
        },
    }
    var rawMex,err = proto.Marshal(&addConfReq)
    if err != nil{
        log.Panicln(err)
    }
    rawMex = encodeGeneric(conn, rawMex)
    readReturn(conn)
}

func addConf(conn net.Conn, leader string){
    var addConfReq protobuf.ChangeConfReq = protobuf.ChangeConfReq{
        Op: protobuf.AdminOp_CHANGE_CONF_ADD,
        Conf: &protobuf.ClusterConf{
            Conf: []string{"10.0.0.29"},
            // Conf: []string{"10.0.0.195"},
            Leader: &leader,
            
        },
    }
    var rawMex,err = proto.Marshal(&addConfReq)
    if err != nil{
        log.Panicln(err)
    }
    rawMex = encodeGeneric(conn, rawMex)
    readReturn(conn)
}

func newConf(conn net.Conn, leader string){
    var newConfReq protobuf.ChangeConfReq = protobuf.ChangeConfReq{
        Op: protobuf.AdminOp_CHANGE_CONF_NEW,
        Conf: &protobuf.ClusterConf{
            // Conf: []string{leader},
            // Conf: []string{leader, "10.0.0.195","10.0.0.20"},
            Conf: []string{leader,"10.0.0.29"},
            Leader: &leader,
            
        },
    }
    var rawMex,err = proto.Marshal(&newConfReq)
    if err != nil {
        panic(err)
    }
    rawMex = encodeGeneric(conn, rawMex)
    readReturn(conn)
}

func encodeGeneric(conn net.Conn, rawMex []byte) []byte{
    var genericMex protobuf.Entry 
    var err error
    genericMex.OpType = protobuf.MexType_NEW_CONF
    genericMex.Payload = rawMex
    rawMex,err = proto.Marshal(&genericMex)
    if err != nil{
        log.Panicln(err)
    }

    conn.Write(rawMex)
    return rawMex
}


func readReturn(conn net.Conn){
    var rawMex = make([]byte, 2048)
    var returnValue protobuf.ClientReturnValue = protobuf.ClientReturnValue{}

    conn.Read(rawMex)
    proto.Unmarshal(rawMex,&returnValue)
    log.Println(returnValue.Description)
    log.Println(returnValue.ExitStatus)
}
