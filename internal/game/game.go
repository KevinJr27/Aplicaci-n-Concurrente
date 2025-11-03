package game

import (
	"fmt"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"bakery/internal/assets"
	"bakery/internal/models"
	"bakery/internal/workers"
)

const (
	BeltSpeed       = 1.5
	MoveSpeed       = 2.5
	ButtonX, ButtonY = 20, 20
	ButtonW, ButtonH = 160, 50
)

type Game struct {
	State       *models.GameState
	Assets      *assets.ImageAssets
	OrdersChan  chan int
	WgPipeline  sync.WaitGroup
	Initialized bool
}

func (g *Game) Update() error {
	if !g.Initialized {
		g.StartPipeline()
		g.Initialized = true
	}

	// Animate conveyor belt
	g.State.Mu.Lock()
	g.State.BeltOffset += BeltSpeed

	if g.State.BeltOffset > 20 {
		g.State.BeltOffset = 0
	}

	g.State.Mu.Unlock()

	// Handle button click 
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= ButtonX && x <= ButtonX+ButtonW && y >= ButtonY && y <= ButtonY+ButtonH {
			g.MakeOrder()
		}
	}

	// Apply updates (non-blocking) 
	select {
	case updatedCake := <-g.State.UpdatesChan:
		g.State.Mu.Lock()
		found := false
		for i, c := range g.State.Cakes {
			if c.ID == updatedCake.ID {
				g.State.Cakes[i] = updatedCake
				found = true
				break
			}
		}
		if !found {
			g.State.Cakes = append(g.State.Cakes, updatedCake)
		}
		g.State.Mu.Unlock()
	default:
	}

	// Animate cake movement 
	g.State.Mu.Lock()
	for i := len(g.State.Cakes) - 1; i >= 0; i-- {
		c := g.State.Cakes[i]

		dx := c.TargetX - c.X
		dy := c.TargetY - c.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist > MoveSpeed {
			c.X += (dx / dist) * MoveSpeed
			c.Y += (dy / dist) * MoveSpeed
		} else {
			c.X = c.TargetX
			c.Y = c.TargetY
		}

		if c.Alpha <= 0 {
			g.State.Cakes = append(g.State.Cakes[:i], g.State.Cakes[i+1:]...)
		}
	}
	g.State.Mu.Unlock()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 120, 255, 1}) 

	g.DrawBelt(screen)
	g.DrawStations(screen)

	g.State.Mu.Lock()
	for _, c := range g.State.Cakes {
		g.DrawCake(screen, c)
	}
	g.State.Mu.Unlock()

	g.DrawOrderButton(screen)

	// Stats
	g.State.Mu.Lock()
	stats := fmt.Sprintf("Orders: %d | In process: %d | Completed: %d",
		g.State.TotalOrders,
		len(g.State.Cakes),
		g.State.Completed)
	g.State.Mu.Unlock()

	ebitenutil.DebugPrintAt(screen, stats, 20, 700)
	ebitenutil.DebugPrintAt(screen, "Click 'Make Order' to add cakes", 20, 720)
}

// HELPER METHODS
func (g *Game) DrawBelt(screen *ebiten.Image) {
	xBelt := float64(assets.XBelt - assets.BeltWidth/2)
	yBelt := float64(assets.YReception - 50)
	beltHeight := float64(assets.YExit - assets.YReception + 100)

	ebitenutil.DrawRect(screen, xBelt, yBelt, float64(assets.BeltWidth), beltHeight,
		color.RGBA{90, 90, 90, 255})

	g.State.Mu.Lock()
	offset := g.State.BeltOffset
	g.State.Mu.Unlock()

	for y := yBelt - 20; y < yBelt+beltHeight; y += 20 {
		lineY := y + offset
		if lineY > yBelt+beltHeight {
			lineY -= beltHeight
		}
		ebitenutil.DrawRect(screen, xBelt, lineY, float64(assets.BeltWidth), 3,
			color.RGBA{70, 70, 70, 255})
	}

	// Borders
	ebitenutil.DrawRect(screen, xBelt-2, yBelt, 2, beltHeight, color.RGBA{60, 60, 60, 255})
	ebitenutil.DrawRect(screen, xBelt+float64(assets.BeltWidth), yBelt, 2, beltHeight, color.RGBA{60, 60, 60, 255})
}

func (g *Game) DrawStations(screen *ebiten.Image) {
	stations := []struct {
		y     float64
		label string
		img   *ebiten.Image
	}{
		{assets.YReception, "RECEPTION", g.Assets.IconReception},
		{assets.YMixing, "MIXING", g.Assets.IconMixer},
		{assets.YOven, "OVEN (3)", g.Assets.IconOven},
		{assets.YDecorate, "DECORATION", g.Assets.IconDecorator},
		{assets.YPackaging, "PACKAGING", g.Assets.IconPackager},
	}

	for _, st := range stations {
		ebitenutil.DrawRect(screen, assets.XStation-80, st.y-30, 160, 60,
			color.RGBA{255, 255, 255, 220})

		if st.img != nil {
			op := &ebiten.DrawImageOptions{}
			b := st.img.Bounds()
			scaleX := assets.IconSize / float64(b.Dx())
			scaleY := assets.IconSize / float64(b.Dy())
			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(assets.XStation-70, st.y-20)
			screen.DrawImage(st.img, op)
		}

		ebitenutil.DebugPrintAt(screen, st.label, int(assets.XStation-20), int(st.y)-10)
	}
}

func (g *Game) DrawCake(screen *ebiten.Image, c *models.Cake) {
	var img *ebiten.Image
	switch c.Status {
	case "Pending":
		img = g.Assets.CakePending
	case "Mixed":
		img = g.Assets.CakeMixed
	case "Baked":
		img = g.Assets.CakeBaked
	case "Decorated":
		img = g.Assets.CakeDecorated
	case "Ready":
		img = g.Assets.CakeReady
	default:
		img = g.Assets.CakePending
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		b := img.Bounds()
		scale := assets.CakeSize / float64(b.Dx())
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(-assets.CakeSize/2, -assets.CakeSize/2)
		op.GeoM.Translate(c.X, c.Y)
		op.ColorScale.ScaleAlpha(float32(c.Alpha))
		screen.DrawImage(img, op)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("#%d", c.ID), int(c.X)-8, int(c.Y)-5)
}

func (g *Game) DrawOrderButton(screen *ebiten.Image) {
	g.State.Mu.Lock()
	busy := g.State.ReceptionBusy
	g.State.Mu.Unlock()

	var img *ebiten.Image
	if busy && g.Assets.ButtonOrderBusy != nil {
		img = g.Assets.ButtonOrderBusy
	} else if g.Assets.ButtonOrder != nil {
		img = g.Assets.ButtonOrder
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		b := img.Bounds()
		op.GeoM.Scale(160.0/float64(b.Dx()), 50.0/float64(b.Dy()))
		op.GeoM.Translate(20, 20)
		screen.DrawImage(img, op)
	} else {
		btnColor := color.RGBA{100, 200, 100, 255}
		if busy {
			btnColor = color.RGBA{150, 150, 150, 255}
		}
		ebitenutil.DrawRect(screen, 20, 20, 160, 50, btnColor)
		ebitenutil.DebugPrintAt(screen, "MAKE ORDER", 45, 40)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 750
}

// GAME LOGIC
func (g *Game) StartPipeline() {
	g.OrdersChan = make(chan int, 10)
	workers.StartPipeline(g.OrdersChan, g.State, &g.WgPipeline)
}

func (g *Game) MakeOrder() {
	g.State.Mu.Lock()
	if g.State.ReceptionBusy {
		g.State.Mu.Unlock()
		return
	}
	g.State.TotalOrders++
	orderID := g.State.TotalOrders
	g.State.ReceptionBusy = true
	g.State.Mu.Unlock()

	go func() {
		g.OrdersChan <- orderID
		g.State.Mu.Lock()
		g.State.ReceptionBusy = false
		g.State.Mu.Unlock()
	}()
}
