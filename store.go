package letsrest

type RequestStore interface {
	Save(*ClientRequest) (string, error)
	Get(id string) (*ClientRequest, error)
	Delete(id string) error
}

type MapRequestStore struct {
	store map[string]*ClientRequest
}

func (s *MapRequestStore) Save(*ClientRequest) (string, error) {
	return "", nil
}
func (s *MapRequestStore) Get(id string) (*ClientRequest, error) {
	return nil, nil
}
func (s *MapRequestStore) Delete(id string) error {
	return nil
}
