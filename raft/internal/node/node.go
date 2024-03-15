package node

type Node interface
{
    Send() error
    Get_ip() string
}
