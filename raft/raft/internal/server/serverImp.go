package server

import (
	"fmt"
	"log"
	"net"
	genericmessage "raft/internal/genericMessage"
	"raft/internal/node"
	"raft/internal/raft_log"
	"raft/internal/raftstate"
	state "raft/internal/raftstate"
	confpool "raft/internal/raftstate/confPool"
	"raft/internal/rpcs"
	"strings"
	"sync"
    "raft/internal/rpcs/redirection"
)

type pairMex struct{
    payload rpcs.Rpc
    sender string
    workdone chan int
}

type server struct {
    wg             sync.WaitGroup
    _state         state.State
    messageChannel chan pairMex
    listener       net.Listener
    clientList      sync.Map
}

func (s *server) Start() {
    log.Println("Start accepting connections")

    s.wg.Add(1)
    go s.acceptIncomingConn()
    s.wg.Add(1)
    go s.run()

    s.wg.Wait()
    log.Println("run finished")
}

//utility
func (this *server) connectToNodes(serversIp []string, port string) ([]string,error){
    var failedConn []string = make([]string, 0)
    var err error

	for i := 0; i < len(serversIp)-1; i++ {
		var new_node node.Node
        var nodeConn net.Conn
        var erroConn error

        nodeConn,erroConn = net.Dial("tcp",serversIp[i]+":"+port)
        if erroConn != nil {
            log.Println("Failed to connect to node: ", serversIp[i])
            failedConn = append(failedConn, serversIp[i])
            err = erroConn
            continue
        }
        new_node = node.NewNode(serversIp[i], port, nodeConn)
        log.Printf("connected to new node, storing it: %v\n", new_node.GetIp())
        this._state.UpdateNodeList(confpool.ADD,new_node)
        go (*this).internalNodeConnection(new_node)
	}

    return failedConn,err
}

func (s *server) acceptIncomingConn() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Failed on accept: ", err)
			continue
		}

		tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			fmt.Println("Connection is not using TCP")
			continue
		}

		var newConncetionIp string = tcpAddr.IP.String()
		var newConncetionPort string = string(rune(tcpAddr.Port))
        log.Printf("new connec: %v\n", newConncetionIp)
        log.Printf("node with ip %v not found", newConncetionIp)
        var newNode node.Node = node.NewNode(newConncetionIp, newConncetionPort,conn)
        go s.handleConnection(newNode)
	}
}

func (s* server) handleConnection(workingNode node.Node){
    if strings.Contains(workingNode.GetIp(), "10.0.0") {
        s._state.UpdateNodeList(confpool.ADD,workingNode)
        s.internalNodeConnection(workingNode)
        return
    }
    s.externalAgentConnection(workingNode)
}

func (s *server) externalAgentConnection(agent node.Node){
    var leaderIp = s._state.GetLeaderIp(raftstate.PUB)
    var rawMex []byte
    var err error
    var resp rpcs.Rpc 
    var inputMex rpcs.Rpc 
    var clientReq pairMex = pairMex{}

    s.clientList.Store(agent.GetIp(),agent)

    //INFO: send to the client who you think is the leader, (blank means you are not in the conf)
    resp = Redirection.NewredirectionRPC(leaderIp)
    s.encodeAndSend(resp,agent)

    for{
        rawMex,err = agent.Recv()
        if err != nil || rawMex == nil {
            fmt.Printf("error in reading from node %v with error %v\n",agent.GetIp(), rawMex)
            break
        }
        inputMex,err = genericmessage.Decode(rawMex)
        if err != nil{
            log.Println(err)
            break
        }
        clientReq.payload = inputMex
        clientReq.sender = agent.GetIp()
        clientReq.workdone = make(chan int)
        s.messageChannel <- clientReq
        log.Println("waiting agent completion on chann: ",clientReq.workdone)
        <- clientReq.workdone
    }

    log.Println("Done serving client: ", agent.GetIp())
    agent.CloseConnection()
}

func (s *server) internalNodeConnection(workingNode node.Node) {
    var nodeIp = workingNode.GetIp()
    var message []byte
    var errMes error
    var rpcMex rpcs.Rpc
    var nodeReq pairMex = pairMex{}


    for{
        message, errMes = workingNode.Recv()
        if errMes != nil {
            log.Printf("error in reading from node %v with error %v\n",nodeIp, errMes)
            break
        }else if message != nil {
            rpcMex,errMes = genericmessage.Decode(message)
            if errMes != nil {
                log.Println(errMes)
                continue
            }
            nodeReq.payload = rpcMex
            nodeReq.sender = workingNode.GetIp()
            nodeReq.workdone = make(chan int)
            s.messageChannel <- nodeReq
            <- nodeReq.workdone
        }
    }
    workingNode.CloseConnection()
    s._state.UpdateNodeList(confpool.REM,workingNode)
}

func (s *server) run() {
    var mess pairMex
    var timeoutElection,err = s._state.GetTimeoutNotifycationChan(raftstate.TIMER_ELECTION)

    if err != nil{
        log.Panicln(err)
    }

    for {
        select {
        case mess = <-s.messageChannel:
            log.Println("message received: ", mess.payload.ToString())
            go s.newMessageReceived(mess)
        case <- timeoutElection:
            log.Println("election not implemented")
        }
    }
}

func (s *server) newMessageReceived(mess pairMex){
            var rpcCall rpcs.Rpc
            var oldRole raftstate.Role
            var resp rpcs.Rpc
            var byEnc []byte
            var errEn error
            var senderNode node.Node 

            senderNode,errEn = s._state.GetNode(mess.sender)
            if errEn != nil {
                var v,f = s.clientList.Load(mess.sender)
                if !f{
                    log.Println(errEn)
                    return
                }
                senderNode = v.(node.Node)
            }
            log.Println("node founded: ", senderNode.GetIp())
            oldRole = s._state.GetRole()
            rpcCall = mess.payload
            resp = rpcCall.Execute(s._state, senderNode)
            log.Println("finih executing rpc: ",rpcCall.ToString())

            if resp != nil {
                log.Println("sending resp to caller RPC: ", resp.ToString())
                byEnc, errEn = genericmessage.Encode(resp)
                if errEn != nil{
                    log.Panicln("error encoding this rpc: ", resp.ToString())
                }
                senderNode.Send(byEnc)
            }

            if s._state.GetRole() == raftstate.LEADER && oldRole != state.LEADER {
            }
            mess.workdone <- 1
}

//utility
func (s *server) nodeAppendEntryPayload(n node.Node, toAppend []raft_log.LogInstance) rpcs.Rpc{
    panic("not implemented")
}

func (s *server) encodeAndSend(rpcMex rpcs.Rpc, n node.Node){
    var err error
    var rawMex []byte

    rawMex,err = genericmessage.Encode(rpcMex)
    if err != nil {
        log.Panicln("error encoding rpc: ", rpcMex.ToString())
    }

    err = n.Send(rawMex)

    if err != nil {
        log.Println("error sending rpcRawMex: ", err)
    }

}
