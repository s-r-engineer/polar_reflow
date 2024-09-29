package linker

import (
	"sync"
)

type Linker struct {
	chain *Link
	mutex sync.Mutex
}

type Link struct {
	value interface{}
	next  *Link
	prev  *Link
}

func (l *Linker) lockMe() func() {
	l.mutex.Lock()
	return func() {
		l.mutex.Unlock()
	}
}

func (l *Linker) isEmpty() bool {
	return l.chain == nil
}

func (l *Linker) count() (counter int) {
	defer l.lockMe()()
	for piece := l.chain.next; piece != l.chain; piece = piece.next {
		counter += 1
	}
	return counter + 1
}

func (l *Linker) Push(value interface{}) {
	defer l.lockMe()()
	link := Link{value: value}
	if l.isEmpty() {
		link.next = &link
		link.prev = &link
	} else {
		link.next = l.chain
		link.prev = l.chain.prev
		l.chain.prev = &link
		link.prev.next = &link
	}
	l.chain = &link
}

// will return value, delete the link from the chain and return the function to add link back if some shit will happen
func (l *Linker) Pop() (value interface{}, f func()) {
	defer l.lockMe()()
	value = l.chain.value
	l.chain.next.prev = l.chain.prev
	l.chain.prev.next = l.chain.next
	l.chain = l.chain.next
	return value, func() {
		l.Push(value)
	}
}

func CreateLinker() *Linker { return &Linker{} }
