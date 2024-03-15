package address

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const sectors_number = 4

type Node_address interface
{
    Get_ip() string
    Send(mess []byte) error 
}

type node_address struct
{
	sectors [sectors_number]uint8
	port    uint16
    connection *net.Conn
}

type EnumType int

const (
  APPEND_ENTRY EnumType = iota
  REQUEST_VOTE 
  APPEND_RESPONSE
  VOTE_RESPONSE
  COPY_STATE
)

type message_typed struct
{
  Ty EnumType
  Payload []byte
}

func (this node_address) Get_ip() string {
    var ip_addr string = ""
    var num uint8 = this.sectors[0]
    
    for i:= 0;  i< sectors_number-1; i++ {
        num = this.sectors[i]
        ip_addr += strconv.Itoa(int(num))
        ip_addr += "."
    }
    ip_addr += strconv.Itoa(int(num))
    return ip_addr
}


func (this node_address) Send(mess []byte) error {
    mess_json,_ := json.Marshal(mess)
    log.Println(string(mess_json))
    // Access the value pointed to by the connection field
    conn := this.connection
    fmt.Fprintf(*conn, string(mess_json) + "\n")
    return nil
}

func New_node_address(ip_addr string, port uint16) Node_address{
	var sectors_str = strings.Split(ip_addr, ".")
	var node node_address

	for i := 0; i < sectors_number; i++ {
		out, err := strconv.Atoi(sectors_str[i])
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		node.sectors[i] = uint8(out)
	}

	node.port = port
    node.connection = nil

	return &node
}
