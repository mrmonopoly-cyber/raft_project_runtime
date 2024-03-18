package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	messages "raft/internal/messages"
    cutom_mex "raft/internal/node/message"
	"reflect"
	"raft/internal/node"
	state "raft/internal/raftstate"
	"sync"
)

type Server struct {
	_state         state.State
	messageChannel chan messages.Rpc
	otherNodes     *sync.Map
	listener       net.Listener
	wg             sync.WaitGroup
}

func generateID(input string) string {
    // Create a new SHA256 hasher
    hasher := sha256.New()

    // Write the input string to the hasher
    hasher.Write([]byte(input))

    // Get the hashed bytes
    hashedBytes := hasher.Sum(nil)

    // Convert the hashed bytes to a hexadecimal string
    id := hex.EncodeToString(hashedBytes)

    return id
}


func NewServer(term uint64,ip_addr string, port string, serversIp []string) *Server {
	listener, err := net.Listen("tcp", ip_addr + ":" + port)

	if err != nil {
		log.Fatalf("Failed to listen on port %s: %s", port, err)
	}

	var server = &Server{
		_state:         state.NewState(term, port, state.FOLLOWER),
		otherNodes:     &sync.Map{},
		messageChannel: make(chan messages.Rpc),
		listener:       listener,
	}

    for i :=0; i< len(serversIp) -2; i++{
        var new_node node.Node
        var err error
        new_node,err= node.NewNode(serversIp[i],port)
        if err != nil {
            println("error in creating new node: %v", err)
            continue
        }
        server.otherNodes.Store(generateID(serversIp[i]),new_node)
    }
	return server
}

func (s *Server) Start() {
	s.wg.Add(3)

	go s.acceptIncomingConn()

	s.connectToServers()

	s._state.StartElectionTimeout()

	if s._state.Leader() {
		s._state.StartHearthbeatTimeout()
	}

	go s.run()

	go s.handleResponse()

	s.wg.Wait()
}

/*
 * Create a connection between this server to all the others and populate the map containing these connections
 */
func (s *Server) connectToServers() {
    s.otherNodes.Range(func (key any, value interface{}) bool{
        var nodeEle node.Node
        var errEl bool
        nodeEle,errEl = value.(node.Node)
        if !errEl {
            log.Println("invalid object in otherNodes map: ", reflect.TypeOf(nodeEle))
            return false
        }
        var ipAddr string = nodeEle.GetIp()
        var port string = nodeEle.GetPort()
		var conn, err = net.Dial("tcp",ipAddr + ":" + port)
        for err != nil {
            log.Println("Dial error: ", err)
            return false
        }
        if conn != nil {
            nodeEle.AddConnOut(&conn)
        }
        return true
    })
}

func (s *Server) acceptIncomingConn() {
	defer s.wg.Done()
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
        var id_node string = generateID(newConncetionIp)
        var value , found  = s.otherNodes.Load(id_node)

        if found {
            var connectedNode node.Node = value.(node.Node)
            connectedNode.AddConnIn(&conn)
        }else {
            var new_node,_= node.NewNode(newConncetionIp, newConncetionPort)        
            new_node.AddConnIn(&conn)
            s.otherNodes.Store(id_node,new_node)
        }
	}
}

/*
 * Handle incoming messages and send it to the corresponding channel
 */
func (s *Server) handleResponse() {
	defer s.wg.Done()
	// iterating over the connections map and receive byte message
	for {
		s.otherNodes.Range(func(k, conn interface{}) bool {
            var message string
            var errMes error
			message, errMes = conn.(node.Node).Recv()
            if errMes != nil {
                fmt.Printf("error in reading from node %v with error %v", 
                    conn.(node.Node).GetIp(), errMes)
                return false
            }
			s.messageChannel <- *cutom_mex.NewMessage([]byte(message)).ToRpc()
			return true
		})
	}
}


func (s *Server) run() {
  defer s.wg.Done()
  for {
    var mess messages.Rpc
    select {
    case mess = <- s.messageChannel:
      mess.Execute(s.otherNodes, &s._state)

    case <- s._state.HeartbeatTimeout().C:
      node.SendAll(s.otherNodes, mex []byte)  
    case <- s._state.ElectionTimeout().C:
    }
  }
}
