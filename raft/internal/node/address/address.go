package address

import (
	"strconv"
	"strings"
)

type NodeAddress interface {
	GetIp() string
	GetPort() string
}


func NewNodeAddress(ipAddr string, port string) (NodeAddress, error) {
	var sectorsStr = strings.Split(ipAddr, ".")
	var node nodeAddress


	for i := 0; i < sectorsNumber; i++ {
		out, err := strconv.Atoi(sectorsStr[i])
		if err != nil {
			return nil, err
		}
		node.sectors[i] = uint8(out)
	}

	var port_int, _ = strconv.Atoi(port)
	node.port = uint16(port_int)

	return &node, nil
}
