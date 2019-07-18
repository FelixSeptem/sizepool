// Package sizepool implement a fix size pool of struct for Go
package sizepool

import (
	"errors"
	"sync"
	"time"

	"github.com/FelixSeptem/collections/queue"
)

const (
	DEFAULT_POOLSIZE int64 = 1024
)

var (
	ErrNoEnoughItem = errors.New("no enough item in size pool")
)

// define the size pool entity
type sizePool struct {
	initsize int64
	new      func() interface{}
	reset    func(interface{})
	pool     *queue.Queue
}

// init a new size pool and return it
func NewPool(size int64, new func() interface{}, reset func(interface{})) *sizePool {
	if size <= 0 {
		size = DEFAULT_POOLSIZE
	}
	sp := &sizePool{
		initsize: size,
		new:      new,
		reset:    reset,
		pool:     queue.NewQueue(int(size)),
	}
	sp.constructNewItem()
	return sp
}

// construct new item and return it to channel
func (p *sizePool) constructNewItem() {
	wg := sync.WaitGroup{}
	for i := 0; int64(i) < p.initsize; i++ {
		wg.Add(1)
		go func() {
			p.pool.Push(p.new())
			wg.Done()
		}()
	}
	wg.Wait()
}

// get size pool init size
func (p *sizePool) InitSize() int64 {
	s := p.initsize
	return s
}

// get a new item from the size pool, return ErrNoEnoughItem when the size pool don't have any item
func (p *sizePool) Get() (interface{}, error) {
	if p.pool.Len() == 0 {
		return nil, ErrNoEnoughItem
	}

	item := p.pool.Pop()

	return item, nil
}

// try to get a new item from the size pool every interval time,be blocked before get the item
func (p *sizePool) BGet(interval time.Duration) (interface{}, error) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		item := p.pool.Pop()
		if item != nil {
			return item, nil
		}
	}
}

// put the item back to the size pool, before put the item back, will run reset to clean the item
func (p *sizePool) Put(i interface{}) {
	p.reset(i)
	p.pool.Push(i)
}
