package letsrest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapRequestStore_CRUD(t *testing.T) {
	store := NewRequestStore()

	bucket, _ := store.CreateBucket("somename")
	assert.NotEmpty(t, bucket.ID)

	loaded, _ := store.Get(bucket.ID)
	assert.Equal(t, bucket, loaded)

	list, _ := store.List()
	assert.Len(t, list, 1)

	notExisted, _ := store.Get("someNotExistentID")
	assert.Nil(t, notExisted)

	store.Delete(bucket.ID)
	deleted, _ := store.Get(bucket.ID)
	assert.Nil(t, deleted)

	list, _ = store.List()
	assert.Len(t, list, 0)
}
