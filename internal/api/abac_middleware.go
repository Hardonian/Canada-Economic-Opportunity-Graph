package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/telemetry"
)

type securitySubjectContextKey struct{}

// ExtractSubject parses security context attributes from request headers.
func ExtractSubject(r *http.Request) security.SecuritySubject {
	clearanceRaw := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-Clearance-Level")))
	clearance := security.ClassUnclassified
	switch security.ClassificationLevel(clearanceRaw) {
	case security.ClassProtectedA:
		clearance = security.ClassProtectedA
	case security.ClassProtectedB:
		clearance = security.ClassProtectedB
	case security.ClassProtectedC:
		clearance = security.ClassProtectedC
	case security.ClassConfidential:
		clearance = security.ClassConfidential
	case security.ClassSecret:
		clearance = security.ClassSecret
	case security.ClassTopSecret:
		clearance = security.ClassTopSecret
	}

	citizenship := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-Citizenship")))
	if citizenship == "" {
		citizenship = "CA"
	}

	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		userID = "anonymous"
	}

	org := strings.TrimSpace(r.Header.Get("X-Organization"))
	if org == "" {
		org = "PUBLIC"
	}

	var compartments []string
	if compHeader := r.Header.Get("X-Compartments"); compHeader != "" {
		for _, c := range strings.Split(compHeader, ",") {
			c = strings.TrimSpace(strings.ToUpper(c))
			if c != "" {
				compartments = append(compartments, c)
			}
		}
	}

	return security.SecuritySubject{
		UserID:       userID,
		Clearance:    clearance,
		Citizenship:  citizenship,
		Organization: org,
		Compartments: compartments,
	}
}

// SubjectFromContext retrieves the evaluated security subject from the request context.
func SubjectFromContext(ctx context.Context) security.SecuritySubject {
	if sub, ok := ctx.Value(securitySubjectContextKey{}).(security.SecuritySubject); ok {
		return sub
	}
	return security.SecuritySubject{
		UserID:       "anonymous",
		Clearance:    security.ClassUnclassified,
		Citizenship:  "CA",
		Organization: "PUBLIC",
	}
}

// AuthorizeRequest evaluates the request against an ABAC security object.
func AuthorizeRequest(r *http.Request, obj security.SecurityObject, action string) security.ABACDecision {
	sub := ExtractSubject(r)
	evaluator := security.NewABACEvaluator()
	decision := evaluator.Authorize(sub, obj, action)
	if !decision.Permitted {
		telemetry.DefaultCollector.IncABACDenial()
	}
	return decision
}

// RedactProjectForSubject returns a clone of the project with fields redacted if the subject holds insufficient clearance.
func RedactProjectForSubject(p *domain.Project, sub security.SecuritySubject) *domain.Project {
	if p == nil {
		return nil
	}
	clone := *p

	// If classified project metadata exists or project requires Protected B+ for exact CapEx
	subRank := security.ClearanceRank[sub.Clearance]
	if subRank < security.ClearanceRank[security.ClassProtectedA] {
		// Public tier: retain standard verified data
	}
	if subRank < security.ClearanceRank[security.ClassSecret] {
		// If sensitivity marker is SECRET or TOP_SECRET, redact internal audit notes or defense identifiers
		if strings.Contains(strings.ToUpper(clone.Sector), "DEFENCE") || strings.Contains(strings.ToUpper(clone.Subsector), "MILITARY") {
			clone.Summary = "[SUMMARY RESTRICTED - REQUIRES CANADIAN SECRET CLEARANCE]"
		}
	}
	return &clone
}
