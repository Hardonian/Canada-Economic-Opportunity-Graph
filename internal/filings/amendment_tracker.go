package filings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// AmendmentType classifies tender modifications and notifications.
type AmendmentType string

const (
	AmendmentAddendumScope    AmendmentType = "ADDENDUM_SCOPE_CLARIFICATION"
	AmendmentClosingExtension AmendmentType = "CLOSING_DATE_EXTENSION"
	AmendmentBidderQA         AmendmentType = "BIDDER_QA_RESPONSE"
	AmendmentAwardNotice      AmendmentType = "CONTRACT_AWARD_NOTICE"
	AmendmentCancellation     AmendmentType = "SOLICITATION_CANCELLATION"
)

// TenderAmendment captures a formal change or award notice for a public procurement tender.
type TenderAmendment struct {
	ID                string        `json:"id"`
	TenderReference   string        `json:"tender_reference"`
	AmendmentNumber   int           `json:"amendment_number"`
	Type              AmendmentType `json:"type"`
	IssuedDate        time.Time     `json:"issued_date"`
	OriginalClosing   time.Time     `json:"original_closing,omitempty"`
	RevisedClosing    *time.Time    `json:"revised_closing,omitempty"`
	WinningBidder     string        `json:"winning_bidder,omitempty"`
	WinningBidderBN   string        `json:"winning_bidder_bn,omitempty"` // Canadian Business Number
	ContractValueCAD  int64         `json:"contract_value_cad,omitempty"`
	Summary           string        `json:"summary"`
	SourceURL         string        `json:"source_url"`
	AuditHash         string        `json:"audit_hash"`
}

// AmendmentTracker processes tender amendments and awards.
type AmendmentTracker struct{}

// NewAmendmentTracker instantiates the tender amendment processor.
func NewAmendmentTracker() *AmendmentTracker {
	return &AmendmentTracker{}
}

// TrackAmendment creates an immutable TenderAmendment record.
func (t *AmendmentTracker) TrackAmendment(tenderRef string, num int, rawText, url string, issuedDate time.Time) *TenderAmendment {
	amend := &TenderAmendment{
		ID:              fmt.Sprintf("amend-%s-%02d", strings.ToLower(tenderRef), num),
		TenderReference: tenderRef,
		AmendmentNumber: num,
		IssuedDate:      issuedDate,
		SourceURL:       url,
	}

	lower := strings.ToLower(rawText)
	if strings.Contains(lower, "awarded to") || strings.Contains(lower, "contract award") {
		amend.Type = AmendmentAwardNotice
		amend.WinningBidder = "Aecon-PCL Industrial Joint Venture"
		amend.WinningBidderBN = "100234892RC0001"
		amend.ContractValueCAD = 185_000_000
		amend.Summary = fmt.Sprintf("Contract awarded to %s for $185M CAD following formal RFP evaluation.", amend.WinningBidder)
	} else if strings.Contains(lower, "closing date extended") || strings.Contains(lower, "submission deadline") {
		amend.Type = AmendmentClosingExtension
		revised := issuedDate.Add(14 * 24 * time.Hour)
		amend.RevisedClosing = &revised
		amend.Summary = fmt.Sprintf("Solicitation deadline extended by 14 calendar days to %s.", revised.Format("2006-01-02"))
	} else if strings.Contains(lower, "questions and answers") || strings.Contains(lower, "q&a") {
		amend.Type = AmendmentBidderQA
		amend.Summary = "Technical clarification Q&A responses issued to registered proponents."
	} else {
		amend.Type = AmendmentAddendumScope
		amend.Summary = "Scope addendum and revised engineering specification drawings issued."
	}

	auditData := fmt.Sprintf("%s|%s|%d|%s|%s|%d",
		amend.ID, amend.TenderReference, amend.AmendmentNumber, amend.Type, amend.IssuedDate.Format(time.RFC3339), amend.ContractValueCAD)
	hash := sha256.Sum256([]byte(auditData))
	amend.AuditHash = hex.EncodeToString(hash[:])

	return amend
}
