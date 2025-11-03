package assets

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageAssets struct {
	CakePending   *ebiten.Image
	CakeMixed     *ebiten.Image
	CakeBaked     *ebiten.Image
	CakeDecorated *ebiten.Image
	CakeReady     *ebiten.Image

	IconReception *ebiten.Image
	IconMixer     *ebiten.Image
	IconOven      *ebiten.Image
	IconDecorator *ebiten.Image
	IconPackager  *ebiten.Image

	Curtain         *ebiten.Image
	ButtonOrder     *ebiten.Image
	ButtonOrderBusy *ebiten.Image
}

func LoadImages() *ImageAssets {
	img := &ImageAssets{}
	img.CakePending = load("assets/cake_pending.png")
	img.CakeMixed = load("assets/cake_mixed.png")
	img.CakeBaked = load("assets/cake_baked.png")
	img.CakeDecorated = load("assets/cake_decorated.png")
	img.CakeReady = load("assets/cake_ready.png")

	img.IconReception = load("assets/icon_reception.png")
	img.IconMixer = load("assets/icon_mixer.png")
	img.IconOven = load("assets/icon_oven.png")
	img.IconDecorator = load("assets/icon_decorator.png")
	img.IconPackager = load("assets/icon_packager.png")

	img.Curtain = load("assets/curtain.png")
	img.ButtonOrder = load("assets/button_order.png")
	img.ButtonOrderBusy = load("assets/button_order_busy.png")
	return img
}

func load(path string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("Could not load %s, using placeholder", path)
		placeholder := ebiten.NewImage(64, 64)
		placeholder.Fill(color.RGBA{200, 200, 200, 255})
		return placeholder
	}
	return image
}
