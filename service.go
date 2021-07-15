package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}

// Implementation of the Service interface.
type ServiceImpl struct {
	n uint64
	p time.Duration
}

// GetLimits method which returns defined restrictions: n - the number of batches, p - time period.
func (s ServiceImpl) GetLimits() (n uint64, p time.Duration) {
	return s.n, s.p
}

// Process method which processes deals with income batch.
func (s ServiceImpl) Process(ctx context.Context, batch Batch) error {
	if int(s.n) < len(batch) {
		return ErrBlocked
	}
	fmt.Println("Processed ", len(batch))
	<-ctx.Done()
	return nil
}

// Batch is a batch of items.
type Batch []Item

type ServiceParameters struct {
	// Maximum number of items in the batch.
	batchSize int
	// Time to process the batch.
	period time.Duration
	// The number of items to process.
	numberOfItems int
}

// Item is some abstract item.
type Item struct{}

func main() {
	params := ServiceParameters{
		batchSize: 25,
		period: time.Duration(2) * time.Second,
		numberOfItems: 69,
	}

	var s Service
	var batch Batch
	s = ServiceImpl{uint64(params.batchSize), params.period}

	// Define a slice of all items to process.
	allItems := make([]Item, params.numberOfItems)

	// Cursor to remember how many items have been already processed.
	cursor := 0

	// The number of batches we need to process.
	batchAmount := int(math.Ceil(float64(params.numberOfItems / params.batchSize)))

	for i := 0; i <= batchAmount; i++ {
		// Set lower bound for the number of items.
		lb := cursor

		// Set upper bound for the number of items within cursor.
		ub := cursor + params.batchSize

		// If upper bound exceeds the number of items, set it as the number of items.
		if ub > params.numberOfItems {
			ub = params.numberOfItems
		}

		// Slice for the batch with the number of items.
		batch = allItems[lb:ub]

		// Define the context with timeout.
		timeoutCtx, cancel := context.WithTimeout(context.Background(), params.period)
		defer cancel()
		err := s.Process(timeoutCtx, batch)
		if err == ErrBlocked {
			fmt.Println("Service is blocked")
		}
		cursor += params.batchSize
	}
}
