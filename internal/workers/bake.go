package workers

import (
	"bakery/internal/assets"
	"bakery/internal/models"
	"math/rand"
	"time"
)


func Bake(
	in <-chan *models.Cake,
	out chan<- *models.Cake,
	oven *models.Oven,
	state *models.GameState,
) {
	for cake := range in {
		// Move to oven Y
		MoveAlongBelt(cake, assets.YOven, state)

		// Wait for oven slot
		oven.Use(cake)

		// Enter oven station (move left)
		cake.TargetX = assets.XStation
		state.UpdatesChan <- cake

		// Simulate baking time (longer)
		time.Sleep(time.Duration(rand.Intn(2500)+4000) * time.Millisecond)

		cake.Status = "Baked"

		// Return to belt
		cake.TargetX = assets.XBelt
		state.UpdatesChan <- cake
		time.Sleep(200 * time.Millisecond)

		// Release oven slot
		oven.Release(cake)

		out <- cake
	}
	// when input closed, close output
	close(out)
}
