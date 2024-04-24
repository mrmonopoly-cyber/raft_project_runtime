package address

import (
	"strconv"
)

const sectorsNumber = 4

type nodeAddress struct {
	sectors [sectorsNumber]uint8
	port    uint16
}

func (this nodeAddress) GetIp() string {
	var ipAddr string = ""
	var num uint8 = this.sectors[0]
    var i =0

	for i = 0; i < sectorsNumber-1; i++ {
		num = this.sectors[i]
		ipAddr += strconv.Itoa(int(num))
		ipAddr += "."
	}
	num = this.sectors[i]
	ipAddr += strconv.Itoa(int(num))
	return ipAddr
}

func (this nodeAddress) GetPort() string {
	return strconv.Itoa(int(this.port))
}

