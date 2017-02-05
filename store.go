package letsrest

type RequestStore interface {
	Save(*ClientRequest) (*ClientRequest, error)
	Get(id string) (*ClientRequest, error)
	Delete(id string) error
}

type MapRequestStore struct {
	store map[string]*ClientRequest
}

func (s *MapRequestStore) Save(in *ClientRequest) (*ClientRequest, error) {
	return in, nil
}
func (s *MapRequestStore) Get(id string) (*ClientRequest, error) {
	return nil, nil
}
func (s *MapRequestStore) Delete(id string) error {
	return nil
}
