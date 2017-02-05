package letsrest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapRequestStore_CRUD(t *testing.T) {
	store := NewRequestStore()
	cReq := &RequestTask{}

	saved, _ := store.Save(cReq)
	assert.NotEmpty(t, saved.ID)

	loaded, _ := store.Get(cReq.ID)
	assert.Equal(t, saved, loaded)

	notExisted, _ := store.Get("someNotExistentID")
	assert.Nil(t, notExisted)

	store.Delete(cReq.ID)
	deleted, _ := store.Get(cReq.ID)
	assert.Nil(t, deleted)
}
