package storage

import "sync"

type UserCache struct {
	mut sync.RWMutex
	users map[int]string
}

func NewUserCache() UserCache {
	return UserCache{
		mut:   sync.RWMutex{},
		users: make(map[int]string),
	}
}

func (uc *UserCache) GetUser(user int) (string, bool) {
	uc.mut.RLock()
	defer uc.mut.RUnlock()
	if uc.users == nil {
		uc.users = make(map[int]string)
	}
	u, ok := uc.users[user]
	return u, ok
}

func (uc *UserCache) SetUser( userID int, username string) {
	uc.mut.RLock()
	defer uc.mut.RUnlock()
	if uc.users == nil {
		uc.users = make(map[int]string)
	}
	uc.users[userID] =  username
}