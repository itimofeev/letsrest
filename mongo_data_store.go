package letsrest

import (
	"github.com/kataras/iris/core/errors"
	"github.com/rs/xid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func NewMongoDataStore(wp WorkerPool) *MongoDataStore {
	mongoAddress := "192.168.99.100:27017"
	session, err := mgo.Dial(mongoAddress)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})

	ensureIndex(session)

	return &MongoDataStore{
		session: session,
		wp:      wp,
	}
}

type MongoDataStore struct {
	session *mgo.Session
	wp      WorkerPool
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("users")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func (s *MongoDataStore) PutUser(user *User) error {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("users")

	if err := c.Insert(user); err != nil {
		return err
	}
	return nil
}

func (s *MongoDataStore) GetUser(id string) (*User, error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("users")

	var user User
	if err := c.FindId(id).One(&user); err != nil {
		return nil, err
	}

	if user.ID == "" {
		return nil, nil
	}
	return &user, nil
}

func (s *MongoDataStore) GetRequest(id string) (*Request, error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	var request Request
	if err := c.FindId(id).One(&request); err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if request.ID == "" {
		return nil, nil
	}
	return &request, nil
}

func (s *MongoDataStore) CreateRequest(user *User, name string) (*Request, error) {
	request := &Request{
		ID:     s.generateId(),
		Name:   name,
		Status: &ExecStatus{Status: "idle"},
		UserID: user.ID,
	}

	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	if err := c.Insert(request); err != nil {
		return nil, err
	}

	return request, nil
}

func (s *MongoDataStore) generateId() string {
	return strings.Replace(xid.New().String(), "-", "", -1)
}

func (s *MongoDataStore) ExecRequest(id string, data *RequestData) (*Request, error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	var request Request
	if err := c.FindId(id).One(&request); err != nil {
		return nil, err
	}

	if request.ID == "" {
		return nil, errors.New("request.not.found")
	}

	request.RequestData = data
	request.Status.Status = "in_progress"
	request.Status.Error = ""

	if err := c.UpdateId(request.ID, request); err != nil {
		return nil, err
	}

	s.wp.AddRequest(&request, s)

	return &request, nil
}

func (s *MongoDataStore) CopyRequest(user *User, id string) (*Request, error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	var request Request
	if err := c.FindId(id).One(&request); err != nil {
		return nil, err
	}

	if request.ID == "" {
		return nil, errors.New("request.not.found")
	}

	newRequest := &Request{
		ID:     s.generateId(),
		Name:   request.Name,
		Status: &ExecStatus{Status: "idle"},
		UserID: user.ID,
	}
	if request.RequestData != nil {
		newRequest.RequestData = &*request.RequestData
	}

	if err := c.Insert(newRequest); err != nil {
		return nil, err
	}

	return newRequest, nil
}

func (s *MongoDataStore) List(user *User) (requests []*Request, err error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	if err := c.Find(bson.M{"user_id": user.ID}).All(&requests); err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *MongoDataStore) Delete(id string) error {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	if err := c.RemoveId(id); err != nil {
		return err
	}
	return nil
}

func (s *MongoDataStore) SetResponse(id string, response *Response, err error) error {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB("letsrest").C("requests")

	var request Request
	if err := c.FindId(id).One(&request); err != nil {
		return err
	}

	if request.ID == "" {
		return errors.New("request.not.found")
	}

	if err != nil {
		request.Status.Status = "error"
		request.Status.Error = err.Error()
	} else {
		request.Status.Status = "done"
	}
	request.Response = response

	if err := c.UpdateId(request.ID, request); err != nil {
		return err
	}

	return nil
}
