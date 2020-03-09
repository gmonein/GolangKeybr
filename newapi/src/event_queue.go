package main

import (
	"sync"
)

const TypingError = 0
const TypingValid = 1
const Disconect = 2
const Finish = 3

type Event struct {
	UserID    int `json:"user_id"`
	EventType int `json:"event"`
	NextIndex int `json:"next_index"`
}

type EventQueue struct {
	Queue       []*Event
	QueueMutex  []sync.Mutex
	index       int
	indexMutex  sync.Mutex
	Length      int
	initialized bool
}

func (q *EventQueue) Initialize() {
	if q.initialized {
		return
	}
	q.Length = 10000
	q.Queue = make([]*Event, q.Length)
	q.QueueMutex = make([]sync.Mutex, q.Length)
	q.initialized = true
}

func (q *EventQueue) Push(event *Event) {
	if !q.initialized {
		q.Initialize()
	}
	q.indexMutex.Lock()

	q.index = q.nextIndex()
	nextIndex := q.nextIndex()

	q.QueueMutex[q.index].Lock()
	q.QueueMutex[nextIndex].Lock()

	q.Queue[q.index] = event
	q.Queue[nextIndex] = nil

	q.QueueMutex[nextIndex].Unlock()
	q.QueueMutex[q.index].Unlock()
	q.indexMutex.Unlock()
}

func (q *EventQueue) Get(index int) *Event {
	q.QueueMutex[index].Lock()
	res := q.Queue[index]
	q.QueueMutex[index].Unlock()
	return res
}

type QueueSubscription struct {
	queue *EventQueue
	index int
}

func (s *QueueSubscription) Next() *Event {
	nextIndex := s.nextIndex()
	next := s.queue.Get(nextIndex)
	if next != nil {
		s.index = s.nextIndex()
	}
	return next
}

func (s *QueueSubscription) nextIndex() int {
	return (s.index + 1) % s.queue.Length
}

func (q *EventQueue) nextIndex() int {
	return (q.index + 1) % q.Length
}
