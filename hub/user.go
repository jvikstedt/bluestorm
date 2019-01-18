package hub

import (
	"sync"

	"github.com/jvikstedt/bluestorm/network"
)

type UserID string
type Users map[UserID]*User

type User struct {
	id    UserID         // never modify
	agent *network.Agent // never modify

	mu   sync.RWMutex
	room Room
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) GetRoom() Room {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.room
}

func (u *User) Agent() *network.Agent {
	return u.agent
}
