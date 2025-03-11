package discovery

import (
	"testing"
	"time"
)

func TestAddPeer(t *testing.T) {
	pm := NewPeerManager("localhost:8080", []string{})
	pm.AddPeer("localhost:8081")

	if len(pm.GetPeers()) != 1 {
		t.Errorf("Expected 1 peer, got %d", len(pm.GetPeers()))
	}
}

func TestRemovePeer(t *testing.T) {
	pm := NewPeerManager("localhost:8080", []string{})
	pm.AddPeer("localhost:8081")
	pm.RemovePeer("localhost:8081")

	if len(pm.GetPeers()) != 0 {
		t.Errorf("Expected 0 peers, got %d", len(pm.GetPeers()))
	}
}

func TestHeartbeat(t *testing.T) {
	pm := NewPeerManager("localhost:8080", []string{})
	pm.AddPeer("localhost:8081")

	time.AfterFunc(2*time.Second, func() {
		pm.RemovePeer("localhost:8081")
	})

	time.Sleep(3 * time.Second)

	if len(pm.GetPeers()) != 0 {
		t.Errorf("Expected 0 peers after heartbeat, got %d", len(pm.GetPeers()))
	}
}