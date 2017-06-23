package letsrest

import (
	"errors"
	"github.com/speps/go-hashids"
	"sync"
)

type DataStore interface {
	GetRequest(id string) (*Request, error)
	CreateRequest(user *User, id string) (*Request, error)
	ExecRequest(id string, data *RequestData) (*Request, error)
	List(user *User) (requests []*Request, err error)
	Delete(id string) error
	SetResponse(id string, response *Response, err error) error

	PutUser(*User) error
	GetUser(id string) (*User, error)
}

func NewDataStore(r Requester) DataStore {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 20

	store := &MapDataStore{
		requester:      r,
		requestsByUser: make(map[string][]*Request),
		requests:       make(map[string]*Request),
		users:          make(map[string]*User),
		requestCh:      make(chan *Request, 1000),
		hd:             hd,
	}

	go store.ListenForTasks()

	return store
}

type MapDataStore struct {
	sync.RWMutex // protecting maps

	hd *hashids.HashIDData

	requestCh chan *Request

	requestsByUser map[string][]*Request
	requests       map[string]*Request
	users          map[string]*User
	requester      Requester
}

func (s *MapDataStore) ListenForTasks() {
	for request := range s.requestCh {
		resp, err := s.requester.Do(request.RequestData)
		s.SetResponse(request.ID, resp, err)
	}
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

func (s *MapDataStore) GetRequest(id string) (*Request, error) {
	s.RLock()
	defer s.RUnlock()

	request := s.requests[id]
	return request, nil
}

func (s *MapDataStore) CreateRequest(user *User, name string) (*Request, error) {
	id, err := s.generateId()
	Must(err, "s.generateId()")
	request := &Request{
		ID:     id,
		Name:   name,
		Status: &ExecStatus{Status: "idle"},
		UserID: user.ID,
	}

	s.Lock()
	defer s.Unlock()

	s.requests[id] = request
	s.requestsByUser[user.ID] = append(s.requestsByUser[user.ID], request)
	return request, err
}

func (s *MapDataStore) generateId() (string, error) {
	h := hashids.NewWithData(s.hd)
	return h.Encode([]int{len(s.requests)})
}

func (s *MapDataStore) ExecRequest(id string, data *RequestData) (*Request, error) {
	s.Lock()
	defer s.Unlock()
	if request, ok := s.requests[id]; ok {
		request.RequestData = data
		request.Status.Status = "in_progress"
		request.Status.Error = ""
		s.requestCh <- request
		return request, nil
	}

	return nil, errors.New("request not found")
}

func (s *MapDataStore) List(user *User) (requests []*Request, err error) {
	if len(s.requestsByUser[user.ID]) == 0 {
		return make([]*Request, 0), nil
	}

	s.RLock()
	defer s.RUnlock()
	return s.requestsByUser[user.ID][:], nil
}

func (s *MapDataStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	request, ok := s.requests[id]
	if !ok {
		return nil
	}

	delete(s.requests, id)
	userRequests := s.requestsByUser[request.UserID]
	for i := range userRequests {
		if userRequests[i].ID == id {
			s.requestsByUser[request.UserID] = append(userRequests[:i], userRequests[i+1:]...)
		}
	}

	return nil
}

func (s *MapDataStore) SetResponse(id string, response *Response, err error) error {
	s.Lock()
	defer s.Unlock()

	request, ok := s.requests[id]
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
