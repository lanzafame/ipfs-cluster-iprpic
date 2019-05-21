package iprpic

import (
	"context"
	"image"
	"log"
	"sync"
	"time"

	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/lanzafame/bobblehat/sense/screen"
	"github.com/lanzafame/bobblehat/sense/screen/color"
)

type pinEvent int

const (
	unknown pinEvent = iota
	pin
	unpin
)

// DefaultSegments provides a default set of segments.
var DefaultSegments = []*Segment{
	{Bucket: 3125, Rectangle: image.Rect(5, 1, 5, 5), Color: color.RandomPlan9PaletteColor()},
	{Bucket: 625, Rectangle: image.Rect(4, 1, 4, 5), Color: color.RandomPlan9PaletteColor()},
	{Bucket: 125, Rectangle: image.Rect(3, 1, 3, 5), Color: color.RandomPlan9PaletteColor()},
	{Bucket: 25, Rectangle: image.Rect(2, 1, 2, 5), Color: color.RandomPlan9PaletteColor()},
	{Bucket: 5, Rectangle: image.Rect(1, 1, 1, 5), Color: color.RandomPlan9PaletteColor()},
	{Bucket: 1, Rectangle: image.Rect(0, 1, 0, 5), Color: color.RandomPlan9PaletteColor()},
}

// Counter represents the section of the display
// that shows the current pin count.
type Counter struct {
	mu       sync.RWMutex
	pinCount int

	Segments []*Segment

	client client.Client
	fb     *screen.FrameBuffer
}

// NewCounter contructs a Counter.
func NewCounter(fb *screen.FrameBuffer, c client.Client, segments []*Segment) *Counter {
	ctx := context.Background()
	cntr := &Counter{fb: fb, client: c, Segments: segments}

	pins, err := cntr.client.StatusAll(ctx, api.TrackerStatusPinned, false)
	if err != nil {
		log.Print(err)
	}

	cntr.mu.Lock()
	cntr.pinCount = len(pins)
	cntr.mu.Unlock()

	cntr.drawCounter(ctx, fb, pin)

	return cntr
}

// DrawPinCount draws a pixel a color everytime that peer
// gets a new metric. Also, does a flashy animation
// indicating whether a pin or unpin occurred.
func (c *Counter) DrawPinCount(ctx context.Context) {
	pinsCh := make(chan pinEvent, 10)

	go c.pins(ctx, pinsCh)

	for {
		select {
		case pe := <-pinsCh:
			go c.drawEvent(ctx, c.fb, pe)
			go c.drawCounter(ctx, c.fb, pe)
		case <-ctx.Done():
			close(pinsCh)
			return
		}
	}

}

func (c *Counter) logPinEvent(newCount int) pinEvent {
	c.mu.Lock()
	defer c.mu.Unlock()
	if newCount > c.pinCount {
		c.pinCount = newCount
		return pin
	}
	if newCount < c.pinCount {
		c.pinCount = newCount
		return unpin
	}
	return unknown
}

func (c *Counter) pins(ctx context.Context, pinsCh chan pinEvent) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pins, err := c.client.StatusAll(ctx, api.TrackerStatusPinned|api.TrackerStatusPinning|api.TrackerStatusUnpinning|api.TrackerStatusQueued, false)
			if err != nil {
				log.Print(err)
				continue
			}
			event := c.logPinEvent(len(pins))
			if event != unknown {
				pinsCh <- event
			}
		case <-ctx.Done():
			return
		}
	}
}

// Segment represents a column segment of the counter.
type Segment struct {
	mu   sync.RWMutex
	v    int
	oldv int

	Bucket int
	image.Rectangle

	Color color.Pixel565
}

func (c *Counter) segmentPinCount(v int) {
	c.mu.Lock()
	for _, b := range c.Segments {
		b.oldv = b.v
		b.v = v / b.Bucket
		v = v % b.Bucket
	}
	c.mu.Unlock()
}

func (c *Counter) drawColumn(ctx context.Context, fb *screen.FrameBuffer, s *Segment) {
	for y := s.Max.Y; y >= s.Max.Y-s.v; y-- {
		fb.SetPixel(s.Min.X, y, s.Color)
	}
	for y := s.Max.Y - s.v; y >= s.Min.Y; y-- {
		fb.SetPixel(s.Min.X, y, color.Black)
	}
	screen.Draw(fb)
}

func (c *Counter) drawCounter(ctx context.Context, fb *screen.FrameBuffer, event pinEvent) {
	c.segmentPinCount(c.pinCount)
	for i := len(c.Segments) - 1; i >= 0; i-- {
		s := c.Segments[i]
		// flash column when it fills and
		// consolidates into the next column
		// v == 0 => no remainders to modulo operation
		// s.v != s.oldv => previous value wasn't also 0
		if s.v == 0 && s.v < s.oldv {
			switch event {
			case pin:
				drawColumnWave(fb, s.Rectangle, color.White, 30*time.Millisecond, Up)
			case unpin:
				drawColumnWave(fb, s.Rectangle, color.White, 30*time.Millisecond, Down)
			}
		}
		c.drawColumn(ctx, fb, s)
	}
}

func (c *Counter) drawEvent(ctx context.Context, fb *screen.FrameBuffer, event pinEvent) {
	rect := image.Rect(0, 6, 7, 7)

	switch event {
	case pin:
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			fb.SetPixel(x, rect.Min.Y, color.Blue)
			fb.SetPixel(x, rect.Max.Y, color.Blue)
			screen.Draw(fb)
			time.Sleep(20 * time.Millisecond)
			fb.SetPixel(x, rect.Min.Y, color.Black)
			fb.SetPixel(x, rect.Max.Y, color.Black)
			screen.Draw(fb)
		}
	case unpin:
		for x := rect.Max.X; x >= rect.Min.X; x-- {
			fb.SetPixel(x, rect.Min.Y, color.Red)
			fb.SetPixel(x, rect.Max.Y, color.Red)
			screen.Draw(fb)
			time.Sleep(20 * time.Millisecond)
			fb.SetPixel(x, rect.Min.Y, color.Black)
			fb.SetPixel(x, rect.Max.Y, color.Black)
			screen.Draw(fb)
		}
	}
}
