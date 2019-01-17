package hub

import (
	"fmt"
	"sync"
)

type RoomID string
type rooms map[RoomID]*Room

type Room struct {
	id RoomID // never modify

	mu sync.RWMutex
	users
}

func (r *Room) ID() RoomID {
	return r.id
}

func (r *Room) GetUsers() users {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copyUsers := make(users)
	for k, v := range r.users {
		copyUsers[k] = v
	}

	return copyUsers
}

func (r *Room) GetUserByID(uid UserID) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[uid]
	if !ok {
		return nil, fmt.Errorf("Could not find user %s from room %s", uid, r.id)
	}

	return u, nil
}
