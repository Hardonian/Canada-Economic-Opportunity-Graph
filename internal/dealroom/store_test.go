package dealroom

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
)

func actor(user, tenant string) Actor {
	return Actor{UserID: user, TenantID: tenant, Organization: tenant, Clearance: security.ClassProtectedB, Citizenship: "CA"}
}

func TestTenantIsolationEncryptedDocumentsAndRetention(t *testing.T) {
	store, err := NewStore(t.TempDir(), strings.Repeat("x", 48))
	if err != nil {
		t.Fatal(err)
	}
	room, err := store.CreateRoom(actor("auth0|owner", "tenant-a"), "project-a", "Project A diligence", security.ClassProtectedA, 30)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.UploadDocument(actor("auth0|owner", "tenant-a"), room.ID, "Memo", "memo.txt", "text/plain", security.ClassProtectedA, strings.NewReader("confidential diligence"), 1024); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ListDocuments(actor("auth0|intruder", "tenant-b"), room.ID); err != ErrForbidden {
		t.Fatalf("expected tenant isolation denial, got %v", err)
	}
	documents, err := store.ListDocuments(actor("auth0|owner", "tenant-a"), room.ID)
	if err != nil || len(documents) != 1 {
		t.Fatalf("expected one document: %v", err)
	}
	_, contents, err := store.ReadDocument(actor("auth0|owner", "tenant-a"), room.ID, documents[0].ID)
	if err != nil || string(contents) != "confidential diligence" {
		t.Fatalf("unexpected read: %v %q", err, contents)
	}
	if documents[0].SHA256 == "" {
		t.Fatal("expected document checksum")
	}
}

func TestPurgeExpiredRemovesEntireRoomAndEncryptedContent(t *testing.T) {
	store, err := NewStore(t.TempDir(), strings.Repeat("x", 48))
	if err != nil {
		t.Fatal(err)
	}
	owner := actor("auth0|owner", "tenant-a")
	room, err := store.CreateRoom(owner, "project-a", "Expired diligence", security.ClassProtectedA, 30)
	if err != nil {
		t.Fatal(err)
	}
	document, err := store.UploadDocument(owner, room.ID, "Memo", "memo.txt", "text/plain", security.ClassProtectedA, strings.NewReader("expired diligence"), 1024)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddComment(owner, room.ID, document.ID, "remove with the room"); err != nil {
		t.Fatal(err)
	}

	store.state.Rooms[room.ID].CreatedAt = time.Now().UTC().AddDate(0, 0, -31)
	removed, err := store.PurgeExpired(time.Now().UTC())
	if err != nil || removed != 1 {
		t.Fatalf("expected one room to be purged, got %d: %v", removed, err)
	}
	if _, err := store.GetRoom(owner, room.ID); err != ErrNotFound {
		t.Fatalf("expected expired room to be unavailable, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(store.documents, document.ID+".enc")); !os.IsNotExist(err) {
		t.Fatalf("expected encrypted document to be deleted, got %v", err)
	}
	if _, exists := store.state.Comments[room.ID]; exists {
		t.Fatal("expected collaboration comments to be deleted")
	}
}
