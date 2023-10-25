// Copyright The OpenTelemetry Authors
// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/collector/exporter/exporterhelper/internal"

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/collector/component"
	intrequest "go.opentelemetry.io/collector/exporter/exporterhelper/internal/request"
	"go.opentelemetry.io/collector/exporter/exporterhelper/queue"
)

// boundedMemoryQueue implements a producer-consumer exchange similar to a ring buffer queue,
// where the queue is bounded and if it fills up due to slow consumers, the new items written by
// the producer are dropped.
type boundedMemoryQueue struct {
	stopWG       sync.WaitGroup
	size         *atomic.Uint32
	stopped      *atomic.Bool
	items        chan *intrequest.Request
	capacity     uint32
	numConsumers int
	callback     func(item *intrequest.Request)
}

// NewBoundedMemoryQueue constructs the new queue of specified capacity, and with an optional
// callback for dropped items (e.g. useful to emit metrics).
func NewBoundedMemoryQueue(capacity int, numConsumers int, callback func(item *intrequest.Request)) queue.Queue {
	return &boundedMemoryQueue{
		items:        make(chan *intrequest.Request, capacity),
		stopped:      &atomic.Bool{},
		size:         &atomic.Uint32{},
		capacity:     uint32(capacity),
		numConsumers: numConsumers,
		callback:     callback,
	}
}

// StartConsumers starts a given number of goroutines consuming items from the queue
// and passing them into the consumer callback.
func (q *boundedMemoryQueue) Start(_ context.Context, _ component.Host) error {
	var startWG sync.WaitGroup
	for i := 0; i < q.numConsumers; i++ {
		q.stopWG.Add(1)
		startWG.Add(1)
		go func() {
			startWG.Done()
			defer q.stopWG.Done()
			for item := range q.items {
				q.size.Add(^uint32(0))
				q.callback(item)
			}
		}()
	}
	startWG.Wait()
	return nil
}

// Produce is used by the producer to submit new item to the queue. Returns false in case of queue overflow.
func (q *boundedMemoryQueue) Produce(item *intrequest.Request) bool {
	if q.stopped.Load() {
		return false
	}

	// we might have two concurrent backing queues at the moment
	// their combined size is stored in q.size, and their combined capacity
	// should match the capacity of the new queue
	if q.size.Load() >= q.capacity {
		return false
	}

	q.size.Add(1)
	select {
	case q.items <- item:
		return true
	default:
		// should not happen, as overflows should have been captured earlier
		q.size.Add(^uint32(0))
		return false
	}
}

// Stop stops all consumers, as well as the length reporter if started,
// and releases the items channel. It blocks until all consumers have stopped.
func (q *boundedMemoryQueue) Stop() {
	q.stopped.Store(true) // disable producer
	close(q.items)
	q.stopWG.Wait()
}

// Size returns the current size of the queue
func (q *boundedMemoryQueue) Size() int {
	return int(q.size.Load())
}

func (q *boundedMemoryQueue) Capacity() int {
	return int(q.capacity)
}

func (q *boundedMemoryQueue) IsPersistent() bool {
	return false
}
