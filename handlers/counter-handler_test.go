package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vibhuukhare/distributed-counter/discovery"
)

func TestIncrement(t *testing.T) {
	
	selfAddress := "localhost:8080"
	peers := []string{"localhost:8081", "localhost:8082"} // mock peers

	peerManager := discovery.NewPeerManager(selfAddress, peers)

	counterHandler := NewCounterHandler(peerManager)

	req, err := http.NewRequest("POST", "/api/increment", bytes.NewBuffer([]byte(`{"source": "test"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(counterHandler.Increment)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if _, exists := response["count"]; !exists {
		t.Errorf("response body does not contain 'count' field")
	}
}