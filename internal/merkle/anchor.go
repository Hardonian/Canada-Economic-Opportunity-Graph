package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// ProofStep represents one node on the Merkle audit path.
type ProofStep struct {
	Hash   string `json:"hash"`
	IsLeft bool   `json:"is_left"`
}

// InclusionProof provides cryptographic proof that a leaf exists in the root.
type InclusionProof struct {
	LeafHash   string      `json:"leaf_hash"`
	LeafIndex  int         `json:"leaf_index"`
	AuditPath  []ProofStep `json:"audit_path"`
	RootHash   string      `json:"root_hash"`
	VerifiedAt time.Time   `json:"verified_at"`
}

// StateAnchor manages hourly cryptographic state anchors and inclusion proofs.
type StateAnchor struct{}

// NewStateAnchor creates a state anchor generator.
func NewStateAnchor() *StateAnchor {
	return &StateAnchor{}
}

// BuildTreeFromLeaves constructs an explicit binary Merkle tree and returns all levels.
func (sa *StateAnchor) BuildTreeFromLeaves(leafHashes []string) ([][]string, string) {
	if len(leafHashes) == 0 {
		empty := sha256.Sum256([]byte("empty-tree"))
		h := hex.EncodeToString(empty[:])
		return [][]string{{h}}, h
	}

	var levels [][]string
	curr := make([]string, len(leafHashes))
	copy(curr, leafHashes)
	levels = append(levels, curr)

	for len(curr) > 1 {
		var next []string
		for i := 0; i < len(curr); i += 2 {
			if i+1 < len(curr) {
				combined := sha256.Sum256([]byte(curr[i] + curr[i+1]))
				next = append(next, hex.EncodeToString(combined[:]))
			} else {
				// Odd leaf: duplicate or promote
				combined := sha256.Sum256([]byte(curr[i] + curr[i]))
				next = append(next, hex.EncodeToString(combined[:]))
			}
		}
		levels = append(levels, next)
		curr = next
	}

	return levels, curr[0]
}

// GenerateProof produces an inclusion proof for a leaf at a given index.
func (sa *StateAnchor) GenerateProof(leafIndex int, leafHashes []string) (*InclusionProof, error) {
	if leafIndex < 0 || leafIndex >= len(leafHashes) {
		return nil, fmt.Errorf("invalid leaf index: %d", leafIndex)
	}

	levels, root := sa.BuildTreeFromLeaves(leafHashes)
	var steps []ProofStep
	idx := leafIndex

	for level := 0; level < len(levels)-1; level++ {
		currLevel := levels[level]
		var siblingHash string
		isLeft := false

		if idx%2 == 0 {
			// Sibling is on the right
			if idx+1 < len(currLevel) {
				siblingHash = currLevel[idx+1]
			} else {
				siblingHash = currLevel[idx] // duplicated odd
			}
			isLeft = false
		} else {
			// Sibling is on the left
			siblingHash = currLevel[idx-1]
			isLeft = true
		}

		steps = append(steps, ProofStep{
			Hash:   siblingHash,
			IsLeft: isLeft,
		})
		idx /= 2
	}

	return &InclusionProof{
		LeafHash:   leafHashes[leafIndex],
		LeafIndex:  leafIndex,
		AuditPath:  steps,
		RootHash:   root,
		VerifiedAt: time.Now().UTC(),
	}, nil
}

// VerifyInclusionProof validates the cryptographic hash path from leaf to root.
func VerifyInclusionProof(proof *InclusionProof) bool {
	if proof == nil || proof.RootHash == "" {
		return false
	}

	curr := proof.LeafHash
	for _, step := range proof.AuditPath {
		var combined [32]byte
		if step.IsLeft {
			combined = sha256.Sum256([]byte(step.Hash + curr))
		} else {
			combined = sha256.Sum256([]byte(curr + step.Hash))
		}
		curr = hex.EncodeToString(combined[:])
	}

	return curr == proof.RootHash
}
