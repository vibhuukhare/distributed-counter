package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vibhuukhare/distributed-counter/discovery"
)

type PeerHandler struct {
	PeerManager *discovery.PeerManager
}

func NewPeerHandler(pm *discovery.PeerManager) *PeerHandler {
	return &PeerHandler{PeerManager: pm}
}

func (ph *PeerHandler) RegisterPeer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Peer string `json:"peer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ph.PeerManager.AddPeer(req.Peer)
	log.Println("Registered new peer:", req.Peer)
	w.WriteHeader(http.StatusOK)
}

func (ph *PeerHandler) GetPeers(w http.ResponseWriter, r *http.Request) {
	peers := ph.PeerManager.GetPeers()
	json.NewEncoder(w).Encode(peers)
}
