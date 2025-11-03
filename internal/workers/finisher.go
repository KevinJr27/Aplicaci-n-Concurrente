package workers

import (
	"time"

	"bakery/internal/assets"
	"bakery/internal/models"
	"sync"
)

// Finisher: moves cakes to exit and performs fade-out. Signals pipeline WaitGroup when done.
func Finisher(
	in <-chan *models.Cake,
	state *models.GameState,
	pipelineWg *sync.WaitGroup, // this wg was incremented by pipeline start
) {
	defer pipelineWg.Done()

	for cake := range in {
		// Move to exit Y
		MoveAlongBelt(cake, assets.YExit, state)

		// Fade out effect
		for alpha := 1.0; alpha > 0; alpha -= 0.1 {
			cake.Alpha = alpha
			state.UpdatesChan <- cake
			time.Sleep(50 * time.Millisecond)
		}

		// increase completed counter
		state.Mu.Lock()
		state.Completed++
		state.Mu.Unlock()
	}
}
