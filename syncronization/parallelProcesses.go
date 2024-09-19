package syncronization

import "sync"

func CreateWGInstance() (func(), func(), func()) {
	var parallelShit sync.WaitGroup
	return func() {
		parallelShit.Add(1)
	}, func() { parallelShit.Done() }, func() { parallelShit.Wait() }
}
