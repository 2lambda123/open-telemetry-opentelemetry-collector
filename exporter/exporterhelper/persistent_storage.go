// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exporterhelper

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"

	"go.opencensus.io/metric/metricdata"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/extension/storage"
)

// persistentStorage provides an interface for request storage operations
type persistentStorage interface {
	// put appends the request to the storage
	put(req request) error
	// get returns the next available request; note that the channel is unbuffered
	get() <-chan request
	// size returns the current size of the persistent storage with items waiting for processing
	size() uint64
	// stop gracefully stops the storage
	stop()
}

// persistentContiguousStorage provides a persistent queue implementation backed by file storage extension
//
// Write index describes the position at which next item is going to be stored.
// Read index describes which item needs to be read next.
// When Write index = Read index, no elements are in the queue.
//
// The items currently processed by consumers are not deleted until the processing is finished.
// Their list is stored under a separate key.
//
//
//   ┌───────file extension-backed queue───────┐
//   │                                         │
//   │     ┌───┐     ┌───┐ ┌───┐ ┌───┐ ┌───┐   │
//   │ n+1 │ n │ ... │ 4 │ │ 3 │ │ 2 │ │ 1 │   │
//   │     └───┘     └───┘ └─x─┘ └─|─┘ └─x─┘   │
//   │                       x     |     x     │
//   └───────────────────────x─────|─────x─────┘
//      ▲              ▲     x     |     x
//      │              │     x     |     xxxx deleted
//      │              │     x     |
//    write          read    x     └── currently processed item
//    index          index   x
//                           xxxx deleted
//
type persistentContiguousStorage struct {
	logger      *zap.Logger
	queueName   string
	client      storage.Client
	unmarshaler requestUnmarshaler

	putChan  chan struct{}
	stopChan chan struct{}
	stopOnce sync.Once
	capacity uint64

	reqChan chan request

	mu                      sync.Mutex
	readIndex               itemIndex
	writeIndex              itemIndex
	currentlyProcessedItems []itemIndex

	itemsCount uint64
}

type itemIndex uint64

const (
	zapKey           = "key"
	zapQueueNameKey  = "queueName"
	zapErrorCount    = "errorCount"
	zapNumberOfItems = "numberOfItems"

	readIndexKey               = "ri"
	writeIndexKey              = "wi"
	currentlyProcessedItemsKey = "pi"
)

var (
	errMaxCapacityReached   = errors.New("max capacity reached")
	errValueNotSet          = errors.New("value not set")
	errKeyNotPresentInBatch = errors.New("key was not present in get batchStruct")
)

// newPersistentContiguousStorage creates a new file-storage extension backed queue;
// queueName parameter must be a unique value that identifies the queue.
// The queue needs to be initialized separately using initPersistentContiguousStorage.
func newPersistentContiguousStorage(ctx context.Context, queueName string, capacity uint64, logger *zap.Logger, client storage.Client, unmarshaler requestUnmarshaler) *persistentContiguousStorage {
	pcs := &persistentContiguousStorage{
		logger:      logger,
		client:      client,
		queueName:   queueName,
		unmarshaler: unmarshaler,
		capacity:    capacity,
		putChan:     make(chan struct{}, capacity),
		reqChan:     make(chan request),
		stopChan:    make(chan struct{}),
	}

	initPersistentContiguousStorage(ctx, pcs)
	unprocessedReqs := pcs.retrieveUnprocessedItems(context.Background())

	err := currentlyProcessedBatchesGauge.UpsertEntry(func() int64 {
		return int64(pcs.numberOfCurrentlyProcessedItems())
	}, metricdata.NewLabelValue(pcs.queueName))
	if err != nil {
		logger.Error("failed to create number of currently processed items metric", zap.Error(err))
	}

	// We start the loop first so in case there are more elements in the persistent storage than the capacity,
	// it does not get blocked on initialization

	go pcs.loop()

	// Make sure the leftover requests are handled
	pcs.enqueueUnprocessedReqs(unprocessedReqs)
	// Make sure the communication channel is loaded up
	for i := uint64(0); i < pcs.size(); i++ {
		pcs.putChan <- struct{}{}
	}

	return pcs
}

func initPersistentContiguousStorage(ctx context.Context, pcs *persistentContiguousStorage) {
	var writeIndex itemIndex
	var readIndex itemIndex
	batch, err := newBatch(pcs).get(readIndexKey, writeIndexKey).execute(ctx)

	if err == nil {
		readIndex, err = batch.getItemIndexResult(readIndexKey)
	}

	if err == nil {
		writeIndex, err = batch.getItemIndexResult(writeIndexKey)
	}

	if err != nil {
		pcs.logger.Error("failed getting read/write index, starting with new ones",
			zap.String(zapQueueNameKey, pcs.queueName),
			zap.Error(err))
		pcs.readIndex = 0
		pcs.writeIndex = 0
	} else {
		pcs.readIndex = readIndex
		pcs.writeIndex = writeIndex
	}

	atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))
}

func (pcs *persistentContiguousStorage) enqueueUnprocessedReqs(reqs []request) {
	if len(reqs) > 0 {
		errCount := 0
		for _, req := range reqs {
			if pcs.put(req) != nil {
				errCount++
			}
		}
		if errCount > 0 {
			pcs.logger.Error("errors occurred while moving items for processing back to queue",
				zap.String(zapQueueNameKey, pcs.queueName),
				zap.Int(zapNumberOfItems, len(reqs)), zap.Int(zapErrorCount, errCount))

		} else {
			pcs.logger.Info("moved items for processing back to queue",
				zap.String(zapQueueNameKey, pcs.queueName),
				zap.Int(zapNumberOfItems, len(reqs)))

		}
	}
}

// loop is the main loop that handles fetching items from the persistent buffer
func (pcs *persistentContiguousStorage) loop() {
	for {
		select {
		case <-pcs.stopChan:
			return
		case <-pcs.putChan:
			req, found := pcs.getNextItem(context.Background())
			if found {
				pcs.reqChan <- req
			}
		}
	}
}

// get returns the request channel that all the requests will be send on
func (pcs *persistentContiguousStorage) get() <-chan request {
	return pcs.reqChan
}

// size returns the number of currently available items, which were not picked by consumers yet
func (pcs *persistentContiguousStorage) size() uint64 {
	return atomic.LoadUint64(&pcs.itemsCount)
}

// numberOfCurrentlyProcessedItems returns the count of batches for which processing started but hasn't finish yet
func (pcs *persistentContiguousStorage) numberOfCurrentlyProcessedItems() int {
	pcs.mu.Lock()
	defer pcs.mu.Unlock()
	return len(pcs.currentlyProcessedItems)
}

func (pcs *persistentContiguousStorage) stop() {
	pcs.logger.Debug("stopping persistentContiguousStorage", zap.String(zapQueueNameKey, pcs.queueName))
	pcs.stopOnce.Do(func() {
		close(pcs.stopChan)
		_ = currentlyProcessedBatchesGauge.UpsertEntry(func() int64 {
			return int64(pcs.numberOfCurrentlyProcessedItems())
		}, metricdata.NewLabelValue(pcs.queueName))
	})
}

// put marshals the request and puts it into the persistent queue
func (pcs *persistentContiguousStorage) put(req request) error {
	// Nil requests are ignored
	if req == nil {
		return nil
	}

	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	if pcs.size() >= pcs.capacity {
		pcs.logger.Warn("maximum queue capacity reached", zap.String(zapQueueNameKey, pcs.queueName))
		return errMaxCapacityReached
	}

	itemKey := pcs.itemKey(pcs.writeIndex)
	pcs.writeIndex++
	atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))

	ctx := context.Background()
	_, err := newBatch(pcs).setItemIndex(writeIndexKey, pcs.writeIndex).setRequest(itemKey, req).execute(ctx)

	// Inform the loop that there's some data to process
	pcs.putChan <- struct{}{}

	return err
}

// getNextItem pulls the next available item from the persistent storage; if none is found, returns (nil, false)
func (pcs *persistentContiguousStorage) getNextItem(ctx context.Context) (request, bool) {
	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	if pcs.readIndex != pcs.writeIndex {
		index := pcs.readIndex
		// Increase here, so even if errors happen below, it always iterates
		pcs.readIndex++
		atomic.StoreUint64(&pcs.itemsCount, uint64(pcs.writeIndex-pcs.readIndex))

		pcs.updateReadIndex(ctx)
		pcs.itemProcessingStart(ctx, index)

		batch, err := newBatch(pcs).get(pcs.itemKey(index)).execute(ctx)
		if err != nil {
			return nil, false
		}

		req, err := batch.getRequestResult(pcs.itemKey(index))
		if err != nil || req == nil {
			return nil, false
		}

		req.setOnProcessingFinished(func() {
			pcs.mu.Lock()
			defer pcs.mu.Unlock()
			pcs.itemProcessingFinish(ctx, index)
		})
		return req, true
	}

	return nil, false
}

// retrieveUnprocessedItems gets the items for which processing was not finished, cleans the storage
// and moves the items back to the queue
func (pcs *persistentContiguousStorage) retrieveUnprocessedItems(ctx context.Context) []request {
	var reqs []request
	var processedItems []itemIndex

	pcs.mu.Lock()
	defer pcs.mu.Unlock()

	pcs.logger.Debug("checking if there are items left by consumers", zap.String(zapQueueNameKey, pcs.queueName))
	batch, err := newBatch(pcs).get(currentlyProcessedItemsKey).execute(ctx)
	if err == nil {
		processedItems, err = batch.getItemIndexArrayResult(currentlyProcessedItemsKey)
	}
	if err != nil {
		pcs.logger.Error("could not fetch items left by consumers", zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
		return reqs
	}

	if len(processedItems) > 0 {
		pcs.logger.Info("fetching items left for processing by consumers",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Int(zapNumberOfItems, len(processedItems)))
	} else {
		pcs.logger.Debug("no items left for processing by consumers")
	}

	reqs = make([]request, len(processedItems))
	keys := make([]string, len(processedItems))
	retrieveBatch := newBatch(pcs)
	cleanupBatch := newBatch(pcs)
	for i, it := range processedItems {
		keys[i] = pcs.itemKey(it)
		retrieveBatch.get(keys[i])
		cleanupBatch.delete(keys[i])
	}

	_, retrieveErr := retrieveBatch.execute(ctx)
	_, cleanupErr := cleanupBatch.execute(ctx)

	if retrieveErr != nil {
		pcs.logger.Warn("failed retrieving items left by consumers", zap.String(zapQueueNameKey, pcs.queueName), zap.Error(retrieveErr))
	}

	if cleanupErr != nil {
		pcs.logger.Debug("failed cleaning items left by consumers", zap.String(zapQueueNameKey, pcs.queueName), zap.Error(cleanupErr))
	}

	if retrieveErr != nil {
		return reqs
	}

	for i, key := range keys {
		req, err := retrieveBatch.getRequestResult(key)
		if err != nil {
			pcs.logger.Warn("failed unmarshalling item",
				zap.String(zapQueueNameKey, pcs.queueName), zap.String(zapKey, key), zap.Error(err))
		} else {
			reqs[i] = req
		}
	}

	return reqs
}

// itemProcessingStart appends the item to the list of currently processed items
func (pcs *persistentContiguousStorage) itemProcessingStart(ctx context.Context, index itemIndex) {
	pcs.currentlyProcessedItems = append(pcs.currentlyProcessedItems, index)
	_, err := newBatch(pcs).
		setItemIndexArray(currentlyProcessedItemsKey, pcs.currentlyProcessedItems).
		execute(ctx)
	if err != nil {
		pcs.logger.Debug("failed updating currently processed items",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

// itemProcessingFinish removes the item from the list of currently processed items and deletes it from the persistent queue
func (pcs *persistentContiguousStorage) itemProcessingFinish(ctx context.Context, index itemIndex) {
	var updatedProcessedItems []itemIndex
	for _, it := range pcs.currentlyProcessedItems {
		if it != index {
			updatedProcessedItems = append(updatedProcessedItems, it)
		}
	}
	pcs.currentlyProcessedItems = updatedProcessedItems

	_, err := newBatch(pcs).
		setItemIndexArray(currentlyProcessedItemsKey, pcs.currentlyProcessedItems).
		delete(pcs.itemKey(index)).
		execute(ctx)
	if err != nil {
		pcs.logger.Debug("failed updating currently processed items",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

func (pcs *persistentContiguousStorage) updateReadIndex(ctx context.Context) {
	_, err := newBatch(pcs).
		setItemIndex(readIndexKey, pcs.readIndex).
		execute(ctx)

	if err != nil {
		pcs.logger.Debug("failed updating read index",
			zap.String(zapQueueNameKey, pcs.queueName), zap.Error(err))
	}
}

func (pcs *persistentContiguousStorage) itemKey(index itemIndex) string {
	return strconv.FormatUint(uint64(index), 10)
}
