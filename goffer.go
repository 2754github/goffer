// Package goffer provides an interface to pub/sub based buffers.
package goffer

import (
	"errors"
	"sync"
)

var (
	// ErrNotSubscribed indicates that the Buffer is not subscribed.
	ErrNotSubscribed = errors.New("buffer is not subscribed")

	// ErrClosed indicates that the Buffer is closed.
	ErrClosed = errors.New("buffer is closed")
)

// Buffer is a channel-like object that allows
// publishing a single item at a time and
// subscribing to multiple items in chunks.
type Buffer[Item any] interface {
	// Publish enqueues a single item into the Buffer.
	//
	// Publish enqueues and returns nil or returns an error without enqueuing.
	Publish(item Item) error

	// Subscribe dequeues all items from the Buffer. (Each time the Buffer fills up.)
	Subscribe() <-chan []Item

	// Pull dequeues all items from the Buffer. (Even if the Buffer is closed.)
	Pull() []Item

	// Close closes the Buffer.
	Close()
}

type buffer[Item any] struct {
	mu              sync.Mutex
	size            int
	items           []Item
	subscriber      chan []Item
	isNotSubscribed bool
	isClosed        bool
}

// New returns an empty Buffer of Item with size.
//
// If size <= 1, the Buffer size will be 1.
func New[Item any](size int) Buffer[Item] {
	s := max(size, 1)

	return &buffer[Item]{
		mu:              sync.Mutex{},
		size:            s,
		items:           make([]Item, 0, s),
		subscriber:      make(chan []Item, 1),
		isNotSubscribed: true,
		isClosed:        false,
	}
}

func (b *buffer[Item]) Publish(item Item) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isClosed {
		return ErrClosed
	}

	if b.isNotSubscribed {
		return ErrNotSubscribed
	}

	b.items = append(b.items, item)
	if len(b.items) == b.size {
		b.subscriber <- b.items
		b.items = make([]Item, 0, b.size)
	}

	return nil
}

func (b *buffer[Item]) Subscribe() <-chan []Item {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.isNotSubscribed = false

	return b.subscriber
}

func (b *buffer[Item]) Pull() []Item {
	b.mu.Lock()
	defer b.mu.Unlock()

	r := make([]Item, len(b.items))
	copy(r, b.items)
	b.items = make([]Item, 0, b.size)

	return r
}

func (b *buffer[Item]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isClosed {
		return
	}

	b.isClosed = true

	close(b.subscriber)
}
