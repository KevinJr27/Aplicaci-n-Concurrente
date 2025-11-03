package workers

import (
	"math/rand"
	"time"

	"bakery/internal/assets"
	"bakery/internal/models"
)

// Packager: receives decorated cakes, packages them.
func Packager(
	in <-chan *models.Cake,
	out chan<- *models.Cake,
	state *models.GameState,
) {
	for cake := range in {
		// Move to packaging Y
		MoveAlongBelt(cake, assets.YPackaging, state)

		// Enter station
		cake.TargetX = assets.XStation
		state.UpdatesChan <- cake

		// Packaging time
		time.Sleep(time.Duration(rand.Intn(600)+400) * time.Millisecond)

		cake.Status = "Ready"

		// Return to belt
		cake.TargetX = assets.XBelt
		state.UpdatesChan <- cake
		time.Sleep(200 * time.Millisecond)

		out <- cake
	}
	close(out)
}
