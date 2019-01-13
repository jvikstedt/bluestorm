package network

// Processor is used to handle network data
type Processor interface {
	Route(a *Agent, msg interface{}) error
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte) (interface{}, error)
}
