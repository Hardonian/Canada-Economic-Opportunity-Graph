package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/dealroom"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
)

const (
	dealRoomActorHeader     = "X-COG-Actor"
	dealRoomSignatureHeader = "X-COG-Actor-Signature"
	dealRoomTimestampHeader = "X-COG-Actor-Timestamp"
)

func (s *Server) withDealRoomActor(w http.ResponseWriter, r *http.Request) (dealroom.Actor, bool) {
	if s.dealRoom == nil || len(s.dealRoomGatewaySecret) < 32 {
		writeError(w, r, http.StatusServiceUnavailable, "deal_room_not_configured", "Private deal rooms are not configured on this deployment.")
		return dealroom.Actor{}, false
	}
	encoded, signature, timestamp := r.Header.Get(dealRoomActorHeader), r.Header.Get(dealRoomSignatureHeader), r.Header.Get(dealRoomTimestampHeader)
	if len(encoded) == 0 || len(encoded) > 4096 || len(signature) == 0 || len(signature) > 128 {
		return dealRoomUnauthorizedWithWriter(w, r)
	}
	when, err := time.Parse(time.RFC3339, timestamp)
	if err != nil || time.Since(when) > 5*time.Minute || when.After(time.Now().Add(time.Minute)) {
		return dealRoomUnauthorizedWithWriter(w, r)
	}
	mac := hmac.New(sha256.New, []byte(s.dealRoomGatewaySecret))
	_, _ = mac.Write([]byte(timestamp + "." + encoded))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(signature), []byte(expected)) != 1 {
		return dealRoomUnauthorizedWithWriter(w, r)
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return dealRoomUnauthorizedWithWriter(w, r)
	}
	var actor dealroom.Actor
	if err := json.Unmarshal(payload, &actor); err != nil || actor.UserID == "" || actor.TenantID == "" || !validDealRoomClassification(actor.Clearance) {
		return dealRoomUnauthorizedWithWriter(w, r)
	}
	_, _ = s.dealRoom.PurgeExpired(time.Now().UTC())
	return actor, true
}

func dealRoomUnauthorizedWithWriter(w http.ResponseWriter, r *http.Request) (dealroom.Actor, bool) {
	writeError(w, r, http.StatusUnauthorized, "invalid_deal_room_actor", "A trusted authenticated actor assertion is required.")
	return dealroom.Actor{}, false
}

func validDealRoomClassification(level security.ClassificationLevel) bool {
	_, ok := security.ClearanceRank[level]
	return ok
}

func dealRoomStoreError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, dealroom.ErrNotFound):
		writeError(w, r, http.StatusNotFound, "deal_room_not_found", "The requested deal-room resource was not found.")
	case errors.Is(err, dealroom.ErrForbidden):
		writeError(w, r, http.StatusForbidden, "deal_room_forbidden", "You are not authorized to access this deal-room resource.")
	default:
		writeError(w, r, http.StatusBadRequest, "deal_room_invalid_request", "The deal-room request is invalid.")
	}
}

func decodeDealRoomJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if decoder.More() {
		return fmt.Errorf("request contains multiple JSON values")
	}
	return nil
}

func (s *Server) handleDealRooms(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	if r.Method == http.MethodGet {
		rooms, err := s.dealRoom.ListRooms(actor)
		if err != nil {
			dealRoomStoreError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"rooms": rooms})
		return
	}
	var request struct {
		ProjectID      string                       `json:"project_id"`
		Name           string                       `json:"name"`
		Classification security.ClassificationLevel `json:"classification"`
		RetentionDays  int                          `json:"retention_days"`
	}
	if err := decodeDealRoomJSON(r, &request); err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	room, err := s.dealRoom.CreateRoom(actor, request.ProjectID, request.Name, request.Classification, request.RetentionDays)
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"room": room})
}

func (s *Server) handleDealRoom(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	room, err := s.dealRoom.GetRoom(actor, r.PathValue("id"))
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"room": room})
}

func (s *Server) handleDealRoomMembers(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	var request struct {
		UserID string        `json:"user_id"`
		Role   dealroom.Role `json:"role"`
	}
	if err := decodeDealRoomJSON(r, &request); err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	member, err := s.dealRoom.AddMember(actor, r.PathValue("id"), request.UserID, request.Role)
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"member": member})
}

func (s *Server) handleDealRoomDocuments(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	roomID := r.PathValue("id")
	if r.Method == http.MethodGet {
		documents, err := s.dealRoom.ListDocuments(actor, roomID)
		if err != nil {
			dealRoomStoreError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"documents": documents})
		return
	}
	if err := r.ParseMultipartForm(s.maxDealRoomDocumentBytes); err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	defer file.Close()
	classification := security.ClassificationLevel(strings.TrimSpace(r.FormValue("classification")))
	contentType := header.Header.Get("Content-Type")
	if contentType == "" || strings.EqualFold(contentType, "application/octet-stream") {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(header.Filename)))
	}
	document, err := s.dealRoom.UploadDocument(actor, roomID, r.FormValue("title"), header.Filename, contentType, classification, file, s.maxDealRoomDocumentBytes)
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"document": document})
}

func (s *Server) handleDealRoomDocumentDownload(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	document, contents, err := s.dealRoom.ReadDocument(actor, r.PathValue("id"), r.PathValue("documentID"))
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", document.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=deal-room-document")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(contents)), 10))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(contents)
}

func (s *Server) handleDealRoomComments(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	roomID := r.PathValue("id")
	if r.Method == http.MethodGet {
		comments, err := s.dealRoom.ListComments(actor, roomID)
		if err != nil {
			dealRoomStoreError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"comments": comments})
		return
	}
	var request struct {
		DocumentID string `json:"document_id"`
		Body       string `json:"body"`
	}
	if err := decodeDealRoomJSON(r, &request); err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	comment, err := s.dealRoom.AddComment(actor, roomID, request.DocumentID, request.Body)
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"comment": comment})
}

func (s *Server) handleDealRoomAudit(w http.ResponseWriter, r *http.Request) {
	actor, ok := s.withDealRoomActor(w, r)
	if !ok {
		return
	}
	audit, err := s.dealRoom.ListAudit(actor, r.PathValue("id"))
	if err != nil {
		dealRoomStoreError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"audit": audit})
}
