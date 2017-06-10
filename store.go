package letsrest

import (
	"errors"
	"github.com/speps/go-hashids"
	"sync"
)

type RequestStore interface {
	CreateRequest(name string) (*Request, error)
	ExecRequest(id string, req *RequestData) (*Request, error)
	Get(id string) (*Request, error)
	Delete(id string) error
	List() ([]*Request, error)

	SetResponse(id string, response *Response, err error) error
}

func NewRequestStore(requester Requester) *MapRequestStore {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 20

	store := &MapRequestStore{
		store:     make(map[string]*Request),
		hd:        hd,
		requester: requester,
		requestCh: make(chan *Request, 1000),
	}
	go store.ListenForTasks()
	return store
}

type MapRequestStore struct {
	sync.RWMutex
	store     map[string]*Request
	requests  []*Request
	hd        *hashids.HashIDData
	requester Requester
	requestCh chan (*Request)
}

func (s *MapRequestStore) ListenForTasks() {
	for request := range s.requestCh {
		resp, err := s.requester.Do(request.RequestData)
		s.SetResponse(request.ID, resp, err)
	}
}

func (s *MapRequestStore) CreateRequest(name string) (*Request, error) {
	id, err := s.generateId()
	Must(err, "s.generateId()")
	request := &Request{ID: id, Name: name, Status: &ExecStatus{Status: "idle"}}

	s.Lock()
	defer s.Unlock()

	s.store[id] = request
	s.requests = append(s.requests, request)
	return request, err
}

func (s *MapRequestStore) generateId() (string, error) {
	h := hashids.NewWithData(s.hd)
	return h.Encode([]int{len(s.store)})
}

func (s *MapRequestStore) ExecRequest(id string, data *RequestData) (*Request, error) {
	s.Lock()
	defer s.Unlock()
	if request, ok := s.store[id]; ok {
		request.RequestData = data
		request.Status.Status = "in_progress"
		request.Status.Error = ""
		s.requestCh <- request
		return request, nil
	}

	return nil, errors.New("request not found")
}

func (s *MapRequestStore) Get(id string) (request *Request, err error) {
	s.RLock()
	defer s.RUnlock()

	if data, ok := s.store[id]; ok {
		return data, nil
	}
	return nil, nil
}

func (s *MapRequestStore) List() (requests []*Request, err error) {
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
	s.Lock()
	defer s.Unlock()

	request, ok := s.store[id]
	if !ok {
		return errors.New("request.not.found")
	}
	if err != nil {
		request.Status.Status = "error"
		request.Status.Error = err.Error()
	} else {
		request.Status.Status = "done"
	}
	request.Response = response
	return nil
}
