package network

// Conn presents a connection between server and a client
type Conn interface {
	Read() ([]byte, error)
	Write(b []byte)
	Close()
}
