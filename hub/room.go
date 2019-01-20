package hub

import (
	"fmt"
	"log"
	"sync"
)

type Room interface {
	ID() RoomID
	deleteUser(uid UserID)
	addUser(user User)
}

type RoomID string
type rooms map[RoomID]Room

type BaseRoom struct {
	id RoomID // never modify

	mu    sync.RWMutex
	users Users
}

func NewBaseRoom(id RoomID) *BaseRoom {
	return &BaseRoom{
		id:    id,
		users: make(Users),
	}
}

func (r *BaseRoom) ID() RoomID {
	return r.id
}

func (r *BaseRoom) deleteUser(uid UserID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, uid)
}

func (r *BaseRoom) addUser(user User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.users == nil {
		r.users = make(Users)
	}
	r.users[user.ID()] = user
}

func (r *BaseRoom) GetUsers() Users {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copyUsers := make(Users)
	for k, v := range r.users {
		copyUsers[k] = v
	}

	return copyUsers
}

func (r *BaseRoom) GetUsersWithRead(callback func(Users)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	callback(r.users)
}

func (r *BaseRoom) GetUsersWithReadWrite(callback func(Users)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	callback(r.users)
}

func (r *BaseRoom) GetUserByID(uid UserID) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[uid]
	if !ok {
		return nil, fmt.Errorf("Could not find user %s from room %s", uid, r.id)
	}

	return u, nil
}

func (r *BaseRoom) Broadcast(msg interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if err := user.WriteMsg(msg); err != nil {
			log.Println(err)
		}
	}
}

func (r *BaseRoom) BroadcastExceptOne(uid UserID, msg interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID() == uid {
			continue
		}
		if err := user.WriteMsg(msg); err != nil {
			log.Println(err)
		}
	}
}
