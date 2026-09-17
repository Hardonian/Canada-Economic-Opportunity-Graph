package verifier

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDistributedTransport_GossipAndQuorum(t *testing.T) {
	// Node 1: Alberta Energy Regulator
	node1 := NewNode("node-aer-01")
	transport1 := NewDistributedTransport(node1)

	// Node 2: Hydro-Québec
	node2 := NewNode("node-hq-02")
	transport2 := NewDistributedTransport(node2)

	// Create test HTTP server for Node 2
	server2 := httptest.NewServer(transport2)
	defer server2.Close()

	// Register Node 2 as peer on Node 1
	transport1.RegisterPeer(&PeerNode{
		ID:       "node-hq-02",
		Name:     "Hydro-Québec Verification Node",
		Endpoint: server2.URL,
		Active:   true,
		LastSeen: time.Now().UTC(),
	})

	// Node 1 observes milestone and signs attestation
	att := &Attestation{
		NodeID:       "node-aer-01",
		MilestoneID:  "ms-hydro-01",
		ProjectID:    "proj-la-romaine",
		EvidenceHash: "a1b2c3d4e5f67890",
		ObservedAt:   time.Now().UTC(),
		Signature:    "sig-aer-01",
	}

	// Node 1 broadcasts to network
	ctx := context.Background()
	delivered := transport1.BroadcastAttestation(ctx, att)
	if delivered != 1 {
		t.Fatalf("expected broadcast delivery to 1 peer, got %d", delivered)
	}

	// Verify Node 2 ingested the attestation via gossip
	node2.mu.RLock()
	found := false
	for _, a := range node2.attstore {
		if a.MilestoneID == "ms-hydro-01" && a.NodeID == "node-aer-01" {
			found = true
			break
		}
	}
	node2.mu.RUnlock()

	if !found {
		t.Fatal("expected Node 2 to ingest gossiped attestation from Node 1")
	}

	// Verify peer list on Node 1
	peers := transport1.ListPeers()
	if len(peers) != 1 || peers[0].ID != "node-hq-02" {
		t.Fatalf("unexpected peers list: %v", peers)
	}
}
