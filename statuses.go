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
	peer "github.com/libp2p/go-libp2p-peer"
)

// Statuses represents the section of the display
// that shows the status of the other peers in the
// cluster.
type Statuses struct {
	Self peer.ID

	mu    sync.RWMutex
	Peers []*Peer

	client client.Client
	fb     *screen.FrameBuffer
}

// NewStatuses constructs a Statuses type with a Cluster
// client.
func NewStatuses(fb *screen.FrameBuffer, self peer.ID, c client.Client) *Statuses {
	ps := make([]*Peer, 0, 10)
	return &Statuses{fb: fb, Self: self, Peers: ps, client: c}
}

// AddNewPeer adds a newly found peer to the Statuses Peers
// slice and assigns it a new position. Returns false if
// there is no more room in the slice.
func (s *Statuses) AddNewPeer(pid peer.ID) bool {
	log.Print("adding new peer")
	s.mu.RLock()
	if len(s.Peers) == cap(s.Peers) {
		return false
	}
	s.mu.RUnlock()

	p := &Peer{ID: pid}
	// p.Color = color.RandomPlan9PaletteColor()
	p.Color = color.Magenta

	s.mu.RLock()
	numPeers := len(s.Peers)
	maxPeers := 5
	s.mu.RUnlock()
	x := (numPeers / maxPeers) + 6
	y := (numPeers % maxPeers) + 1
	p.Point = image.Point{X: x, Y: y}

	s.mu.Lock()
	s.Peers = append(s.Peers, p)
	s.mu.Unlock()

	return true
}

// DropExpiredPeer removes a peer from the Statuses Peers
// slice.
func (s *Statuses) DropExpiredPeer(pid peer.ID) bool {
	log.Print("dropping peer")
	var i int
	var p *Peer
	s.mu.RLock()
	for i, p = range s.Peers {
		if p.ID == pid {
			s.mu.RUnlock()
			goto DELETE
		}
	}
	s.mu.RUnlock()
	return false

DELETE:
	s.mu.Lock()
	copy(s.Peers[i:], s.Peers[i+1:])
	s.Peers[len(s.Peers)-1] = nil // or the zero value of T
	s.Peers = s.Peers[:len(s.Peers)-1]
	s.mu.Unlock()
	return true
}

// Full returns if the Peers array is full.
func (s *Statuses) Full() bool {
	s.mu.RLock()
	b := len(s.Peers) == cap(s.Peers)
	s.mu.RUnlock()
	return b
}

// Peer returns the peer object from the Peers array.
func (s *Statuses) Peer(pid peer.ID) *Peer {
	s.mu.RLock()
	for _, p := range s.Peers {
		if p == nil {
			continue
		}
		if p.ID == pid {
			s.mu.RUnlock()
			return p
		}
	}
	s.mu.RUnlock()
	return nil
}

// Peer contains the peer id and its position on the
// LED display.
type Peer struct {
	peer.ID
	image.Point
	Color color.Pixel565
}

// DrawPeerStatuses draws a pixel a color everytime that peer
// gets a new metric.
func (s *Statuses) DrawPeerStatuses(ctx context.Context) {
	metricsCh := make(chan *api.Metric, 10)

	go s.metrics(ctx, metricsCh)

	for {
		select {
		case m := <-metricsCh:
			s.mu.RLock()
			p := s.Peer(m.Peer)
			if p != nil {
				p.flicker(ctx, s.fb)
			}
			s.mu.RUnlock()
		case <-ctx.Done():
			close(metricsCh)
			return
		}
	}

}

func (p *Peer) flicker(ctx context.Context, fb *screen.FrameBuffer) {
	fb.SetPixel(p.X, p.Y, p.Color)
	screen.Draw(fb)

	ticker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			fb.SetPixel(p.X, p.Y, color.Black)
			screen.Draw(fb)
			return
		case <-ctx.Done():
			screen.Clear()
			return
		}
	}
}

func (s *Statuses) metrics(ctx context.Context, metricsCh chan *api.Metric) {
	// make call to cluster
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics, err := s.client.Metrics(ctx, "ping")
			if err != nil {
				log.Print(err)
				continue
			}

			for _, m := range metrics {
				if m.Peer != s.Self {
					if !m.Expired() {
						if !s.Full() && s.Peer(m.Peer) == nil {
							s.AddNewPeer(m.Peer)
						}
						if s.Peer(m.Peer) != nil {
							metricsCh <- m
						}
					} else {
						s.DropExpiredPeer(m.Peer)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
