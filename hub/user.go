package hub

import "sync"

type UserID string
type users map[UserID]*User

type User struct {
	id UserID // never modify

	mu   sync.RWMutex
	room *Room
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) GetRoom() *Room {
	u.mu.RLock()
	defer u.mu.RUnlock()

	return u.room
}
