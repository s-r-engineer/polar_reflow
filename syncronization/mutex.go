package syncronization

import "sync"

func CreateMutexInstance() (func(), func()) {
	var mutex sync.Mutex

	return func() {
			mutex.Lock()
		}, func() {
			mutex.Unlock()
		}
}
