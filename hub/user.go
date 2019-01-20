package hub

import (
	"sync"
)

type User interface {
	ID() UserID
	setRoom(r Room)
	GetRoom() Room
	WriteMsg(i interface{}) error
}

type UserID string
type Users map[UserID]User

type BaseUser struct {
	id UserID // never modify

	mu   sync.RWMutex
	room Room
}

func NewBaseUser(uid UserID) *BaseUser {
	return &BaseUser{
		id: uid,
	}
}

func (u *BaseUser) ID() UserID {
	return u.id
}

func (u *BaseUser) GetRoom() Room {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.room
}

func (u *BaseUser) setRoom(r Room) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.room = r
}

func (u *BaseUser) WriteMsg(i interface{}) error {
	panic("not implemented")
}
