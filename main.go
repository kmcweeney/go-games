package main

import (
	"image"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "game",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	win.SetSmooth(true)
	if err != nil {
		panic(err)
	}
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())

	spritesheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}
	tree := pixel.NewSprite(spritesheet, pixel.R(0, 0, 32, 32))
	win.Clear(colornames.Greenyellow)
	tree.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	angle := 0.0
	last := time.Now()
	size := 0
	reverse := false
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Firebrick)
		mat := pixel.IM
		//mat = mat.Rotated(pixel.ZV, angle)
		//fmt.Println("rotating")

		if size > 50 {
			reverse = true
		}
		if size == 1 {
			reverse = false
		}
		if reverse {
			//angle -= 3 * dt
			//mat = mat.Rotated(pixel.ZV, angle)
			//mat = mat.ScaledXY(pixel.ZV, pixel.V(float64(size)*dt, float64(size)*dt))
			size--
			//fmt.Println(size)
		} else {
			//angle += 3 * dt
			//mat = mat.Rotated(pixel.ZV, angle)

			size++
			//fmt.Println(size)
		}
		//mat = mat.Moved(win.Bounds().Center())
		angle += 3 * dt
		mat = mat.Rotated(pixel.ZV, angle)
		mat = mat.ScaledXY(pixel.ZV, pixel.V(float64(size)*0.03, float64(size)*0.03))
		mat = mat.Moved(win.Bounds().Center())

		sprite.Draw(win, mat)
		tree.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func main() {
	//pixelgl.Run(run)
	//ticTacToe()
	c4()
}
