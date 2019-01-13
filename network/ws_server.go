package network

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

type connSet map[*websocket.Conn]struct{}

// WSServer is a websocket server
type WSServer struct {
	Addr string
	Processor
	OnConnect    func(*Agent)
	OnDisconnect func(*Agent)
	ln           net.Listener

	conns      connSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	s.wg.Add(1)
	defer s.wg.Done()

	wsConn := s.addNewConnection(conn)

	agent := &Agent{
		ID:        GenerateID(),
		Conn:      wsConn,
		Processor: s.Processor,
	}
	if s.OnConnect != nil {
		s.OnConnect(agent)
	}
	agent.run()
	if s.OnDisconnect != nil {
		s.OnDisconnect(agent)
	}

	s.removeConnection(wsConn)
}

func (s *WSServer) addNewConnection(conn *websocket.Conn) *WSConn {
	log.Println("Adding new connection")

	s.mutexConns.Lock()
	defer s.mutexConns.Unlock()

	if s.conns == nil {
		s.conns = connSet{}
	}
	s.conns[conn] = struct{}{}

	return newWSConn(conn)
}

func (s *WSServer) removeConnection(wsConn *WSConn) {
	log.Println("Removing connection")
	wsConn.Close()

	s.mutexConns.Lock()
	defer s.mutexConns.Unlock()

	delete(s.conns, wsConn.conn)
	log.Println("Connection removed")
}

// Start starts a server (non blocking)
func (s *WSServer) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.ln = ln

	httpServer := &http.Server{
		Addr:    s.Addr,
		Handler: s,
	}

	log.Printf("Starting server %s\n", s.Addr)
	go httpServer.Serve(ln)

	return nil
}

// Stop stops a server, blocks until connections are closed
func (s *WSServer) Stop() {
	s.ln.Close()

	s.mutexConns.Lock()
	log.Printf("Closing %d connections\n", len(s.conns))
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
}
