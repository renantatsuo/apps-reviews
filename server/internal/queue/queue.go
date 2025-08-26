package queue

import (
	"github.com/mattdeak/gopq"
)

// Queue is the interface for the queue
type Queue interface {
	// Enqueue an item to the queue
	Enqueue(item []byte) error
	// Dequeue an item from the queue
	Dequeue() ([]byte, error)
	// Close the queue
	Close() error
}

type queue struct {
	queue *gopq.Queue
}

func New(connStr string) Queue {
	return &queue{
		queue: connect(connStr),
	}
}

func connect(connStr string) *gopq.Queue {
	queue, err := gopq.NewSimpleQueue(connStr)
	if err != nil {
		panic(err)
	}
	return queue
}

func (q *queue) Enqueue(item []byte) error {
	return q.queue.Enqueue(item)
}

func (q *queue) Dequeue() ([]byte, error) {
	msg, err := q.queue.Dequeue()
	if err != nil {
		return nil, err
	}
	return msg.Item, nil
}

func (q *queue) Close() error {
	return q.queue.Close()
}
