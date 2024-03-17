package address

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
  "bufio"
)

const sectorsNumber = 4

type NodeAddress interface
{
    GetIp() string
    GetPort() string
    Send(mess []byte) error
    Receive() ([]byte, error) 
    HandleConnOut(conn *net.Conn)
    HandleConnIn(conn *net.Conn)
}

type nodeAddress struct
{
	sectors [sectorsNumber] uint8
	port    uint16
  connIn *net.Conn
  connOut *net.Conn
}

type EnumType int

const (
  APPEND_ENTRY EnumType = iota
  REQUEST_VOTE 
  APPEND_RESPONSE
  VOTE_RESPONSE
  COPY_STATE
)

type messageTyped struct
{
  Ty EnumType
  Payload []byte
}

func (this nodeAddress) GetIp() string {
    var ipAddr string = ""
    var num uint8 = this.sectors[0]
    
    for i:= 0;  i< sectorsNumber-1; i++ {
        num = this.sectors[i]
        ipAddr += strconv.Itoa(int(num))
        ipAddr += "."
    }
    ipAddr += strconv.Itoa(int(num))
    return ipAddr
}

func (this nodeAddress) GetPort() string {
    return string(rune(this.port))
}

func (this nodeAddress) Send(mess []byte) error {
    mess_json,_ := json.Marshal(mess)
    log.Println(string(mess_json))
    // Access the value pointed to by the connection field
    conn := this.connOut
    fmt.Fprintf(*conn, string(mess_json) + "\n")
    return nil
}

func (this nodeAddress) Receive() ([]byte, error) {
  mess, err := bufio.NewReader(*this.connIn).ReadString('\n')
  return []byte(mess), err
} 

func (this nodeAddress) HandleConnOut(conn *net.Conn) {
  this.connOut = conn
}

func (this nodeAddress) HandleConnIn(conn *net.Conn) {
  this.connIn = conn
}

func NewNodeAddress(ipAddr string, port string) NodeAddress{
    var sectorsStr = strings.Split(ipAddr, ".")
    var node nodeAddress

    for i := 0; i < sectorsNumber; i++ {
        out, err := strconv.Atoi(sectorsStr[i])
        if err != nil {
            fmt.Println("Error:", err)
            return nil
        }
        node.sectors[i] = uint8(out)
    }

    var port_int, _ = strconv.Atoi(port)
    node.port = uint16(port_int)
    node.connOut = nil
    node.connIn = nil

    return &node
}
