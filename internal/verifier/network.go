package verifier

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// PeerNode represents a remote verifier auditor in the decentralized network.
type PeerNode struct {
	ID        NodeID    `json:"id"`
	Name      string    `json:"name"`      // e.g. "Alberta Energy Regulator", "Hydro-Québec", "IESO"
	Endpoint  string    `json:"endpoint"`  // https://peer-url:port/api/verifier
	PublicKey string    `json:"public_key"`
	Active    bool      `json:"active"`
	LastSeen  time.Time `json:"last_seen"`
}

// NetworkGossipEnvelope wraps a signed attestation for transport across peer nodes.
type NetworkGossipEnvelope struct {
	Attestation *Attestation `json:"attestation"`
	SenderNode  NodeID       `json:"sender_node"`
	SentAt      time.Time    `json:"sent_at"`
	Hops        int          `json:"hops"`
	Signature   string       `json:"signature"`
}

// DistributedTransport manages cross-provincial verifier gossip and consensus.
type DistributedTransport struct {
	localNode *Node
	peers     map[NodeID]*PeerNode
	client    *http.Client
	mu        sync.RWMutex
}

// NewDistributedTransport creates a network transport for a local verifier node.
func NewDistributedTransport(localNode *Node) *DistributedTransport {
	return &DistributedTransport{
		localNode: localNode,
		peers:     make(map[NodeID]*PeerNode),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// RegisterPeer adds an auditor peer node to the gossip network.
func (dt *DistributedTransport) RegisterPeer(peer *PeerNode) {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	dt.peers[peer.ID] = peer
}

// ListPeers returns all registered auditor peers.
func (dt *DistributedTransport) ListPeers() []*PeerNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	out := make([]*PeerNode, 0, len(dt.peers))
	for _, p := range dt.peers {
		out = append(out, p)
	}
	return out
}

// BroadcastAttestation disseminates a signed attestation to all active network peers.
func (dt *DistributedTransport) BroadcastAttestation(ctx context.Context, att *Attestation) int {
	dt.mu.RLock()
	peers := make([]*PeerNode, 0, len(dt.peers))
	for _, p := range dt.peers {
		if p.Active {
			peers = append(peers, p)
		}
	}
	dt.mu.RUnlock()

	env := NetworkGossipEnvelope{
		Attestation: att,
		SenderNode:  dt.localNode.id,
		SentAt:      time.Now().UTC(),
		Hops:        1,
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", att.MilestoneID, att.EvidenceHash, env.SentAt.Unix())))
	env.Signature = hex.EncodeToString(h[:])

	data, err := json.Marshal(env)
	if err != nil {
		return 0
	}

	successCount := 0
	var wg sync.WaitGroup
	var countMu sync.Mutex

	for _, peer := range peers {
		if peer.ID == dt.localNode.id {
			continue
		}
		wg.Add(1)
		go func(p *PeerNode) {
			defer wg.Done()
			req, err := http.NewRequestWithContext(ctx, "POST", p.Endpoint+"/attestation", bytes.NewReader(data))
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Verifier-Node", string(dt.localNode.id))

			resp, err := dt.client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
				countMu.Lock()
				successCount++
				countMu.Unlock()
			}
		}(peer)
	}

	wg.Wait()
	return successCount
}

// HandleIncomingGossip processes an attestation received over the network from a peer.
func (dt *DistributedTransport) HandleIncomingGossip(env *NetworkGossipEnvelope) error {
	if env == nil || env.Attestation == nil {
		return fmt.Errorf("invalid nil gossip envelope")
	}

	dt.mu.Lock()
	if peer, ok := dt.peers[env.SenderNode]; ok {
		peer.LastSeen = time.Now().UTC()
	}
	dt.mu.Unlock()

	// Ingest attestation into local verifier node
	dt.localNode.mu.Lock()
	defer dt.localNode.mu.Unlock()

	for _, existing := range dt.localNode.attstore {
		if existing.NodeID == env.Attestation.NodeID && existing.MilestoneID == env.Attestation.MilestoneID {
			return nil // idempotent duplicate
		}
	}

	dt.localNode.attstore = append(dt.localNode.attstore, env.Attestation)
	if dt.localNode.store != nil {
		dt.localNode.store.Save(env.Attestation)
	}
	return nil
}

// ServeHTTP exposes the HTTP endpoint for peer gossiping.
func (dt *DistributedTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var env NetworkGossipEnvelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "malformed envelope: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dt.HandleIncomingGossip(&env); err != nil {
		http.Error(w, "gossip rejected: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}
