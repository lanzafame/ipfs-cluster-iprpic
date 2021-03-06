package main

import (
	"context"
	"log"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/lanzafame/bobblehat/sense/screen"
	iprpic "github.com/lanzafame/ipfs-cluster-iprpic"
)

func main() {
	log.Print("loading...")
	ctx, cancel := context.WithCancel(context.Background())
	go iprpic.SignalHandler(ctx, cancel)

	cfg := &client.Config{}
	c, err := client.NewDefaultClient(cfg)
	if err != nil {
		panic(err)
	}

	p, err := c.ID(ctx)
	if err != nil {
		panic(err)
	}

	fb := screen.NewFrameBuffer()
	s := iprpic.NewStatuses(fb, p.ID, c)
	ctr := iprpic.NewCounter(fb, c, iprpic.DefaultSegments)

	go s.DrawPeerStatuses(ctx)
	ctr.DrawPinCount(ctx)
}
