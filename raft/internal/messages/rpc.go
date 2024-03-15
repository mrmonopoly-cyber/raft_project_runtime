package messages

type Rpc interface {
  GetTerm() uint64
  Encode() ([]byte, error)
  Decode(b []byte) error
  ToString() string
  Manage()
}
