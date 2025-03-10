package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type PeerManager struct {
	mu       sync.Mutex
	peers    map[string]time.Time
	selfAddr string
}

func NewPeerManager(selfAddr string, initialPeers []string) *PeerManager {
	pm := &PeerManager{
		peers:    make(map[string]time.Time),
		selfAddr: selfAddr,
	}

	// Add initial peers and register with them
	for _, peer := range initialPeers {
		if peer != "" && peer != selfAddr {
			pm.AddPeer(peer)
			pm.registerWithPeer(peer)
		}
	}
	return pm
}

func (pm *PeerManager) AddPeer(peer string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if peer != pm.selfAddr {
		pm.peers[peer] = time.Now()
		log.Println("Added peer:", peer)
	}
}

func (pm *PeerManager) RemovePeer(peer string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.peers, peer)
	log.Println("Removed peer:", peer)
}

func (pm *PeerManager) GetPeers() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var peers []string
	for peer := range pm.peers {
		peers = append(peers, peer)
	}
	return peers
}

// Periodically checks if peers are alive
func (pm *PeerManager) HeartBeatCheck() {
	for {
		time.Sleep(5 * time.Second)

		pm.mu.Lock()
		for peer := range pm.peers {
			go func(p string) {
				resp, err := http.Get(fmt.Sprintf("http://%s/api/health", p))
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Println("Peer failed:", p)
					pm.RemovePeer(p)
				} else {
					pm.mu.Lock()
					pm.peers[p] = time.Now()
					pm.mu.Unlock()
				}
			}(peer)
		}
		pm.mu.Unlock()
	}
}

// Registers the current instance with an existing peer
func (pm *PeerManager) registerWithPeer(peer string) {
	log.Println("Registering with peer:", peer)
	data := map[string]string{"peer": pm.selfAddr}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://"+peer+"/api/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to register with peer:", peer, err)
		return
	}
	defer resp.Body.Close()

	log.Println("Successfully registered with peer:", peer)
}
