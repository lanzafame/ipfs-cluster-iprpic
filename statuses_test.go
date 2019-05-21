package iprpic

import (
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
)

var (
	SelfPeerID, _ = peer.IDB58Decode("QmPhGrTF7Lx5pnVztbvo3TFuqJwPcngmcawbXBu1EVDuxz")
	PeerID1, _    = peer.IDB58Decode("QmXZrtE5jQwXNqCJMfHUTQkvhQ4ZAnqMnmzFMJfLewuabc")
	PeerID2, _    = peer.IDB58Decode("QmUZ13osndQ5uL4tPWHXe3iBgBgq9gfewcBMSCAuMBsDJ6")
	PeerID3, _    = peer.IDB58Decode("QmPGDFvBkgWhvzEK9qaTWrWurSwqXNmhnK3hgELPdZZNPa")
	PeerID4, _    = peer.IDB58Decode("QmZ8naDy5mEz4GLuQwjWt9MPYqHTBbsm8tQBrNSjiq6zBc")
	PeerID5, _    = peer.IDB58Decode("QmZVAo3wd8s5eTTy2kPYs34J9PvfxpKPuYsePPYGjgRRjg")
	PeerID6, _    = peer.IDB58Decode("QmR8Vu6kZk7JvAN2rWVWgiduHatgBq2bb15Yyq8RRhYSbx")
	PeerID7, _    = peer.IDB58Decode("QmVHrDHzAfoqcwFMk6m5zs437NEmVYYZEK8ZwghwXvacCW")
	PeerID8, _    = peer.IDB58Decode("QmZRC5qgMQMtFMGJc3ShWAMHk3sNbuhLXnnqd1bpGHbF18")
	PeerID9, _    = peer.IDB58Decode("QmTwBrGo3ak6T6MXfMtnjLBE6ZJe5HsTxtnN6MRaTWSCUk")
	PeerID10, _   = peer.IDB58Decode("QmNvM9taXgxmW95fcN3Ei9Gaw2pQFRjsMM73UGqaRRq73b")
	PeerID11, _   = peer.IDB58Decode("Qmf5J2Qz2dcqM7oyq9aqyNnY6J2im1d2PrTiei8kVAmaC7")
	PeerID12, _   = peer.IDB58Decode("QmVZsvCH9rCg4FPn39uJRqe2e3cfhHrULCGR4pdCrhkxZb")
	PeerID13, _   = peer.IDB58Decode("QmeCxnEjdqZFLy8m3bpaajhJ4fdEe465d88ZhnpgsAeenY")
	PeerID14, _   = peer.IDB58Decode("QmcYuXg35vBajSBwRpGrXkLKmy325cxWEE4GdGDtqCwaja")
	PeerID15, _   = peer.IDB58Decode("QmNdegGvTPPeDngrc6cEq5x1R7mZZqRuNPuJeszQYtKzng")
	PeerID16, _   = peer.IDB58Decode("QmaCwMULXuG8jgyguLAvt9qchBjfY55KsbYrw4CLmknmTF")
	PeerID17, _   = peer.IDB58Decode("QmaVdDgXKmb1VTHQSAaf5iJAQ1VGzQmcRWbvBYPUJnitgk")

	Peers = []peer.ID{
		PeerID1, PeerID2, PeerID3, PeerID4, PeerID5, PeerID6, PeerID7, PeerID8, PeerID9,
		PeerID10, PeerID11, PeerID12, PeerID13, PeerID14, PeerID15, PeerID16, PeerID17,
	}
)

func TestStatuses_AddNewPeer(t *testing.T) {
	statuses := NewStatuses(SelfPeerID)

	for i := 0; i < len(Peers)-1; i++ {
		if ok := statuses.AddNewPeer(Peers[i]); !ok {
			t.Log("should be ok adding 16 peers")
			t.Fail()
		}
	}

	if ok := statuses.AddNewPeer(Peers[len(Peers)-1]); ok {
		t.Log("should not be ok adding 17th peer")
		t.Fail()
	}
}
