package lakehouse

import (
	"fmt"
	"sync"
	"time"
)

// QualityRule defines an individual data quality SLA invariant.
type QualityRule struct {
	Field    string
	RuleType string // "NOT_NULL", "POSITIVE_NUMBER", "GEO_BOUNDS", "STRING_LENGTH"
}

// QuarantineEngine inspects records before ingestion and quarantines malformed data.
type QuarantineEngine struct {
	mu         sync.RWMutex
	records    []*QuarantineRecord
	rules      map[string][]QualityRule // table -> rules
}

// NewQuarantineEngine initializes the data quality firewall.
func NewQuarantineEngine() *QuarantineEngine {
	qe := &QuarantineEngine{
		records: make([]*QuarantineRecord, 0),
		rules:   make(map[string][]QualityRule),
	}
	qe.registerDefaultRules()
	return qe
}

func (qe *QuarantineEngine) registerDefaultRules() {
	qe.rules["projects"] = []QualityRule{
		{Field: "id", RuleType: "NOT_NULL"},
		{Field: "name", RuleType: "STRING_LENGTH"},
		{Field: "capex_cad", RuleType: "POSITIVE_NUMBER"},
		{Field: "latitude", RuleType: "GEO_BOUNDS"},
	}
}

// Validate inspects a raw row. If violations are found, it stores the record in quarantine.
func (qe *QuarantineEngine) Validate(table string, raw map[string]interface{}) (bool, *QuarantineRecord) {
	qe.mu.Lock()
	defer qe.mu.Unlock()

	rules := qe.rules[table]
	var violations []string

	for _, r := range rules {
		val, exists := raw[r.Field]
		switch r.RuleType {
		case "NOT_NULL":
			if !exists || val == nil || fmt.Sprintf("%v", val) == "" {
				violations = append(violations, fmt.Sprintf("Field %s must not be null", r.Field))
			}
		case "POSITIVE_NUMBER":
			if exists && val != nil {
				if fval, ok := toFloat64(val); ok {
					if fval < 0 {
						violations = append(violations, fmt.Sprintf("Field %s must be non-negative (observed: %v)", r.Field, val))
					}
				}
			}
		case "GEO_BOUNDS":
			if exists && val != nil {
				if fval, ok := toFloat64(val); ok {
					if fval < -90.0 || fval > 90.0 {
						violations = append(violations, fmt.Sprintf("Field %s must be valid latitude [-90, 90] (observed: %v)", r.Field, val))
					}
				}
			}
		case "STRING_LENGTH":
			if exists && val != nil {
				s := fmt.Sprintf("%v", val)
				if len(s) < 2 {
					violations = append(violations, fmt.Sprintf("Field %s is too short (min length 2)", r.Field))
				}
			}
		}
	}

	if len(violations) > 0 {
		qr := &QuarantineRecord{
			RecordID:    fmt.Sprintf("QR-%d", time.Now().UnixNano()),
			SourceTable: table,
			RawData:     raw,
			Violations:  violations,
			Severity:    "CRITICAL",
			DetectedAt:  time.Now().UTC(),
		}
		qe.records = append(qe.records, qr)
		return false, qr
	}

	return true, nil
}

// GetQuarantinedRecords returns all isolated records.
func (qe *QuarantineEngine) GetQuarantinedRecords() []*QuarantineRecord {
	qe.mu.RLock()
	defer qe.mu.RUnlock()

	result := make([]*QuarantineRecord, len(qe.records))
	copy(result, qe.records)
	return result
}
