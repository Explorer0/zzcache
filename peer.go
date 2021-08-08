package zzcache

type Peer interface {
	Get(group string, key string) ([]byte, error)
	Set(group string, key string, value []byte)	error
	GetName() string
}


