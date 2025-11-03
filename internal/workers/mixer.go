package workers

import (
	"bakery/internal/assets"
	"bakery/internal/models"
	"math/rand"
	"sync"
	"time"
)

// MixerWorker: single mixer that processes cakes from 'in' to 'out'
func MixerWorker(
	workerID int,
	in <-chan *models.Cake,
	out chan<- *models.Cake,
	state *models.GameState,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for cake := range in {
		// Move to mixing Y on the belt
		MoveAlongBelt(cake, assets.YMixing, state)

		// Enter station (move left to station)
		cake.TargetX = assets.XStation
		state.UpdatesChan <- cake
		time.Sleep(time.Duration(rand.Intn(400)+200) * time.Millisecond)

		cake.Status = "Mixed"

		// Return to belt
		cake.TargetX = assets.XBelt
		state.UpdatesChan <- cake
		time.Sleep(200 * time.Millisecond)

		out <- cake
	}
}

// MixCoordinator: pool of mixers. Closes out when all mixers finish.
func MixCoordinator(
	in <-chan *models.Cake,
	out chan<- *models.Cake,
	numWorkers int,
	state *models.GameState,
) {
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 1; i <= numWorkers; i++ {
		go MixerWorker(i, in, out, state, &wg)
	}

	// Wait for all mixer workers to finish, then close the next stage
	wg.Wait()
	close(out)
}
