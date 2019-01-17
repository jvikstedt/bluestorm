package network

import (
	"fmt"
	"log"
	"sync"
)

// Agent handles reading from the connection
type Agent struct {
	id        string
	conn      Conn
	processor Processor

	muMetadata sync.RWMutex
	metadata   map[string]interface{}
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

func (a *Agent) WriteMsg(msg interface{}) error {
	b, err := a.processor.Marshal(msg)
	if err != nil {
		return err
	}
	a.conn.Write(b)

	return nil
}

func (a *Agent) SetValue(key string, value interface{}) {
	a.muMetadata.Lock()
	defer a.muMetadata.Unlock()

	a.metadata[key] = value
}

func (a *Agent) GetValue(key string) (interface{}, error) {
	a.muMetadata.RLock()
	defer a.muMetadata.RUnlock()

	value, ok := a.metadata[key]
	if !ok {
		return nil, fmt.Errorf("Could not find value by key %s", key)
	}

	return value, nil
}
