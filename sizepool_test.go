package sizepool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeConn struct {
	Addr    string
	Timeout time.Duration
}

func newConn() interface{} {
	// to simulate a expensive construct process
	time.Sleep(time.Millisecond * 30)
	return new(fakeConn)
}

func resetConn(i interface{}) {
	c := i.(*fakeConn)
	c.Addr = ""
	c.Timeout = 0
}

func TestNewPool(t *testing.T) {
	p := NewPool(1024, newConn, resetConn)
	assert.Equal(t, int64(1024), p.InitSize())
	assert.Equal(t, 1024, p.pool.Len())
}

func TestSizePool_InitSize(t *testing.T) {
	p := NewPool(1024, newConn, resetConn)
	assert.Equal(t, int64(1024), p.InitSize())
}

func TestSizePool_Get(t *testing.T) {
	p := NewPool(1024, newConn, resetConn)

	for i := 0; i < 1024; i++ {
		item, err := p.Get()
		assert.NoError(t, err)
		assert.Equal(t, new(fakeConn), item.(*fakeConn))
	}

	item, err := p.Get()
	assert.Error(t, ErrNoEnoughItem, err)
	assert.Equal(t, nil, item)
}

func TestSizePool_Put(t *testing.T) {
	p := NewPool(1024, newConn, resetConn)

	p.Put(new(fakeConn))
	assert.Equal(t, 1025, p.pool.Len())
	assert.Equal(t, new(fakeConn), p.pool.Back().Value.(*fakeConn))
}

func TestSizePool_BGet(t *testing.T) {
	p := NewPool(1024, newConn, resetConn)

	i, err := p.BGet(time.Millisecond * 50)
	assert.NoError(t, err)
	assert.Equal(t, new(fakeConn), i.(*fakeConn))
}

func BenchmarkNewPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPool(1024, newConn, resetConn)
	}
}

func BenchmarkSizePool_Get(b *testing.B) {
	p := NewPool(65535, newConn, resetConn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Get()
	}
}

func BenchmarkSizePool_BGet(b *testing.B) {
	p := NewPool(65535, newConn, resetConn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.BGet(time.Nanosecond)
	}
}

func BenchmarkSizePool_Put(b *testing.B) {
	p := NewPool(1, newConn, resetConn)
	fc := new(fakeConn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Put(fc)
	}
}
