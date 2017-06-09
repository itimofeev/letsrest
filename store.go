package letsrest

import (
	"errors"
	"github.com/speps/go-hashids"
	"sync"
)

type RequestStore interface {
	CreateBucket(name string) (*Bucket, error)
	CreateRequest(*Bucket, *Request) error
	Get(id string) (*Bucket, error)
	Delete(id string) error
	List() ([]Bucket, error)

	SetResponse(id string, response *Response, err error) error
}

func NewRequestStore() *MapRequestStore {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 20

	return &MapRequestStore{store: make(map[string]Bucket), hd: hd}
}

type MapRequestStore struct {
	sync.RWMutex
	store   map[string]Bucket
	buckets []Bucket
	hd      *hashids.HashIDData
}

func (s *MapRequestStore) CreateBucket(name string) (*Bucket, error) {
	id, err := s.generateId()
	Must(err, "s.generateId()")
	bucket := Bucket{ID: id, Status: &ExecStatus{Status: "in_progress"}}

	s.Lock()
	defer s.Unlock()

	s.store[id] = bucket
	s.buckets = append(s.buckets, bucket)
	return &bucket, err
}

func (s *MapRequestStore) generateId() (string, error) {
	h := hashids.NewWithData(s.hd)
	return h.Encode([]int{len(s.store)})
}

func (s *MapRequestStore) CreateRequest(*Bucket, *Request) error {
	return nil
}

func (s *MapRequestStore) Get(id string) (bucket *Bucket, err error) {
	s.RLock()
	defer s.RUnlock()

	if data, ok := s.store[id]; ok {
		return &data, nil
	}
	return nil, nil
}

func (s *MapRequestStore) List() (taskList []Bucket, err error) {
	s.RLock()
	defer s.RUnlock()

	return s.buckets[:], nil
}

func (s *MapRequestStore) Delete(id string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.store, id)
	for i := range s.buckets {
		if s.buckets[i].ID == id {
			s.buckets = append(s.buckets[:i], s.buckets[i+1:]...)
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
