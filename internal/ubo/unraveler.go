package ubo

import (
	"strings"
)

// SecrecyJurisdictions defines offshore jurisdictions with opaque beneficial ownership registries.
var SecrecyJurisdictions = map[string]bool{
	"KY": true, "CAYMAN ISLANDS": true,
	"VG": true, "BRITISH VIRGIN ISLANDS": true, "BVI": true,
	"CY": true, "CYPRUS": true,
	"PA": true, "PANAMA": true,
	"SC": true, "SEYCHELLES": true,
	"MH": true, "MARSHALL ISLANDS": true,
	"LI": true, "LIECHTENSTEIN": true,
}

// HoldingNode represents a corporate entity in an ownership tree.
type HoldingNode struct {
	EntityID     string  `json:"entity_id"`
	Name         string  `json:"name"`
	Jurisdiction string  `json:"jurisdiction"`
	DirectShare  float64 `json:"direct_share"` // 0.0 - 1.0
	ParentID     string  `json:"parent_id,omitempty"`
}

// UnraveledOwner represents an end-of-chain ultimate beneficial owner.
type UnraveledOwner struct {
	UltimateEntityID     string   `json:"ultimate_entity_id"`
	Name                 string   `json:"name"`
	Jurisdiction         string   `json:"jurisdiction"`
	EffectiveShare       float64  `json:"effective_share"` // Multiplied through chain
	ChainHops            int      `json:"chain_hops"`
	OffshoreSecrecyHops  int      `json:"offshore_secrecy_hops"`
	IntermediaryEntities []string `json:"intermediary_entities"`
	IsOffshoreObfuscated bool     `json:"is_offshore_obfuscated"`
}

// UBOUnraveler resolves multi-tier cross-border ownership trees down to natural persons or SOEs.
type UBOUnraveler struct{}

// NewUBOUnraveler creates an ownership unraveler.
func NewUBOUnraveler() *UBOUnraveler {
	return &UBOUnraveler{}
}

// UnravelChain walks an ownership hierarchy and computes effective shares.
func (u *UBOUnraveler) UnravelChain(targetID string, holdings []HoldingNode) []UnraveledOwner {
	// Build child -> parent mapping
	parentMap := make(map[string]HoldingNode)
	for _, h := range holdings {
		parentMap[h.EntityID] = h
	}

	// Find leaf nodes (entities that are not parents of any other in the tree)
	isParent := make(map[string]bool)
	for _, h := range holdings {
		if h.ParentID != "" {
			isParent[h.ParentID] = true
		}
	}

	var leafNodes []HoldingNode
	for _, h := range holdings {
		if !isParent[h.EntityID] && h.ParentID != "" {
			leafNodes = append(leafNodes, h)
		}
	}

	var results []UnraveledOwner

	for _, leaf := range leafNodes {
		curr := leaf
		effectiveShare := 1.0
		hops := 0
		secrecyHops := 0
		var intermediaries []string

		for {
			effectiveShare *= curr.DirectShare
			hops++
			if SecrecyJurisdictions[strings.ToUpper(curr.Jurisdiction)] {
				secrecyHops++
			}

			if curr.ParentID == "" || curr.ParentID == targetID {
				break
			}
			intermediaries = append(intermediaries, curr.Name)

			parent, ok := parentMap[curr.ParentID]
			if !ok {
				break
			}
			curr = parent
		}

		results = append(results, UnraveledOwner{
			UltimateEntityID:     leaf.EntityID,
			Name:                 leaf.Name,
			Jurisdiction:         leaf.Jurisdiction,
			EffectiveShare:       effectiveShare,
			ChainHops:            hops,
			OffshoreSecrecyHops:  secrecyHops,
			IntermediaryEntities: intermediaries,
			IsOffshoreObfuscated: secrecyHops > 0,
		})
	}

	return results
}
