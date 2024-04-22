package rpcs

type Rpc interface {
	ToString() string
	//Execute() *Rpc
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
