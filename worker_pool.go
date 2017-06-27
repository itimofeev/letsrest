package letsrest

type CanSetResponse interface {
	SetResponse(id string, response *Response, err error) error
}

type requestWithStore struct {
	request *Request
	r       CanSetResponse
}

type WorkerPool interface {
	AddRequest(*Request, CanSetResponse)
}

func NewWorkerPool(requester Requester) *ChanWorkerPool {
	pool := &ChanWorkerPool{
		requester: requester,
		requestCh: make(chan *requestWithStore, 1000),
	}
	go pool.ListenForTasks()
	return pool
}

type ChanWorkerPool struct {
	requester Requester

	requestCh chan *requestWithStore
}

func (wp *ChanWorkerPool) AddRequest(request *Request, r CanSetResponse) {
	wp.requestCh <- &requestWithStore{request: request, r: r}
}

func (wp *ChanWorkerPool) ListenForTasks() {
	for rs := range wp.requestCh {
		resp, err := wp.requester.Do(rs.request.RequestData)
		rs.r.SetResponse(rs.request.ID, resp, err)
	}
}
