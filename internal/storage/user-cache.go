package storage

import "sync"

type UserCache struct {
	mut sync.RWMutex
	users map[string]string
}

func NewUserCache() UserCache {
	return UserCache{
		mut:   sync.RWMutex{},
		users: make(map[string]string),
	}
}

func (uc *UserCache) GetUser(user string) (string, bool) {
	uc.mut.RLock()
	defer uc.mut.RUnlock()
	if uc.users == nil {
		uc.users = make(map[string]string)
	}
	u, ok := uc.users[user]
	return u, ok
}

func (uc *UserCache) SetUser(username string, userID string) {
	uc.mut.RLock()
	defer uc.mut.RUnlock()
	if uc.users == nil {
		uc.users = make(map[string]string)
	}
	uc.users[userID] =  username
}