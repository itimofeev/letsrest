package letsrest

import (
	"errors"
	"sync"
)

type RequestStore interface {
	Save(*ClientRequest) (*ClientRequest, error)
	Get(id string) (*ClientRequest, error)
	Delete(id string) error

	SetResponse(id string, response *ClientResponse) error
	GetResponse(id string) (*ClientResponse, error)
}

func NewRequestStore() *MapRequestStore {
	return &MapRequestStore{store: make(map[string]*RequestData)}
}

type MapRequestStore struct {
	sync.RWMutex
	store map[string]*RequestData
}

func (s *MapRequestStore) Save(in *ClientRequest) (*ClientRequest, error) {
	s.Lock()
	defer s.Unlock()

	id, err := GenerateRandomString(10)
	in.ID = id
	s.store[id] = &RequestData{ID: id, Request: in}
	return in, err
}
func (s *MapRequestStore) Get(id string) (cReq *ClientRequest, err error) {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data != nil {
		cReq = data.Request
	}
	return
}

func (s *MapRequestStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
	return nil
}

func (s *MapRequestStore) SetResponse(id string, response *ClientResponse) error {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data == nil {
		return errors.New("request.not.found")
	}
	data.Response = response
	return nil
}

func (s *MapRequestStore) GetResponse(id string) (*ClientResponse, error) {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data == nil {
		return nil, nil
	}
	return data.Response, nil
}
