package workers

import (
	"math/rand"
	"time"

	"bakery/internal/assets"
	"bakery/internal/models"
)

// Decorator: receives baked cakes, decorates them and sends forward.
func Decorator(
	in <-chan *models.Cake,
	out chan<- *models.Cake,
	state *models.GameState,
) {
	for cake := range in {
		// Move to decorate Y
		MoveAlongBelt(cake, assets.YDecorate, state)

		// Small approach into station
		time.Sleep(300 * time.Millisecond)

		cake.TargetX = assets.XStation
		state.UpdatesChan <- cake

		// Decorate time
		time.Sleep(time.Duration(rand.Intn(800)+700) * time.Millisecond)

		cake.Status = "Decorated"

		// Return to belt
		cake.TargetX = assets.XBelt
		state.UpdatesChan <- cake
		time.Sleep(200 * time.Millisecond)

		out <- cake
	}
	// close next stage when done
	close(out)
}
