package letsrest

import (
	"errors"
	"github.com/speps/go-hashids"
	"sync"
)

type RequestStore interface {
	CreateRequest(name string) (*Request, error)
	ExecRequest(id string, req *RequestData) error
	Get(id string) (*Request, error)
	Delete(id string) error
	List() ([]Request, error)

	SetResponse(id string, response *Response, err error) error
}

func NewRequestStore() *MapRequestStore {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 20

	return &MapRequestStore{store: make(map[string]*Request), hd: hd}
}

type MapRequestStore struct {
	sync.RWMutex
	store    map[string]*Request
	requests []Request
	hd       *hashids.HashIDData
}

func (s *MapRequestStore) CreateRequest(name string) (*Request, error) {
	id, err := s.generateId()
	Must(err, "s.generateId()")
	request := Request{ID: id, Status: &ExecStatus{Status: "idle"}}

	s.Lock()
	defer s.Unlock()

	s.store[id] = &request
	s.requests = append(s.requests, request)
	return &request, err
}

func (s *MapRequestStore) generateId() (string, error) {
	h := hashids.NewWithData(s.hd)
	return h.Encode([]int{len(s.store)})
}

func (s *MapRequestStore) ExecRequest(id string, r *RequestData) error {
	s.RLock()
	defer s.RUnlock()
	if data, ok := s.store[id]; ok {
		data.RequestData = r
		return nil
	}

	return errors.New("request not found")
}

func (s *MapRequestStore) Get(id string) (request *Request, err error) {
	s.RLock()
	defer s.RUnlock()

	if data, ok := s.store[id]; ok {
		return data, nil
	}
	return nil, nil
}

func (s *MapRequestStore) List() (requests []Request, err error) {
	s.RLock()
	defer s.RUnlock()

	return s.requests[:], nil
}

func (s *MapRequestStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
	for i := range s.requests {
		if s.requests[i].ID == id {
			s.requests = append(s.requests[:i], s.requests[i+1:]...)
		}
	}

	return nil
}

func (s *MapRequestStore) SetResponse(id string, response *Response, err error) error {
	s.RLock()
	defer s.RUnlock()

	data, ok := s.store[id]
	if !ok {
		return errors.New("request.not.found")
	}
	if err != nil {
		data.Status.Status = "error"
		data.Status.Error = err.Error()
	} else {
		data.Status.Status = "done"
	}
	data.Response = response
	return nil
}
