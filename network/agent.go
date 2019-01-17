package network

import "log"

// Agent handles reading from the connection
type Agent struct {
	id        string
	conn      Conn
	processor Processor
}

// Run Starts reading from the client (blocks until error from the read)
func (a *Agent) run() {
	for {
		data, err := a.conn.Read()
		if err != nil {
			log.Printf("Read message error: %v", err)
			break
		}

		msg, err := a.processor.Unmarshal(data)
		if err != nil {
			log.Printf("Unmarshal message error: %v", err)
			continue
		}

		err = a.processor.Route(a, msg)
		if err != nil {
			log.Printf("Routing error: %v", err)
		}
	}
}

func (a *Agent) ID() string {
	return a.id
}

func (a *Agent) Conn() Conn {
	return a.conn
}

func (a *Agent) Processor() Processor {
	return a.processor
}
