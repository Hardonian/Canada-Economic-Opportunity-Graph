package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/dealroom"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/security"
)

func signedDealRoomRequest(t *testing.T, method, target string, body *bytes.Reader, actor dealroom.Actor, secret string) *http.Request {
	t.Helper()
	var request *http.Request
	if body == nil {
		request = httptest.NewRequest(method, target, nil)
	} else {
		request = httptest.NewRequest(method, target, body)
	}
	payload, err := json.Marshal(actor)
	if err != nil {
		t.Fatal(err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "." + encoded))
	request.Header.Set(dealRoomActorHeader, encoded)
	request.Header.Set(dealRoomTimestampHeader, timestamp)
	request.Header.Set(dealRoomSignatureHeader, base64.RawURLEncoding.EncodeToString(mac.Sum(nil)))
	return request
}

func TestDealRoomGatewayEnforcesTenantIsolationAndEncryptedUpload(t *testing.T) {
	secret := strings.Repeat("gateway-secret-", 3)
	rooms, err := dealroom.NewStore(t.TempDir(), strings.Repeat("encryption-secret-", 3))
	if err != nil {
		t.Fatal(err)
	}
	options := DefaultOptions()
	options.DealRoom = rooms
	options.DealRoomGatewaySecret = secret
	options.MaxDealRoomDocumentBytes = 1024 * 1024
	server, err := NewServerWithOptions(database.NewMemoryStore(), options)
	if err != nil {
		t.Fatal(err)
	}
	owner := dealroom.Actor{UserID: "auth0|owner", TenantID: "tenant-a", Organization: "Tenant A", Clearance: security.ClassProtectedB, Citizenship: "CA"}
	createBody := bytes.NewReader([]byte(`{"project_id":"project-a","name":"Financing diligence","classification":"PROTECTED_A","retention_days":365}`))
	createRequest := signedDealRoomRequest(t, http.MethodPost, "/api/v1/deal-rooms", createBody, owner, secret)
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	server.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create room: %d %s", createResponse.Code, createResponse.Body.String())
	}
	var created struct {
		Room dealroom.Room `json:"room"`
	}
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	var multipartBody bytes.Buffer
	writer := multipart.NewWriter(&multipartBody)
	_ = writer.WriteField("title", "Private investment memo")
	_ = writer.WriteField("classification", "PROTECTED_A")
	file, err := writer.CreateFormFile("file", "memo.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.Write([]byte("private diligence material"))
	_ = writer.Close()
	uploadRequest := signedDealRoomRequest(t, http.MethodPost, "/api/v1/deal-rooms/"+created.Room.ID+"/documents", bytes.NewReader(multipartBody.Bytes()), owner, secret)
	uploadRequest.Header.Set("Content-Type", writer.FormDataContentType())
	uploadResponse := httptest.NewRecorder()
	server.ServeHTTP(uploadResponse, uploadRequest)
	if uploadResponse.Code != http.StatusCreated {
		t.Fatalf("upload document: %d %s", uploadResponse.Code, uploadResponse.Body.String())
	}
	var uploaded struct {
		Document dealroom.Document `json:"document"`
	}
	if err := json.Unmarshal(uploadResponse.Body.Bytes(), &uploaded); err != nil {
		t.Fatal(err)
	}

	intruder := dealroom.Actor{UserID: "auth0|intruder", TenantID: "tenant-b", Organization: "Tenant B", Clearance: security.ClassProtectedB, Citizenship: "CA"}
	readRequest := signedDealRoomRequest(t, http.MethodGet, "/api/v1/deal-rooms/"+created.Room.ID+"/documents/"+uploaded.Document.ID, nil, intruder, secret)
	readResponse := httptest.NewRecorder()
	server.ServeHTTP(readResponse, readRequest)
	if readResponse.Code != http.StatusForbidden {
		t.Fatalf("expected cross-tenant denial, got %d %s", readResponse.Code, readResponse.Body.String())
	}

	ownerRead := signedDealRoomRequest(t, http.MethodGet, "/api/v1/deal-rooms/"+created.Room.ID+"/documents/"+uploaded.Document.ID, nil, owner, secret)
	ownerResponse := httptest.NewRecorder()
	server.ServeHTTP(ownerResponse, ownerRead)
	if ownerResponse.Code != http.StatusOK || ownerResponse.Body.String() != "private diligence material" {
		t.Fatalf("owner download failed: %d %q", ownerResponse.Code, ownerResponse.Body.String())
	}
}
