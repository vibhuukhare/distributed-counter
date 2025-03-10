package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/vibhuukhare/distributed-counter/discovery"
)

type CounterHandler struct {
	mu          sync.Mutex
	counter     int
	PeerHandler *PeerHandler
	PeerManager *discovery.PeerManager
}

func NewCounterHandler(pm *discovery.PeerManager) *CounterHandler {
	return &CounterHandler{
		counter:     0,
		PeerHandler: NewPeerHandler(pm),
		PeerManager: pm,
	}
}

type IncrementRequest struct {
	Source string `json:"source"`
}

func (ch *CounterHandler) Increment(w http.ResponseWriter, r *http.Request) {
	ch.mu.Lock()
	ch.counter++
	currentCount := ch.counter
	ch.mu.Unlock()

	data, _ := json.Marshal(map[string]int{"increment": 1})

	// Detect if request is from a peer, without this when incrementing the counter, an infinitwe loop was triggererd
	// as one node was calling another and this was happening continously, so written a check for that
	if r.Header.Get("X-From-Peer") == "true" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"count": currentCount})
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// Broadcast increment to peers
	for _, peer := range ch.PeerManager.GetPeers() {
		go func(p string) {
			log.Printf("Propagating increment to peer: %s", p)

			req, err := http.NewRequest("POST", "http://"+p+"/api/increment", bytes.NewBuffer(data))
			if err != nil {
				log.Println("Failed to create request for peer:", p)
				return
			}

			// Add header to prevent infinite loop
			req.Header.Set("X-From-Peer", "true")
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Println("Failed to propagate increment to:", p, err)
				return
			}
			defer resp.Body.Close()
			log.Printf("Increment successfully propagated to: %s, status: %d", p, resp.StatusCode)
		}(peer)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": currentCount})
}

func (ch *CounterHandler) GetCount(w http.ResponseWriter, r *http.Request) {
	ch.mu.Lock()
	count := ch.counter
	ch.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
