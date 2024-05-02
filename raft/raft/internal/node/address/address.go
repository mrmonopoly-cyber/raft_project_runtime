package address

import (
	"fmt"
	"strconv"
	"strings"
)

type NodeAddress interface {
	GetIp() string
	GetPort() string
}


func NewNodeAddress(ipAddr string, port string) NodeAddress {
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

	return &node
}
