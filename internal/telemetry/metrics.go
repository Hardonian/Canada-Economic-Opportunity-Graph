package telemetry

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector tracks operational SLIs and exports Prometheus telemetry.
type MetricsCollector struct {
	httpRequestsTotal       sync.Map // key: "METHOD /path status" -> *atomic.Uint64
	httpRequestDurations    sync.Map // key: "/path" -> *durationStats
	storageWALWrites        atomic.Uint64
	storageCheckpoints      atomic.Uint64
	verifierAttestations    atomic.Uint64
	verifierQuorumsReached  atomic.Uint64
	lakehouseQueries        atomic.Uint64
	activeAuditorPeers      atomic.Int64
}

type durationStats struct {
	count atomic.Uint64
	sumMs atomic.Uint64 // total milliseconds
}

// Global default telemetry collector.
var DefaultCollector = NewMetricsCollector()

// NewMetricsCollector constructs a Prometheus metrics registry.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

// IncHTTP increments the total request counter for a route.
func (mc *MetricsCollector) IncHTTP(method, path string, status int, duration time.Duration) {
	key := fmt.Sprintf("%s %s %d", method, path, status)
	val, _ := mc.httpRequestsTotal.LoadOrStore(key, &atomic.Uint64{})
	counter := val.(*atomic.Uint64)
	counter.Add(1)

	// Record duration
	dVal, _ := mc.httpRequestDurations.LoadOrStore(path, &durationStats{})
	stats := dVal.(*durationStats)
	stats.count.Add(1)
	stats.sumMs.Add(uint64(duration.Milliseconds()))
}

// IncWALWrite tracks durable transactions appended to disk.
func (mc *MetricsCollector) IncWALWrite() {
	mc.storageWALWrites.Add(1)
}

// IncCheckpoint tracks state snapshots compressed and synced.
func (mc *MetricsCollector) IncCheckpoint() {
	mc.storageCheckpoints.Add(1)
}

// IncAttestation records a signed milestone verification.
func (mc *MetricsCollector) IncAttestation(quorumReached bool) {
	mc.verifierAttestations.Add(1)
	if quorumReached {
		mc.verifierQuorumsReached.Add(1)
	}
}

// SetAuditorPeers records the number of active provincial auditor nodes.
func (mc *MetricsCollector) SetAuditorPeers(n int) {
	mc.activeAuditorPeers.Store(int64(n))
}

// IncLakehouseQuery tracks OLAP queries executed over the zero-copy engine.
func (mc *MetricsCollector) IncLakehouseQuery() {
	mc.lakehouseQueries.Add(1)
}

// RenderPrometheus generates standard Prometheus exposition text format.
func (mc *MetricsCollector) RenderPrometheus() string {
	var sb strings.Builder

	sb.WriteString("# HELP cog_http_requests_total Total number of HTTP requests processed by endpoint\n")
	sb.WriteString("# TYPE cog_http_requests_total counter\n")
	mc.httpRequestsTotal.Range(func(key, value any) bool {
		parts := strings.Split(key.(string), " ")
		if len(parts) == 3 {
			count := value.(*atomic.Uint64).Load()
			sb.WriteString(fmt.Sprintf("cog_http_requests_total{method=\"%s\",path=\"%s\",status=\"%s\"} %d\n",
				parts[0], parts[1], parts[2], count))
		}
		return true
	})

	sb.WriteString("\n# HELP cog_storage_wal_writes_total Total number of transactions committed to Write-Ahead Log\n")
	sb.WriteString("# TYPE cog_storage_wal_writes_total counter\n")
	sb.WriteString(fmt.Sprintf("cog_storage_wal_writes_total %d\n", mc.storageWALWrites.Load()))

	sb.WriteString("\n# HELP cog_storage_checkpoints_total Total number of state snapshots compressed and synced\n")
	sb.WriteString("# TYPE cog_storage_checkpoints_total counter\n")
	sb.WriteString(fmt.Sprintf("cog_storage_checkpoints_total %d\n", mc.storageCheckpoints.Load()))

	sb.WriteString("\n# HELP cog_verifier_attestations_total Total number of signed milestone attestations\n")
	sb.WriteString("# TYPE cog_verifier_attestations_total counter\n")
	sb.WriteString(fmt.Sprintf("cog_verifier_attestations_total %d\n", mc.verifierAttestations.Load()))

	sb.WriteString("\n# HELP cog_verifier_quorums_reached_total Total number of milestone quorums successfully verified\n")
	sb.WriteString("# TYPE cog_verifier_quorums_reached_total counter\n")
	sb.WriteString(fmt.Sprintf("cog_verifier_quorums_reached_total %d\n", mc.verifierQuorumsReached.Load()))

	sb.WriteString("\n# HELP cog_active_auditor_peers Current number of connected provincial auditor nodes\n")
	sb.WriteString("# TYPE cog_active_auditor_peers gauge\n")
	sb.WriteString(fmt.Sprintf("cog_active_auditor_peers %d\n", mc.activeAuditorPeers.Load()))

	sb.WriteString("\n# HELP cog_lakehouse_queries_total Total zero-copy lakehouse OLAP queries executed\n")
	sb.WriteString("# TYPE cog_lakehouse_queries_total counter\n")
	sb.WriteString(fmt.Sprintf("cog_lakehouse_queries_total %d\n", mc.lakehouseQueries.Load()))

	return sb.String()
}

// Handler returns an http.Handler that serves the /metrics Prometheus scrape endpoint.
func (mc *MetricsCollector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mc.RenderPrometheus()))
	})
}
