package queue

import "errors"

var (
	// ErrQueueEmpty represents the error returned when a queue is empty.
	ErrQueueEmpty = errors.New("queue is empty")

	// ErrQueueFull represents the error returned when a queue is full.
	ErrQueueFull = errors.New("queue is full")
)

// ChannelQueue is a generic type that represents a queue implemented using channels.
// It can hold elements of any type.
type ChannelQueue[T any] struct {
	channel chan T
}

// NewQueue creates a new ChannelQueue with the specified capacity.
// The capacity determines the maximum number of items the queue can hold.
// The returned ChannelQueue will have an underlying channel with the given capacity.
func NewQueue[T any](capacity int) ChannelQueue[T] {
	return ChannelQueue[T]{
		channel: make(chan T, capacity),
	}
}

// Enqueue adds an item to the queue. If the queue is full, it returns an ErrQueueFull error.
func (q *ChannelQueue[T]) Enqueue(item T) error {
	select {
	case q.channel <- item:
		return nil
	default:
		return ErrQueueFull
	}
}

// Dequeue removes and returns an item from the queue.
// If the queue is empty, it returns a zero value of type T and an ErrQueueEmpty error.
func (q *ChannelQueue[T]) Dequeue() (T, error) {
	select {
	case item := <-q.channel:
		return item, nil
	default:
		var zeroVal T
		return zeroVal, ErrQueueEmpty
	}
}
