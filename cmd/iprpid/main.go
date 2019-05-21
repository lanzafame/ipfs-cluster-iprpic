package main

import (
	"context"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/lanzafame/bobblehat/sense/screen"
	iprpic "github.com/lanzafame/ipfs-cluster-iprpic"
	peer "github.com/libp2p/go-libp2p-peer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go iprpic.SignalHandler(ctx, cancel)

	pid, err := peer.IDB58Decode("QmPhGrTF7Lx5pnVztbvo3TFuqJwPcngmcawbXBu1EVDuxz")
	if err != nil {
		panic(err)
	}

	cfg := &client.Config{}
	c, err := client.NewDefaultClient(cfg)
	if err != nil {
		panic(err)
	}

	fb := screen.NewFrameBuffer()

	s := iprpic.NewStatuses(fb, pid, c)
	ctr := iprpic.NewCounter(fb, c, iprpic.DefaultSegments)

	go s.DrawPeerStatuses(ctx)
	ctr.DrawPinCount(ctx)
}
