package bluestorm

import (
	"log"
	"os"
	"os/signal"

	"github.com/jvikstedt/bluestorm/hub"
	"github.com/jvikstedt/bluestorm/network"
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

func CloseOnSignal(sig ...os.Signal) chan bool {
	closeSig := make(chan bool)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, sig...)
		<-sigint

		closeSig <- true
	}()
	return closeSig
}

func OnConnectHelper(hubManager *hub.Manager, defaultRoomID hub.RoomID) func(*network.Agent) {
	return func(agent *network.Agent) {
		log.Printf("Agent connectedted %s\n", agent.ID())
		defaultRoom, err := hubManager.GetRoom(defaultRoomID)
		if err != nil {
			log.Println(err)
			return
		}

		if err := hubManager.UserToRoom(hub.UserID(agent.ID()), defaultRoomID, agent); err != nil {
			log.Println(err)
			agent.Conn().Close()
			return
		}

		agent.SetValue("room", defaultRoom)
	}
}

func OnDisconnectHelper(hubManager *hub.Manager) func(*network.Agent) {
	return func(agent *network.Agent) {
		log.Printf("Agent disconnected %s\n", agent.ID())
		if err := hubManager.RemoveUser(hub.UserID(agent.ID())); err != nil {
			log.Println(err)
		}
	}
}
