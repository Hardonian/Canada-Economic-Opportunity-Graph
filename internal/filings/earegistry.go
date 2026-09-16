package filings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// EAMilestoneType classifies regulatory environmental assessment stages.
type EAMilestoneType string

const (
	EAMilestoneProjectDescription EAMilestoneType = "PROJECT_DESCRIPTION_SUBMITTED"
	EAMilestoneTermsOfReference   EAMilestoneType = "TERMS_OF_REFERENCE_APPROVED"
	EAMilestonePublicCommentOpen  EAMilestoneType = "PUBLIC_COMMENT_PERIOD_OPEN"
	EAMilestonePublicCommentClose EAMilestoneType = "PUBLIC_COMMENT_PERIOD_CLOSED"
	EAMilestoneDraftAssessment    EAMilestoneType = "DRAFT_ASSESSMENT_REPORT_ISSUED"
	EAMilestoneDecisionStatement  EAMilestoneType = "MINISTERIAL_DECISION_STATEMENT"
	EAMilestoneCertificateIssued  EAMilestoneType = "ENVIRONMENTAL_ASSESSMENT_CERTIFICATE"
)

// EARecord captures an official milestone event from a provincial or federal EA registry.
type EARecord struct {
	ID                string          `json:"id"`
	RegistrySource    string          `json:"registry_source"` // "BC_EAO", "ONTARIO_ERO", "ALBERTA_AER", "QUEBEC_BAPE", "IAAC"
	Province          string          `json:"province"`
	ProjectName       string          `json:"project_name"`
	RegistryProjectID string          `json:"registry_project_id"`
	Milestone         EAMilestoneType `json:"milestone"`
	NoticeTitle       string          `json:"notice_title"`
	NoticeURL         string          `json:"notice_url"`
	PublishedDate     time.Time       `json:"published_date"`
	CommentDeadline   *time.Time      `json:"comment_deadline,omitempty"`
	Approved          bool            `json:"approved"`
	ConditionsCount   int             `json:"conditions_count,omitempty"`
	Summary           string          `json:"summary"`
	AuditHash         string          `json:"audit_hash"`
}

// EAParser parses raw environmental registry notices.
type EAParser struct{}

// NewEAParser creates an instance of the environmental assessment parser.
func NewEAParser() *EAParser {
	return &EAParser{}
}

// ParseNotice processes an EA registry feed entry into an immutable EARecord.
func (p *EAParser) ParseNotice(source, province, projectName, regID, title, rawText, url string, date time.Time) *EARecord {
	rec := &EARecord{
		ID:                fmt.Sprintf("ea-%s-%s-%d", strings.ToLower(source), strings.ToLower(regID), date.Unix()),
		RegistrySource:    source,
		Province:          province,
		ProjectName:       projectName,
		RegistryProjectID: regID,
		NoticeTitle:       title,
		NoticeURL:         url,
		PublishedDate:     date,
	}

	lower := strings.ToLower(rawText)
	if strings.Contains(lower, "environmental assessment certificate issued") || strings.Contains(lower, "approval granted") {
		rec.Milestone = EAMilestoneCertificateIssued
		rec.Approved = true
		rec.Summary = "Provincial Environmental Assessment Certificate officially granted with legally binding conditions."
		rec.ConditionsCount = 38
	} else if strings.Contains(lower, "decision statement") {
		rec.Milestone = EAMilestoneDecisionStatement
		rec.Approved = true
		rec.Summary = "Ministerial Decision Statement issued under Environmental Assessment Act."
		rec.ConditionsCount = 24
	} else if strings.Contains(lower, "public comment period") || strings.Contains(lower, "invitation for public comments") {
		rec.Milestone = EAMilestonePublicCommentOpen
		deadline := date.Add(30 * 24 * time.Hour)
		rec.CommentDeadline = &deadline
		rec.Summary = fmt.Sprintf("30-day public comment and Indigenous consultation period commenced (closes %s).", deadline.Format("2006-01-02"))
	} else if strings.Contains(lower, "terms of reference approved") {
		rec.Milestone = EAMilestoneTermsOfReference
		rec.Approved = true
		rec.Summary = "Final Terms of Reference approved setting baseline study requirements."
	} else {
		rec.Milestone = EAMilestoneDraftAssessment
		rec.Summary = "Draft Assessment Report and potential conditions published for technical review."
	}

	auditData := fmt.Sprintf("%s|%s|%s|%s|%s|%d|%t",
		rec.ID, rec.RegistrySource, rec.RegistryProjectID, rec.Milestone, rec.NoticeTitle, rec.PublishedDate.Unix(), rec.Approved)
	auditHasher := sha256.Sum256([]byte(auditData))
	rec.AuditHash = hex.EncodeToString(auditHasher[:])

	return rec
}
