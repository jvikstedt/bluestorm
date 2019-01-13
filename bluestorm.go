package bluestorm

import (
	"log"
)

// Server specifies required functions for the server
type Server interface {
	Start() error // Should not block
	Stop()
}

// Servers is collection of Server
type Servers []Server

// Run runs and stops servers
// You can trigger stop by passing value to closeSig channel
func Run(closeSig chan bool, servers Servers) {
	for _, server := range servers {
		err := server.Start()
		if err != nil {
			log.Println(err)
		}
	}
	<-closeSig
	for i, server := range servers {
		log.Printf("Stopping server %d/%d\n", i+1, len(servers))
		server.Stop()
	}
	log.Println("All servers stopped!")
}
