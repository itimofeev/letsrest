package letsrest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapRequestStore_CRUD(t *testing.T) {
	store := NewDataStore(&testRequester{})
	user1 := &User{ID: "1"}
	store.PutUser(user1)

	request, _ := store.CreateRequest(user1, "somename")
	assert.NotEmpty(t, request.ID)

	loaded, _ := store.GetRequest(request.ID)
	assert.Equal(t, request, loaded)

	list, _ := store.List(user1)
	assert.Len(t, list, 1)

	notExisted, _ := store.Get("someNotExistentID")
	assert.Nil(t, notExisted)

	store.Delete(request.ID)
	deleted, _ := store.Get(request.ID)
	assert.Nil(t, deleted)

	list, _ = store.List(user1)
	assert.Len(t, list, 0)
}
