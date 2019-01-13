package network

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WSConn is a websocket connection
type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	writeCh   chan []byte
	closeFlag bool
}

func newWSConn(conn *websocket.Conn) *WSConn {
	wsConn := &WSConn{
		conn:    conn,
		writeCh: make(chan []byte, 10),
	}

	go func() {
		for bytes := range wsConn.writeCh {
			log.Printf("Sending text: %s\n", string(bytes))
			if err := conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				break
			}
		}

		conn.Close()
		wsConn.Lock()
		wsConn.closeFlag = true
		wsConn.Unlock()
		log.Println("Writer closed")
	}()

	return wsConn
}

func (wsConn *WSConn) Read() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

func (wsConn *WSConn) Write(b []byte) {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}

	wsConn.writeCh <- b
}

// Close sets closeFlag as true to prevent further writing to client
func (wsConn *WSConn) Close() {
	wsConn.Lock()
	if wsConn.closeFlag {
		return
	}
	wsConn.closeFlag = true
	wsConn.Unlock()

	close(wsConn.writeCh)
}
