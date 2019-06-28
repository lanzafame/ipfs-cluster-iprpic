package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/lanzafame/bobblehat/sense/screen"
	iprpic "github.com/lanzafame/ipfs-cluster-iprpic"
	peer "github.com/libp2p/go-libp2p-peer"
)

func main() {
	var p = flag.String("p", "", "peer id of current cluster peer")
	flag.Parse()

	if *p == "" {
		fmt.Println("please provide peer id of the cluster peer running on this machine")
		os.Exit(-1)
	}

	log.Print("loading...")
	ctx, cancel := context.WithCancel(context.Background())
	go iprpic.SignalHandler(ctx, cancel)

	pid, err := peer.IDB58Decode(*p)
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
