package models

import "sync"

type Oven struct {
	Capacity int
	InUse    int
	mu       sync.Mutex
	cond     *sync.Cond
}

func NewOven(capacity int) *Oven {
	o := &Oven{Capacity: capacity}
	o.cond = sync.NewCond(&o.mu)
	return o
}

func (o *Oven) Use(cake *Cake) {
	o.mu.Lock()
	for o.InUse >= o.Capacity {
		o.cond.Wait()
	}
	o.InUse++
	cake.InOven = true
	o.mu.Unlock()
}

func (o *Oven) Release(cake *Cake) {
	o.mu.Lock()
	o.InUse--
	cake.InOven = false
	o.mu.Unlock()
	o.cond.Signal()
}
