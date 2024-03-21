package address

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const sectorsNumber = 4

type NodeAddress interface {
	GetIp() string
	GetPort() string
}

type nodeAddress struct {
	sectors [sectorsNumber]uint8
	port    uint16
}

func (this nodeAddress) GetIp() string {
	var ipAddr string = ""
	var num uint8 = this.sectors[0]

	for i := 0; i < sectorsNumber-1; i++ {
		num = this.sectors[i]
		ipAddr += strconv.Itoa(int(num))
		ipAddr += "."
	}
	ipAddr += strconv.Itoa(int(num))
	return ipAddr
}

func (this nodeAddress) GetPort() string {
	return strconv.Itoa(int(this.port))
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
