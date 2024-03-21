package address

import (
	"fmt"
	"strconv"
	"strings"
)

const sectorsNumber = 4

type NodeAddress interface
{
    GetIp() string
    GetPort() string
}

type nodeAddress struct
{
    sectors [sectorsNumber] uint8
    port    uint16
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

    return &node
}
