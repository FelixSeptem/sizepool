package sizepool

import (
	"sync"
	"time"
)

// define the size pool entity
type sizePoolChan struct {
	initsize int64
	new      func() interface{}
	reset    func(interface{})
	pool     chan interface{}
}

// init a new size pool and return it
func NewChanPool(size int64, new func() interface{}, reset func(interface{})) *sizePoolChan {
	if size <= 0 {
		size = DEFAULT_POOLSIZE
	}
	sp := &sizePoolChan{
		initsize: size,
		new:      new,
		reset:    reset,
		pool:     make(chan interface{}, size),
	}
	sp.constructNewItem()
	return sp
}

// construct new item and return it to channel
func (p *sizePoolChan) constructNewItem() {
	wg := sync.WaitGroup{}
	for i := 0; int64(i) < p.initsize; i++ {
		wg.Add(1)
		go func() {
			p.pool <- p.new()
			wg.Done()
		}()
	}
	wg.Wait()
}

// get size pool init size
func (p *sizePoolChan) InitSize() int64 {
	return p.initsize
}

// try to get a new item from the size pool be blocked before get the item
func (p *sizePoolChan) BGet(timeout time.Duration) (interface{}, error) {
	select {
	case i := <-p.pool:
		return i, nil
	case <-time.Tick(timeout):
		return nil, ErrNoEnoughItem
	}
}

// put the item back to the size pool, before put the item back, will run reset to clean the item
func (p *sizePoolChan) Put(i interface{}) {
	p.reset(i)
	go func() {
		p.pool <- i
	}()
}

func (p *sizePoolChan) Close() {
	close(p.pool)
}
