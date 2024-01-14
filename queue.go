package sync_queue

import (
	"sync"

	"gopkg.in/eapache/queue.v1"
)

type SyncQueue struct {
	buffer *queue.Queue
	cond   *sync.Cond
	closed bool
}

func NewSyncQueue() *SyncQueue {
	return &SyncQueue{
		buffer: queue.New(),
		cond:   sync.NewCond(&sync.Mutex{}),
		closed: false,
	}
}

func (q *SyncQueue) Pop() (v interface{}) {
	buffer := q.buffer

	q.cond.L.Lock()

	for buffer.Length() == 0 && !q.closed {
		q.cond.Wait()
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}

	q.cond.L.Unlock()
	return
}

func (q *SyncQueue) TryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.cond.L.Lock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if q.closed {
		ok = true
	}

	q.cond.L.Unlock()
	return
}

func (q *SyncQueue) Push(v interface{}) {
	q.cond.L.Lock()

	if !q.closed {
		q.buffer.Add(v)
		q.cond.Signal()
	}

	q.cond.L.Unlock()
}

func (q *SyncQueue) Len() (l int) {
	q.cond.L.Lock()
	l = q.buffer.Length()
	q.cond.L.Unlock()
	return
}

func (q *SyncQueue) Close() {
	q.cond.L.Lock()
	if !q.closed {
		q.closed = true
		q.cond.Broadcast()
	}
	q.cond.L.Unlock()
}
