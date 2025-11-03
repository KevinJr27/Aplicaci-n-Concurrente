package workers

import (
	"bakery/internal/models"
	"sync"
)

func StartPipeline(
	orders <-chan int,
	state *models.GameState,
	pipelineWg *sync.WaitGroup,
) {
	
	pending := make(chan *models.Cake, 10)
	mixed := make(chan *models.Cake, 10)
	baked := make(chan *models.Cake, 10)
	decorated := make(chan *models.Cake, 10)
	ready := make(chan *models.Cake, 10)

	oven := models.NewOven(3)

	pipelineWg.Add(1)

	var receptionWg sync.WaitGroup
	receptionWg.Add(1)

	go StartReceptionist(orders, pending, state, &receptionWg)

	go MixCoordinator(pending, mixed, 5, state)
	go Bake(mixed, baked, oven, state)
	go Decorator(baked, decorated, state)
	go Packager(decorated, ready, state)
	go Finisher(ready, state, pipelineWg)
}

func StartReceptionist(
	orders <-chan int,
	out chan<- *models.Cake,
	state *models.GameState,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	Receptionist(orders, out, state, wg)
}