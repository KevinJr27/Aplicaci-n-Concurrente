package workers

import (
	"math/rand"
	"time"

	"bakery/internal/assets"
	"bakery/internal/models"
)

// MoveAlongBelt animates a cake moving to the belt (x center) and a target Y.
// It sends an update to the state's UpdatesChan so the renderer can animate it.
func MoveAlongBelt(c *models.Cake, destY float64, state *models.GameState) {
	c.TargetX = assets.XBelt
	c.TargetY = destY
	state.UpdatesChan <- c
	// small randomized wait to simulate travel time
	time.Sleep(time.Duration(rand.Intn(200)+300) * time.Millisecond)
}
