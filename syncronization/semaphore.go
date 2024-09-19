package syncronization

import (
	"context"

	"golang.org/x/sync/semaphore"
)

func CreateSemaphoreInstance(n int) (func() error, func()) {
	newSemaphore := semaphore.NewWeighted(int64(n))
	return func() error { return newSemaphore.Acquire(context.Background(), int64(1)) }, func() { newSemaphore.Release(int64(1)) }
}
