package workers

import (
	"bakery/internal/assets"
	models "bakery/internal/models"
	"math/rand"
	"time"

	"sync"
)

func Receptionist(
	orders <-chan int,
	out chan<- *models.Cake,
	state *models.GameState,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for id := range orders {
		state.Mu.Lock()
		state.ReceptionBusy = true
		state.Mu.Unlock()

		time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)

		cake := &models.Cake{
			ID:      id,
			Status:  "Pending",
			X:       assets.XStation,
			Y:       assets.YReception,
			TargetX: assets.XStation,
			TargetY: assets.YReception,
			Alpha:   1.0,
		}

		state.NewCakesChan <- cake

		state.Mu.Lock()
		state.ReceptionBusy = false
		state.Mu.Unlock()

		out <- cake
	}
}
