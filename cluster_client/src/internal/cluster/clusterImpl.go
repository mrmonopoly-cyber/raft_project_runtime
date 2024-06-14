package cluster

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"raft/client/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"raft/client/src/internal/rpcs"
	changeconfigreq "raft/client/src/internal/rpcs/change_config_req"
	inforequest "raft/client/src/internal/rpcs/info_request"
	inforesponse "raft/client/src/internal/rpcs/info_response"
	"raft/client/src/internal/utility"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
	"libvirt.org/go/libvirt"
)

type ConfigChangeOp int 

const (
  NEW ConfigChangeOp = 0
  CHANGE ConfigChangeOp = 1
)


type clusterImpl struct {
  IPs []utility.Pair[string,string]
  leaderIP string
  conn net.Conn
}

/*
 * Enstablish a connection to a random node in the cluster 
 * (only at beginning when no leader exists)
*/
func (this *clusterImpl) EnstablishConnection() {
    var conn net.Conn
    var err error
    var addr string = this.IPs[0].Fst

    conn, err = net.Dial("tcp",addr+":8080")
    if err != nil {
      log.Panicf("error dialing conection %v\n",err)
    }

    this.conn = conn
}

/*
 * Connect to cluster leader
*/
func (this *clusterImpl) ConnectToLeader() {
    var conn net.Conn
    var err error
    var mex []byte
    var leaderIp protobuf.LeaderIp
    var addr string = this.IPs[0].Fst

    for {
        conn,err = net.Dial("tcp",addr+":8080")
        if err != nil {
            log.Panicf("error dialing conection %v\n",err)
        }
        mex, err = recv(conn)
        err = proto.Unmarshal(mex,&leaderIp)
        if err != nil {
            panic("failed decoding ip of leader")
        }

        log.Println("new leader :", leaderIp.Ip)
        if strings.Contains(leaderIp.Ip,"ok") {
            break
        }
        conn.Close()
        this.leaderIP = leaderIp.Ip
        err = nil
    }

    this.conn = conn
}

/*
 * Removing a node by its public IPs
*/
func (this *clusterImpl) RemoveNode(IP string) error {
  var found bool = slices.ContainsFunc(this.IPs, func(element utility.Pair[string, string]) bool {
                                              return IP == element.Fst}) 

  if found == false {
    return errors.New("IP address not found")
  }

  var IPs []utility.Pair[string,string]
  for _,i := range this.IPs {
    if i.Fst != IP {
      IPs = append(IPs, i)
    }
  }
  this.IPs = IPs

  return nil
}

/*
 * Make an info request of type "config" and send it to the cluster
 * in order to get the current configuration
*/
func (this *clusterImpl) GetConfig() rpcs.Rpc {
  var req, resp rpcs.Rpc
  var reqByte, respByte []byte
  var err, errResp, errDec error

  req = inforequest.NewInfoRequest()
  resp = inforesponse.NewInfoResonse()

  reqByte, err = req.Encode()
  if err != nil {
    log.Println(err)
  }
  this.conn.Write(reqByte)

  respByte, errResp = recv(this.conn)
  if errResp != nil {
    log.Println(errResp)
  }
  
  errDec = resp.Decode(respByte)
  if errDec != nil {
    log.Println(errDec)
  }

  return resp
}

/*
 * Send a new configuration to the cluster, it distinguishes 
 * between a NEW configuration (i.e. starting config) 
 * and a CHANGE configuration (i.e. modified config)
*/
func (this *clusterImpl) SendConfig(op ConfigChangeOp) {
  var req rpcs.Rpc
  var reqByte []byte
  var err error
  var config protobuf.ClusterConf

  switch op {
    case CHANGE: 
      config = protobuf.ClusterConf{
        Conf: this.getPrivateIPs(),
        Leader: &this.leaderIP,
      }
      req = changeconfigreq.NewConfigRequest(config, protobuf.AdminOp_CHANGE_CONF_CHANGE)

    default: 
      config = protobuf.ClusterConf{
        Conf: this.getPrivateIPs(),
        Leader: nil,
      }
      req = changeconfigreq.NewConfigRequest(config, protobuf.AdminOp_CHANGE_CONF_NEW)
   }

  reqByte, err = req.Encode()
  if err != nil {
    log.Println(err)
  }
  this.conn.Write(reqByte)
}

/*
 * Receiver method
*/
func recv(conn net.Conn) ([]byte, error) {

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

func (this *clusterImpl) getPrivateIPs() []string {
  var IPs []string
  for _,i := range this.IPs {
    IPs = append(IPs, i.Snd)
  }
  return IPs
}

func (this *clusterImpl) getPublicIPs() []string {
  var IPs []string
  for _,i := range this.IPs {
    IPs = append(IPs, i.Fst)
  }
  return IPs
}
