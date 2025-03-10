package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/vibhuukhare/distributed-counter/discovery"

	"github.com/vibhuukhare/distributed-counter/handlers"
)

func main() {
	host := "localhost"

	// Read port from command-line argument
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	selfAddr := fmt.Sprintf("%s:%s", host, *port)

	// Get peer list from environment variable
	rawPeers := os.Getenv("PEERS")
	var initialPeers []string
	if rawPeers != "" {
		initialPeers = strings.Split(rawPeers, ",")
	}

	// Initialize PeerManager
	peerManager := discovery.NewPeerManager(selfAddr, initialPeers)

	// Initialize handlers
	counterHandler := handlers.NewCounterHandler(peerManager)

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", counterHandler.PeerHandler.RegisterPeer)
	mux.HandleFunc("/api/peers", counterHandler.PeerHandler.GetPeers)
	mux.HandleFunc("/api/increment", counterHandler.Increment)
	mux.HandleFunc("/api/count", counterHandler.GetCount)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Start heartbeat checks, will run in the background
	go peerManager.HeartBeatCheck()

	// Start server
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	go func() {
		log.Println("Server started on", selfAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
