package main

import (
	"log"
	"math/rand"
	"time"

	"bakery/internal/assets"
	"bakery/internal/models"
	"bakery/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load image assets
	imgAssets := assets.LoadImages()

	// Initialize game state
	state := models.NewGameState()

	// Create the main game
	bakeryGame := &game.Game{
		State:       state,
		Assets:      imgAssets,
		OrdersChan:  make(chan int, 20),
	}

	// Ebiten window setup
	ebiten.SetWindowSize(800, 750)
	ebiten.SetWindowTitle("Concurrent Bakery - Pipeline Simulation")

	// Run game
	if err := ebiten.RunGame(bakeryGame); err != nil {
		log.Fatal(err)
	}
}
