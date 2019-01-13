package network

import "log"

// Agent handles reading from the connection
type Agent struct {
	ID string
	Conn
	Processor Processor
}

// Run Starts reading from the client (blocks until error from the read)
func (a *Agent) run() {
	for {
		data, err := a.Conn.Read()
		if err != nil {
			log.Printf("Read message error: %v", err)
			break
		}

		msg, err := a.Processor.Unmarshal(data)
		if err != nil {
			log.Printf("Unmarshal message error: %v", err)
			continue
		}

		err = a.Processor.Route(a, msg)
		if err != nil {
			log.Printf("Routing error: %v", err)
		}
	}
}
