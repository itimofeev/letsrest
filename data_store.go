package letsrest

import "sync"

type DataStore interface {
	RequestStore(user *User) RequestStore

	PutUser(*User) error
	GetUser(id string) (*User, error)
}

func NewDataStore(r Requester) DataStore {
	return &MapDataStore{
		r:            r,
		requestStore: make(map[string]RequestStore),
		users:        make(map[string]*User),
	}
}

type MapDataStore struct {
	sync.RWMutex // protecting maps

	requestStore map[string]RequestStore
	users        map[string]*User
	r            Requester
}

func (s *MapDataStore) RequestStore(user *User) RequestStore {
	s.RLock()
	if store, ok := s.requestStore[user.ID]; ok {
		s.RUnlock()
		return store
	}
	s.RUnlock()

	s.Lock()
	defer s.Unlock()

	store := NewRequestStore(s.r)
	s.requestStore[user.ID] = store

	return store
}

func (s *MapDataStore) PutUser(user *User) error {
	s.Lock()
	defer s.Unlock()

	s.users[user.ID] = user
	return nil
}

func (s *MapDataStore) GetUser(id string) (*User, error) {
	s.RLock()
	defer s.RUnlock()

	if user, ok := s.users[id]; ok {
		return user, nil
	}
	return nil, nil
}
