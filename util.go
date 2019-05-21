package iprpic

import (
	"image"
	"time"

	"github.com/lanzafame/bobblehat/sense/screen"
	"github.com/lanzafame/bobblehat/sense/screen/color"
)

// Direction is the direction that an animation should go
// across the screen.
type Direction int

// Directions
const (
	Up Direction = iota + 1
	Down
	Left
	Right
)

func drawRect(fb *screen.FrameBuffer, r image.Rectangle, c color.Pixel565) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			fb.SetPixel(x, y, c)
		}
	}
	screen.Draw(fb)
}

func drawColumnWave(fb *screen.FrameBuffer, r image.Rectangle, c color.Pixel565, speed time.Duration, dir Direction) {
	switch dir {
	case Up:
		for y := r.Max.Y; y > r.Min.Y; y-- {
			fb.SetPixel(r.Min.X, y, c)
			screen.Draw(fb)
			time.Sleep(speed)
			if y < r.Max.Y {
				fb.SetPixel(r.Min.X, y+1, color.Black)
				screen.Draw(fb)
			}
		}
		fb.SetPixel(r.Min.X, r.Min.Y+1, color.Black)
		screen.Draw(fb)
	case Down:
		for y := r.Min.Y + 1; y <= r.Max.Y; y++ {
			fb.SetPixel(r.Min.X, y, c)
			screen.Draw(fb)
			time.Sleep(speed)
			if y > r.Min.Y {
				fb.SetPixel(r.Min.X, y-1, color.Black)
				screen.Draw(fb)
			}
		}
		fb.SetPixel(r.Min.X, r.Max.Y, color.Black)
		screen.Draw(fb)
	}
}
