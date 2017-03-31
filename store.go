package letsrest

import (
	"errors"
	"github.com/speps/go-hashids"
	"sync"
)

type RequestStore interface {
	Save(*RequestTask) (*RequestTask, error)
	Get(id string) (*RequestTask, error)
	Delete(id string) error
	List() ([]RequestTask, error)

	SetResponse(id string, response *Response, err error) error
	GetResponse(id string) (*Result, error)
}

func NewRequestStore() *MapRequestStore {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 20

	return &MapRequestStore{store: make(map[string]*RequestData), hd: hd}
}

// объект для хранения в store
type RequestData struct {
	ID string

	Info     *Info
	Request  *RequestTask
	Response *Response
}

type MapRequestStore struct {
	sync.RWMutex
	store    map[string]*RequestData
	taskList []RequestData
	hd       *hashids.HashIDData
}

func (s *MapRequestStore) Save(in *RequestTask) (*RequestTask, error) {
	s.Lock()
	defer s.Unlock()

	id, err := s.generateId()
	in.ID = id
	rd := RequestData{ID: id, Request: in, Info: &Info{Status: "in_progress"}}
	s.store[id] = &RequestData{ID: id, Request: in, Info: &Info{Status: "in_progress"}}
	s.taskList = append(s.taskList, rd)
	return in, err
}

func (s *MapRequestStore) generateId() (string, error) {
	h := hashids.NewWithData(s.hd)
	return h.Encode([]int{len(s.store)})
}

func (s *MapRequestStore) Get(id string) (cReq *RequestTask, err error) {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data != nil {
		cReq = data.Request
	}
	return
}
func (s *MapRequestStore) List() (taskList []RequestTask, err error) {
	s.RLock()
	defer s.RUnlock()

	taskList = make([]RequestTask, len(s.taskList))
	for i, value := range s.taskList {
		taskList[len(s.taskList)-i-1] = *value.Request
	}

	return
}

func (s *MapRequestStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
	return nil
}

func (s *MapRequestStore) SetResponse(id string, response *Response, err error) error {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data == nil {
		return errors.New("request.not.found")
	}
	if err != nil {
		data.Info.Status = "error"
		data.Info.Error = err.Error()
	} else {
		data.Info.Status = "done"
	}
	data.Response = response
	return nil
}

func (s *MapRequestStore) GetResponse(id string) (*Result, error) {
	s.RLock()
	defer s.RUnlock()

	data := s.store[id]
	if data == nil {
		return nil, nil
	}
	return &Result{Response: data.Response, Status: data.Info}, nil
}
