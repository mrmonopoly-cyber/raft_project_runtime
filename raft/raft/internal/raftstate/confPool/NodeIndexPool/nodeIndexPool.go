package nodeIndexPool

type OP uint8
const (
    ADD OP = iota
    REM OP = iota
)

type INFO_OP uint8
const(
    GET INFO_OP = iota
    SET INFO_OP = iota
)

type INFO_INDEX uint8
const(
    MATCH INFO_INDEX = iota
    NEXTT INFO_INDEX = iota
)

type NodeIndexPool interface{
    UpdateStatusList(op OP, ip string)
    UpdateNodeInfo(op INFO_OP, info INFO_INDEX, ip string) (int,error)
}

func NewLeaederCommonIdx() NodeIndexPool{
    return nil
}
