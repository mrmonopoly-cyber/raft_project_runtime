package address

import (
	"fmt"
	"raft/internal/node"
	"strconv"
	"strings"
)

const sectors_number = 4

type Node_address struct {
	sectors [sectors_number]uint8
	port    uint16
}

// Get_ip implements node.Node.
func (this *Node_address) Get_ip() string {
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


// send implements node.Node.
func (this *Node_address) Send() error {
	panic("not implemented")
}

func New_node_address(ip_addr string, port uint16) node.Node {
	var sectors_str = strings.Split(ip_addr, ".")
	var node Node_address

	for i := 0; i < sectors_number; i++ {
		out, err := strconv.Atoi(sectors_str[i])
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		node.sectors[i] = uint8(out)
	}

	node.port = port

	return &node
}
