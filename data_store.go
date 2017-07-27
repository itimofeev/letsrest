package letsrest

import (
	"errors"
	"github.com/rs/xid"
	"strings"
	"sync"
)
//интерфейс для работы с дс
type DataStore interface {
	CanSetResponse

	GetRequest(id string) (*Request, error)
	CreateRequest(user *User, name string) (*Request, error)
	EditRequest(id, name string) (*Request, error)
	ExecRequest(id string, data *RequestData) (*Request, error)
	CopyRequest(user *User, id string) (*Request, error)
	List(user *User) (requests []*Request, err error)
	Delete(id string) error

	PutUser(*User) error
	GetUser(id string) (*User, error)
}

func NewDataStore(config *Config, wp WorkerPool) DataStore {
	return NewMongoDataStore(config, wp)
}

func NewMapDataStore(wp WorkerPool) *MapDataStore {
	store := &MapDataStore{
		wp:             wp,
		requestsByUser: make(map[string][]*Request),
		requests:       make(map[string]*Request),
		users:          make(map[string]*User),
	}

	return store
}
//хранятся данные в памяти
type MapDataStore struct {
	sync.RWMutex // protecting maps

	requestsByUser map[string][]*Request
	requests       map[string]*Request
	users          map[string]*User
	wp             WorkerPool
}
//реализация методов
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
	request := &Request{
		ID:     s.generateId(),
		Name:   name,
		Status: &ExecStatus{Status: "idle"},
		UserID: user.ID,
	}

	s.Lock()
	defer s.Unlock()

	s.requests[request.ID] = request
	s.requestsByUser[user.ID] = append(s.requestsByUser[user.ID], request)
	return request, nil
}

func (s *MapDataStore) generateId() string {
	return strings.Replace(xid.New().String(), "-", "", -1)
}

func (s *MapDataStore) ExecRequest(id string, data *RequestData) (*Request, error) {
	s.Lock()
	defer s.Unlock()
	if request, ok := s.requests[id]; ok {
		request.RequestData = data
		request.Status.Status = "in_progress"
		request.Status.Error = ""
		s.wp.AddRequest(request, s)
		return request, nil
	}

	return nil, errors.New("request not found")
}

func (s *MapDataStore) EditRequest(id, name string) (*Request, error) {
	s.Lock()
	defer s.Unlock()
	if request, ok := s.requests[id]; ok {
		request.Name = name
		return request, nil
	}

	return nil, errors.New("request not found")
}

func (s *MapDataStore) CopyRequest(user *User, id string) (*Request, error) {
	s.Lock()
	defer s.Unlock()

	if request, ok := s.requests[id]; ok {
		newRequest := &Request{
			ID:     s.generateId(),
			Name:   request.Name,
			Status: &ExecStatus{Status: "idle"},
			UserID: user.ID,
		}
		if request.RequestData != nil {
			newRequest.RequestData = &*request.RequestData
		}
		s.requests[newRequest.ID] = newRequest
		s.requestsByUser[user.ID] = append(s.requestsByUser[user.ID], newRequest)
		return newRequest, nil
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
