package hub

import (
	"fmt"
	"sync"
)

type Manager struct {
	mu    sync.RWMutex
	users Users
	rooms
}

func NewManager() *Manager {
	return &Manager{
		users: make(Users),
		rooms: make(rooms),
	}
}

func (m *Manager) UserToRoom(uid UserID, rid RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[uid]
	if !ok {
		return fmt.Errorf("Could not move user %s to room %s, because user does not exist", uid, rid)
	}

	room, ok := m.rooms[rid]
	if !ok {
		return fmt.Errorf("Could not move user %s to room %s, because room does not exist", uid, rid)
	}

	if oldroom := user.GetRoom(); oldroom != nil {
		oldroom.deleteUser(uid)
	}

	user.setRoom(room)
	room.addUser(user)

	return nil
}

func (m *Manager) AddUser(user User, rid RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.users[user.ID()]
	if ok {
		return fmt.Errorf("Could not add user %s, it already exists", user.ID())
	}

	room, ok := m.rooms[rid]
	if !ok {
		return fmt.Errorf("Could not get room %s, it does not exist", rid)
	}

	user.setRoom(room)
	room.addUser(user)
	m.users[user.ID()] = user

	return nil
}

func (m *Manager) RemoveUser(uid UserID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[uid]
	if !ok {
		return fmt.Errorf("Could not remove user %s, user does not exist", uid)
	}

	if oldroom := user.GetRoom(); oldroom != nil {
		oldroom.deleteUser(uid)
	}

	delete(m.users, uid)

	return nil
}

func (m *Manager) AddRoom(room Room) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.rooms[room.ID()]
	if ok {
		return fmt.Errorf("Could not add room %s, it already exists", room.ID())
	}

	m.rooms[room.ID()] = room
	return nil
}

func (m *Manager) RemoveRoom(rid RoomID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.rooms[rid]
	if !ok {
		return fmt.Errorf("Could not remove room %s, it does not exist", rid)
	}

	// TODO What if room has users, should users be migrated somewhere else?

	delete(m.rooms, rid)

	return nil
}

func (m *Manager) GetRoom(rid RoomID) (Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[rid]
	if !ok {
		return nil, fmt.Errorf("Could not get room %s, it does not exist", rid)
	}

	return room, nil
}

func (m *Manager) GetUser(uid UserID) (User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[uid]
	if !ok {
		return nil, fmt.Errorf("Could not get user %s, it does not exist", uid)
	}

	return user, nil
}

var defaultManager = NewManager()

func DefaultManager() *Manager {
	return defaultManager
}
