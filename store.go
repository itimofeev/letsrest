package letsrest

import "sync"

type RequestStore interface {
	Save(*ClientRequest) (*ClientRequest, error)
	Get(id string) (*ClientRequest, error)
	Delete(id string) error
}

func NewRequestStore() *MapRequestStore {
	return &MapRequestStore{store: make(map[string]*ClientRequest)}
}

type MapRequestStore struct {
	sync.RWMutex
	store map[string]*ClientRequest
}

func (s *MapRequestStore) Save(in *ClientRequest) (*ClientRequest, error) {
	s.Lock()
	defer s.Unlock()

	id, err := GenerateRandomString(10)
	in.ID = id
	s.store[id] = in
	return in, err
}
func (s *MapRequestStore) Get(id string) (cReq *ClientRequest, err error) {
	s.RLock()
	defer s.RUnlock()

	cReq = s.store[id]
	return
}
func (s *MapRequestStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
	return nil
}
