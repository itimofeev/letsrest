package letsrest

import (
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type fakeWp struct {
}

func (wp *fakeWp) AddRequest(*Request, CanSetResponse) {

}

var mgStore = NewMongoDataStore(&fakeWp{})

func Test_MongoDataStore_CreateUser(t *testing.T) {
	user := createTestUser(t)

	requests, err := mgStore.List(user)
	require.Nil(t, err)
	assert.Len(t, requests, 0)
}

func Test_MongoDataStore_CRUDRequest(t *testing.T) {
	user := createTestUser(t)

	request, err := mgStore.CreateRequest(user, "my-cool-name")
	require.Nil(t, err)
	assert.Equal(t, "my-cool-name", request.Name)

	getRequest, err := mgStore.GetRequest(request.ID)
	require.Nil(t, err)
	assert.Equal(t, request, getRequest)

	requests, err := mgStore.List(user)
	require.Nil(t, err)
	assert.Len(t, requests, 1)
	assert.Equal(t, request, requests[0])

	assert.Nil(t, mgStore.Delete(request.ID))
	notExists, err := mgStore.GetRequest(request.ID)
	require.Nil(t, err)
	assert.Nil(t, notExists)
}

func Test_MongoDataStore_ExecRequest(t *testing.T) {
	user := createTestUser(t)

	request, err := mgStore.CreateRequest(user, "my-cool-name")
	require.Nil(t, err)
	assert.Equal(t, "my-cool-name", request.Name)

	data := &RequestData{URL: "someUrl", Method: "someMethod", Headers: []Header{{Name: "headerName", Value: "HeaderValue"}}}
	_, err = mgStore.ExecRequest(request.ID, data)
	require.Nil(t, err)

	getRequest, err := mgStore.GetRequest(request.ID)
	require.Nil(t, err)
	require.NotNil(t, getRequest)

	assert.Equal(t, getRequest.RequestData, data)
	assert.Equal(t, getRequest.Status.Status, "in_progress")
}

func Test_MongoDataStore_SetResponse(t *testing.T) {
	user := createTestUser(t)

	request, err := mgStore.CreateRequest(user, "my-cool-name")
	require.Nil(t, err)
	assert.Equal(t, "my-cool-name", request.Name)

	data := &RequestData{URL: "someUrl", Method: "someMethod", Headers: []Header{{Name: "headerName", Value: "HeaderValue"}}}
	_, err = mgStore.ExecRequest(request.ID, data)
	require.Nil(t, err)

	response := &Response{
		ContentType: "ct",
		Duration:    time.Second,
		StatusCode:  777,
		Headers:     []Header{{Name: "123", Value: "777"}},
	}
	require.Nil(t, mgStore.SetResponse(request.ID, response, nil))

	getRequest, err := mgStore.GetRequest(request.ID)
	require.Nil(t, err)
	require.NotNil(t, getRequest)

	assert.Equal(t, getRequest.Response, response)
	assert.Equal(t, "done", getRequest.Status.Status)
}

func createTestUser(t *testing.T) *User {
	require.NotNil(t, mgStore)

	userId := xid.New().String()

	assert.Nil(t, mgStore.PutUser(&User{ID: userId}))

	user, err := mgStore.GetUser(userId)
	require.Nil(t, err)
	require.NotNil(t, user)
	assert.Equal(t, userId, user.ID)

	return user
}
