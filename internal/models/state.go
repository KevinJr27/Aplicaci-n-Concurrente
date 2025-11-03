package models

import "sync"

type GameState struct {
	Mu                 sync.Mutex
	Cakes              []*Cake
	NextID             int
	TotalOrders        int
	Completed          int
	ReceptionBusy      bool
	BeltOffset         float64
	UpdatesChan        chan *Cake
	NewCakesChan       chan *Cake
}

func NewGameState() *GameState {
	return &GameState{
		Cakes:        make([]*Cake, 0),
		UpdatesChan:  make(chan *Cake, 50),
		NewCakesChan: make(chan *Cake, 50),
	}
}
