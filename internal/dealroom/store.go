// Package dealroom provides tenant-isolated, encrypted deal-room persistence.
// It deliberately lives outside the public opportunity graph: no private room
// metadata, comments, or document content is projected into public APIs.
package dealroom

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("deal room resource not found")
	ErrForbidden = errors.New("deal room access denied")
)

type Role string

const (
	RoleOwner  Role = "OWNER"
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
	RoleViewer Role = "VIEWER"
)

func (r Role) CanWrite() bool  { return r == RoleOwner || r == RoleAdmin || r == RoleMember }
func (r Role) CanManage() bool { return r == RoleOwner || r == RoleAdmin }
func (r Role) Valid() bool {
	return r == RoleOwner || r == RoleAdmin || r == RoleMember || r == RoleViewer
}

type Actor struct {
	UserID       string                       `json:"user_id"`
	TenantID     string                       `json:"tenant_id"`
	Organization string                       `json:"organization"`
	Clearance    security.ClassificationLevel `json:"clearance"`
	Citizenship  string                       `json:"citizenship"`
}

type Membership struct {
	UserID    string    `json:"user_id"`
	Role      Role      `json:"role"`
	AddedBy   string    `json:"added_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID             string                       `json:"id"`
	TenantID       string                       `json:"tenant_id"`
	ProjectID      string                       `json:"project_id"`
	Name           string                       `json:"name"`
	Classification security.ClassificationLevel `json:"classification"`
	RetentionDays  int                          `json:"retention_days"`
	CreatedBy      string                       `json:"created_by"`
	CreatedAt      time.Time                    `json:"created_at"`
	UpdatedAt      time.Time                    `json:"updated_at"`
	Members        map[string]Membership        `json:"members"`
}

type Document struct {
	ID             string                       `json:"id"`
	RoomID         string                       `json:"room_id"`
	TenantID       string                       `json:"tenant_id"`
	Title          string                       `json:"title"`
	FileName       string                       `json:"file_name"`
	ContentType    string                       `json:"content_type"`
	SizeBytes      int64                        `json:"size_bytes"`
	SHA256         string                       `json:"sha256"`
	Classification security.ClassificationLevel `json:"classification"`
	UploadedBy     string                       `json:"uploaded_by"`
	CreatedAt      time.Time                    `json:"created_at"`
	RetentionAt    time.Time                    `json:"retention_at"`
}

type Comment struct {
	ID         string    `json:"id"`
	RoomID     string    `json:"room_id"`
	DocumentID string    `json:"document_id,omitempty"`
	Body       string    `json:"body"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditEvent struct {
	ID         string         `json:"id"`
	TenantID   string         `json:"tenant_id"`
	RoomID     string         `json:"room_id"`
	ActorID    string         `json:"actor_id"`
	Action     string         `json:"action"`
	SubjectID  string         `json:"subject_id"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	OccurredAt time.Time      `json:"occurred_at"`
}

type state struct {
	Rooms     map[string]*Room     `json:"rooms"`
	Documents map[string]*Document `json:"documents"`
	Comments  map[string][]Comment `json:"comments"`
	Audit     []AuditEvent         `json:"audit"`
}

type Store struct {
	mu        sync.RWMutex
	root      string
	documents string
	statePath string
	aead      cipher.AEAD
	state     state
}

func NewStore(root, encryptionSecret string) (*Store, error) {
	if len(strings.TrimSpace(encryptionSecret)) < 32 {
		return nil, fmt.Errorf("deal room encryption key must be at least 32 characters")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, fmt.Errorf("create deal room storage: %w", err)
	}
	documents := filepath.Join(root, "documents")
	if err := os.MkdirAll(documents, 0700); err != nil {
		return nil, fmt.Errorf("create deal room document storage: %w", err)
	}
	key := sha256.Sum256([]byte(encryptionSecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	s := &Store{
		root: root, documents: documents, statePath: filepath.Join(root, "metadata.enc"), aead: aead,
		state: state{Rooms: map[string]*Room{}, Documents: map[string]*Document{}, Comments: map[string][]Comment{}, Audit: []AuditEvent{}},
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	contents, err := os.ReadFile(s.statePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read deal room metadata: %w", err)
	}
	plaintext, err := s.open(contents)
	if err != nil {
		return fmt.Errorf("decrypt deal room metadata: %w", err)
	}
	if err := json.Unmarshal(plaintext, &s.state); err != nil {
		return fmt.Errorf("decode deal room metadata: %w", err)
	}
	if s.state.Rooms == nil {
		s.state.Rooms = map[string]*Room{}
	}
	if s.state.Documents == nil {
		s.state.Documents = map[string]*Document{}
	}
	if s.state.Comments == nil {
		s.state.Comments = map[string][]Comment{}
	}
	if s.state.Audit == nil {
		s.state.Audit = []AuditEvent{}
	}
	return nil
}

func (s *Store) seal(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return append(nonce, s.aead.Seal(nil, nonce, plaintext, nil)...), nil
}

func (s *Store) open(payload []byte) ([]byte, error) {
	if len(payload) < s.aead.NonceSize() {
		return nil, errors.New("encrypted payload is malformed")
	}
	return s.aead.Open(nil, payload[:s.aead.NonceSize()], payload[s.aead.NonceSize():], nil)
}

func (s *Store) persistLocked() error {
	payload, err := json.Marshal(s.state)
	if err != nil {
		return err
	}
	sealed, err := s.seal(payload)
	if err != nil {
		return err
	}
	temporary := s.statePath + "." + uuid.NewString() + ".tmp"
	if err := os.WriteFile(temporary, sealed, 0600); err != nil {
		return err
	}
	if err := os.Rename(temporary, s.statePath); err != nil {
		_ = os.Remove(temporary)
		return err
	}
	return nil
}

func validClassification(level security.ClassificationLevel) bool {
	_, ok := security.ClearanceRank[level]
	return ok
}

func validIdentifier(value string) bool {
	if value == "" || len(value) > 200 {
		return false
	}
	for _, char := range value {
		if unicode.IsLetter(char) || unicode.IsDigit(char) || strings.ContainsRune("._|:-", char) {
			continue
		}
		return false
	}
	return true
}

func cleanText(value string, max int) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) > max {
		value = value[:max]
	}
	return value
}

func (s *Store) auditLocked(actor Actor, roomID, action, subjectID string, metadata map[string]any) {
	s.state.Audit = append(s.state.Audit, AuditEvent{ID: uuid.NewString(), TenantID: actor.TenantID, RoomID: roomID, ActorID: actor.UserID, Action: action, SubjectID: subjectID, Metadata: metadata, OccurredAt: time.Now().UTC()})
	if len(s.state.Audit) > 5_000 {
		s.state.Audit = s.state.Audit[len(s.state.Audit)-5_000:]
	}
}

func (s *Store) roomForActorLocked(actor Actor, roomID string, write, manage bool) (*Room, Membership, error) {
	room := s.state.Rooms[roomID]
	if room == nil {
		return nil, Membership{}, ErrNotFound
	}
	if room.TenantID != actor.TenantID || !validIdentifier(actor.UserID) || !validIdentifier(actor.TenantID) {
		return nil, Membership{}, ErrForbidden
	}
	member, ok := room.Members[actor.UserID]
	if !ok || (write && !member.Role.CanWrite()) || (manage && !member.Role.CanManage()) {
		return nil, Membership{}, ErrForbidden
	}
	if security.ClearanceRank[actor.Clearance] < security.ClearanceRank[room.Classification] {
		return nil, Membership{}, ErrForbidden
	}
	return room, member, nil
}

func cloneRoom(room *Room) *Room {
	clone := *room
	clone.Members = make(map[string]Membership, len(room.Members))
	for id, member := range room.Members {
		clone.Members[id] = member
	}
	return &clone
}

func cloneDocument(document *Document) *Document { clone := *document; return &clone }

func (s *Store) CreateRoom(actor Actor, projectID, name string, classification security.ClassificationLevel, retentionDays int) (*Room, error) {
	if !validIdentifier(actor.UserID) || !validIdentifier(actor.TenantID) || !validIdentifier(projectID) || !validClassification(classification) {
		return nil, errors.New("invalid deal room attributes")
	}
	name = cleanText(name, 160)
	if name == "" {
		return nil, errors.New("deal room name is required")
	}
	if retentionDays < 30 || retentionDays > 3650 {
		return nil, errors.New("retention days must be between 30 and 3650")
	}
	if security.ClearanceRank[actor.Clearance] < security.ClearanceRank[classification] {
		return nil, ErrForbidden
	}
	now := time.Now().UTC()
	room := &Room{ID: uuid.NewString(), TenantID: actor.TenantID, ProjectID: projectID, Name: name, Classification: classification, RetentionDays: retentionDays, CreatedBy: actor.UserID, CreatedAt: now, UpdatedAt: now, Members: map[string]Membership{actor.UserID: {UserID: actor.UserID, Role: RoleOwner, AddedBy: actor.UserID, CreatedAt: now}}}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Rooms[room.ID] = room
	s.auditLocked(actor, room.ID, "ROOM_CREATED", room.ID, map[string]any{"project_id": projectID, "classification": classification})
	if err := s.persistLocked(); err != nil {
		delete(s.state.Rooms, room.ID)
		return nil, err
	}
	return cloneRoom(room), nil
}

func (s *Store) ListRooms(actor Actor) ([]*Room, error) {
	if !validIdentifier(actor.UserID) || !validIdentifier(actor.TenantID) {
		return nil, ErrForbidden
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	rooms := make([]*Room, 0)
	for _, room := range s.state.Rooms {
		if room.TenantID != actor.TenantID {
			continue
		}
		if _, member := room.Members[actor.UserID]; !member {
			continue
		}
		if security.ClearanceRank[actor.Clearance] < security.ClearanceRank[room.Classification] {
			continue
		}
		rooms = append(rooms, cloneRoom(room))
	}
	sort.Slice(rooms, func(i, j int) bool { return rooms[i].UpdatedAt.After(rooms[j].UpdatedAt) })
	return rooms, nil
}

func (s *Store) GetRoom(actor Actor, roomID string) (*Room, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	room, _, err := s.roomForActorLocked(actor, roomID, false, false)
	if err != nil {
		return nil, err
	}
	return cloneRoom(room), nil
}

func (s *Store) AddMember(actor Actor, roomID, userID string, role Role) (*Membership, error) {
	if !validIdentifier(userID) || !role.Valid() || role == RoleOwner {
		return nil, errors.New("invalid member")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	room, _, err := s.roomForActorLocked(actor, roomID, false, true)
	if err != nil {
		return nil, err
	}
	member := Membership{UserID: userID, Role: role, AddedBy: actor.UserID, CreatedAt: time.Now().UTC()}
	room.Members[userID] = member
	room.UpdatedAt = time.Now().UTC()
	s.auditLocked(actor, room.ID, "MEMBER_GRANTED", userID, map[string]any{"role": role})
	if err := s.persistLocked(); err != nil {
		return nil, err
	}
	return &member, nil
}

func (s *Store) UploadDocument(actor Actor, roomID, title, fileName, contentType string, classification security.ClassificationLevel, body io.Reader, maximumBytes int64) (*Document, error) {
	if maximumBytes <= 0 || !validClassification(classification) {
		return nil, errors.New("invalid document upload")
	}
	title, fileName, contentType = cleanText(title, 180), cleanText(fileName, 180), cleanText(contentType, 120)
	if title == "" || fileName == "" || contentType == "" {
		return nil, errors.New("title, file name, and content type are required")
	}
	if !strings.HasPrefix(contentType, "application/pdf") && !strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument") && !strings.HasPrefix(contentType, "text/plain") && !strings.HasPrefix(contentType, "text/csv") {
		return nil, errors.New("document content type is not allowed")
	}
	contents, err := io.ReadAll(io.LimitReader(body, maximumBytes+1))
	if err != nil {
		return nil, err
	}
	if len(contents) == 0 || int64(len(contents)) > maximumBytes {
		return nil, errors.New("document is empty or exceeds the configured limit")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	room, _, err := s.roomForActorLocked(actor, roomID, true, false)
	if err != nil {
		return nil, err
	}
	if security.ClearanceRank[actor.Clearance] < security.ClearanceRank[classification] || security.ClearanceRank[classification] < security.ClearanceRank[room.Classification] {
		return nil, ErrForbidden
	}
	document := &Document{ID: uuid.NewString(), RoomID: room.ID, TenantID: room.TenantID, Title: title, FileName: fileName, ContentType: contentType, SizeBytes: int64(len(contents)), SHA256: fmt.Sprintf("%x", sha256.Sum256(contents)), Classification: classification, UploadedBy: actor.UserID, CreatedAt: time.Now().UTC(), RetentionAt: time.Now().UTC().AddDate(0, 0, room.RetentionDays)}
	sealed, err := s.seal(contents)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(s.documents, document.ID+".enc")
	if err := os.WriteFile(path, sealed, 0600); err != nil {
		return nil, err
	}
	s.state.Documents[document.ID] = document
	room.UpdatedAt = time.Now().UTC()
	s.auditLocked(actor, room.ID, "DOCUMENT_UPLOADED", document.ID, map[string]any{"sha256": document.SHA256, "classification": classification, "size_bytes": document.SizeBytes})
	if err := s.persistLocked(); err != nil {
		delete(s.state.Documents, document.ID)
		_ = os.Remove(path)
		return nil, err
	}
	return cloneDocument(document), nil
}

func (s *Store) ListDocuments(actor Actor, roomID string) ([]*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, _, err := s.roomForActorLocked(actor, roomID, false, false); err != nil {
		return nil, err
	}
	documents := make([]*Document, 0)
	for _, document := range s.state.Documents {
		if document.RoomID == roomID && security.ClearanceRank[actor.Clearance] >= security.ClearanceRank[document.Classification] {
			documents = append(documents, cloneDocument(document))
		}
	}
	sort.Slice(documents, func(i, j int) bool { return documents[i].CreatedAt.After(documents[j].CreatedAt) })
	return documents, nil
}

func (s *Store) ReadDocument(actor Actor, roomID, documentID string) (*Document, []byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room, _, err := s.roomForActorLocked(actor, roomID, false, false)
	if err != nil {
		return nil, nil, err
	}
	document := s.state.Documents[documentID]
	if document == nil || document.RoomID != room.ID || document.TenantID != actor.TenantID {
		return nil, nil, ErrNotFound
	}
	if security.ClearanceRank[actor.Clearance] < security.ClearanceRank[document.Classification] {
		return nil, nil, ErrForbidden
	}
	payload, err := os.ReadFile(filepath.Join(s.documents, document.ID+".enc"))
	if err != nil {
		return nil, nil, err
	}
	contents, err := s.open(payload)
	if err != nil {
		return nil, nil, err
	}
	s.auditLocked(actor, room.ID, "DOCUMENT_DOWNLOADED", document.ID, nil)
	if err := s.persistLocked(); err != nil {
		return nil, nil, err
	}
	return cloneDocument(document), contents, nil
}

func (s *Store) AddComment(actor Actor, roomID, documentID, body string) (*Comment, error) {
	body = cleanText(body, 2000)
	if body == "" {
		return nil, errors.New("comment is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	room, _, err := s.roomForActorLocked(actor, roomID, true, false)
	if err != nil {
		return nil, err
	}
	if documentID != "" {
		document := s.state.Documents[documentID]
		if document == nil || document.RoomID != room.ID {
			return nil, ErrNotFound
		}
	}
	comment := Comment{ID: uuid.NewString(), RoomID: room.ID, DocumentID: documentID, Body: body, CreatedBy: actor.UserID, CreatedAt: time.Now().UTC()}
	s.state.Comments[room.ID] = append(s.state.Comments[room.ID], comment)
	room.UpdatedAt = time.Now().UTC()
	s.auditLocked(actor, room.ID, "COMMENT_ADDED", comment.ID, map[string]any{"document_id": documentID})
	if err := s.persistLocked(); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *Store) ListComments(actor Actor, roomID string) ([]Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, _, err := s.roomForActorLocked(actor, roomID, false, false); err != nil {
		return nil, err
	}
	return append([]Comment(nil), s.state.Comments[roomID]...), nil
}

func (s *Store) ListAudit(actor Actor, roomID string) ([]AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, _, err := s.roomForActorLocked(actor, roomID, false, false); err != nil {
		return nil, err
	}
	items := make([]AuditEvent, 0)
	for _, event := range s.state.Audit {
		if event.RoomID == roomID && event.TenantID == actor.TenantID {
			items = append(items, event)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OccurredAt.After(items[j].OccurredAt) })
	return items, nil
}

// PurgeExpired enforces retention server-side. It is safe to call repeatedly.
// Expiry destroys every private record for an expired room, rather than
// retaining collaboration metadata after the room's documents are gone.
func (s *Store) PurgeExpired(now time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	removedRooms := 0
	for roomID, room := range s.state.Rooms {
		if now.Before(room.CreatedAt.AddDate(0, 0, room.RetentionDays)) {
			continue
		}
		for documentID, document := range s.state.Documents {
			if document.RoomID != roomID {
				continue
			}
			_ = os.Remove(filepath.Join(s.documents, documentID+".enc"))
			delete(s.state.Documents, documentID)
		}
		delete(s.state.Comments, roomID)
		delete(s.state.Rooms, roomID)
		removedRooms++
		s.state.Audit = append(s.state.Audit, AuditEvent{ID: uuid.NewString(), TenantID: room.TenantID, RoomID: roomID, ActorID: "retention-service", Action: "ROOM_PURGED_BY_RETENTION", SubjectID: roomID, OccurredAt: now.UTC()})
	}
	if removedRooms > 0 {
		return removedRooms, s.persistLocked()
	}
	return 0, nil
}
