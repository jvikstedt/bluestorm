package hub

import (
	"fmt"
	"sync"

	"github.com/jvikstedt/bluestorm/network"
)

var DefaultRoomID RoomID = "default"

type Manager struct {
	mu    sync.RWMutex
	users Users
	rooms
}

func NewManager() *Manager {
	manager := &Manager{
		users: make(Users),
		rooms: make(rooms),
	}

	manager.rooms[DefaultRoomID] = &BaseRoom{
		id:    DefaultRoomID,
		users: make(Users),
	}

	return manager
}

func (m *Manager) AddRoom(room Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.rooms[room.ID()]
	if ok {
		return fmt.Errorf("Could not create room %s, it already exists", room.ID())
	}

	m.rooms[room.ID()] = room
	return nil
}

func (m *Manager) RemoveRoom(id RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.rooms[id]
	if !ok {
		return fmt.Errorf("Could not delete room %s, it does not exist", id)
	}

	delete(m.rooms, id)

	return nil
}

func (m *Manager) GetRoom(id RoomID) (Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[id]
	if !ok {
		return nil, fmt.Errorf("Could not get room %s, it does not exist", id)
	}

	return room, nil
}

// UserToRoom move user to room
// agent has to be set only on initial join, otherwise leave as nil
func (m *Manager) UserToRoom(uid UserID, rid RoomID, agent *network.Agent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[uid]
	if !ok {
		user = &User{
			id:    uid,
			agent: agent,
		}
		m.users[uid] = user
	}

	room, ok := m.rooms[rid]
	if !ok {
		room = &BaseRoom{
			id:    rid,
			users: make(Users),
		}
		m.rooms[rid] = room
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	if user.room == room {
		return fmt.Errorf("Could not move user %s to room %s, because user already is in that room", uid, rid)
	}

	if user.room != nil {
		// Remove user from old room
		user.room.DeleteUser(uid)
	}

	// Add user to new room
	user.room = room
	user.room.AddUser(user)

	return nil
}

func (m *Manager) RemoveUser(uid UserID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[uid]
	if !ok {
		return fmt.Errorf("Could not remove user %s, user does not exist", uid)
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	if user.room != nil {
		user.room.DeleteUser(uid)
	}

	delete(m.users, uid)

	return nil
}

var defaultManager = NewManager()

func DefaultManager() *Manager {
	return defaultManager
}
