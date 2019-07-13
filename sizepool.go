// Package sizepool implement a fix size pool of struct for Go
package sizepool

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

const (
	DEFAULT_POOLSIZE int64 = 1024
)

var (
	ErrNoEnoughItem = errors.New("no enough item in size pool")
)

// define the size pool entity
type sizePool struct {
	mu       sync.RWMutex
	initsize int64
	new      func() interface{}
	reset    func(interface{})
	pool     *list.List
}

// init a new size pool and return it
func NewPool(size int64, new func() interface{}, reset func(interface{})) *sizePool {
	if size <= 0 {
		size = DEFAULT_POOLSIZE
	}
	sp := &sizePool{
		mu:       sync.RWMutex{},
		initsize: size,
		new:      new,
		reset:    reset,
		pool:     list.New(),
	}
	sp.constructNewItem()
	return sp
}

// construct new item and return it to channel
func (p *sizePool) constructNewItem() {
	var (
		ch = make(chan interface{}, p.initsize)
	)
	for i := 0; int64(i) < p.initsize; i++ {
		go func() {
			ch <- p.new()
		}()
	}
	for i := 0; int64(i) < p.initsize; i++ {
		p.pool.PushBack(<-ch)
	}
}

// get size pool init size
func (p *sizePool) InitSize() int64 {
	p.mu.RLock()
	s := p.initsize
	p.mu.RUnlock()
	return s
}

// get a new item from the size pool, return ErrNoEnoughItem when the size pool don't have any item
func (p *sizePool) Get() (interface{}, error) {
	p.mu.RLock()
	if p.pool.Len() == 0 {
		return nil, ErrNoEnoughItem
	}
	p.mu.RUnlock()

	p.mu.Lock()
	item := p.pool.Front()
	if item == nil {
		return nil, ErrNoEnoughItem
	}
	p.pool.Remove(item)
	p.mu.Unlock()

	return item.Value, nil
}

// try to get a new item from the size pool every interval time,be blocked before get the item
func (p *sizePool) BGet(interval time.Duration) (interface{}, error) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		p.mu.Lock()
		item := p.pool.Front()
		if item != nil {
			p.pool.Remove(item)
			p.mu.Unlock()
			return item.Value, nil
		}
	}
}

// put the item back to the size pool, before put the item back, will run reset to clean the item
func (p *sizePool) Put(i interface{}) {
	p.reset(i)
	p.mu.Lock()
	p.pool.PushBack(i)
	p.mu.Unlock()
}
