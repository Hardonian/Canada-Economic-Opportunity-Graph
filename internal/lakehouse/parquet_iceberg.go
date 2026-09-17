package lakehouse

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// IcebergTableManager manages Apache Iceberg table metadata and snapshot commits.
type IcebergTableManager struct {
	mu     sync.RWMutex
	tables map[string]*IcebergMetadata
}

// NewIcebergTableManager initializes the lakehouse metadata catalog.
func NewIcebergTableManager() *IcebergTableManager {
	return &IcebergTableManager{
		tables: make(map[string]*IcebergMetadata),
	}
}

// CreateTable creates a new Iceberg table specification.
func (itm *IcebergTableManager) CreateTable(tableName, location string, schema TableSchema, partSpec PartitionSpec) (*IcebergMetadata, error) {
	itm.mu.Lock()
	defer itm.mu.Unlock()

	if _, exists := itm.tables[tableName]; exists {
		return nil, fmt.Errorf("table %s already exists", tableName)
	}

	meta := &IcebergMetadata{
		TableUUID:       fmt.Sprintf("table-uuid-%s", tableName),
		FormatVersion:   2,
		Location:        location,
		CurrentSchemaID: schema.SchemaID,
		Schemas:         []TableSchema{schema},
		PartitionSpec:   partSpec,
		CurrentSnapshot: 0,
		Snapshots:       make([]Snapshot, 0),
	}

	itm.tables[tableName] = meta
	return meta, nil
}

// CommitSnapshot appends an immutable snapshot to the table history.
func (itm *IcebergTableManager) CommitSnapshot(tableName string, recordCount int64, manifestPath string) (*Snapshot, error) {
	itm.mu.Lock()
	defer itm.mu.Unlock()

	meta, ok := itm.tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s not found", tableName)
	}

	snapID := time.Now().UnixNano()
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%s", tableName, snapID, recordCount, manifestPath)))

	snap := Snapshot{
		SnapshotID:   snapID,
		ParentID:     meta.CurrentSnapshot,
		TimestampMs:  time.Now().UnixMilli(),
		ManifestList: manifestPath,
		RecordCount:  recordCount,
		ContentHash:  hex.EncodeToString(h[:]),
	}

	meta.Snapshots = append(meta.Snapshots, snap)
	meta.CurrentSnapshot = snapID

	return &snap, nil
}

// GetMetadata retrieves current Iceberg metadata for a table.
func (itm *IcebergTableManager) GetMetadata(tableName string) (*IcebergMetadata, error) {
	itm.mu.RLock()
	defer itm.mu.RUnlock()

	meta, ok := itm.tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s not found", tableName)
	}
	return meta, nil
}
