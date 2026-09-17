package ontology

import (
	"time"
)

// TemporalEvaluator reconstructs point-in-time entity state using bi-temporal coordinates.
type TemporalEvaluator struct {
	writeEngine *WriteBackEngine
}

// NewTemporalEvaluator creates a temporal query evaluator.
func NewTemporalEvaluator(writeEngine *WriteBackEngine) *TemporalEvaluator {
	return &TemporalEvaluator{
		writeEngine: writeEngine,
	}
}

// GetNodeAsOf reconstructs the properties of a node as of a specific real-world ValidTime and SystemTime.
func (te *TemporalEvaluator) GetNodeAsOf(nodeID string, asOfValid time.Time, asOfSystem time.Time) map[string]interface{} {
	if te.writeEngine == nil {
		return make(map[string]interface{})
	}

	actions := te.writeEngine.GetActionHistory(nodeID)
	result := make(map[string]interface{})

	for _, a := range actions {
		if a.Status != ActionStatusApplied {
			continue
		}
		// Filter by SystemTime (when transaction was known to the database)
		if !asOfSystem.IsZero() && a.SystemTime.After(asOfSystem) {
			continue
		}
		// Filter by ValidTime (when event took real-world effect)
		if !asOfValid.IsZero() && a.ValidTime.After(asOfValid) {
			continue
		}

		// Apply forward mutations up to this point
		for k, v := range a.Parameters {
			result[k] = v
		}
	}

	return result
}
