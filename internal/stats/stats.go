// Package stats measures what the relay is doing — events in, events matched,
// throughput, memory — and will expose it as a /stats JSON endpoint plus a
// live dashboard at /.
//
// SCAFFOLD: the Snapshot shape (the /stats contract) is sketched here; the
// Collector and HTTP handler are TODO.
package stats

// Snapshot is a point-in-time view of the relay's counters. The JSON tags are
// the /stats contract the dashboard and benchmark tooling will read.
type Snapshot struct {
	EventsInTotal int64   `json:"events_in_total"`
	EventsPerSec  float64 `json:"events_per_sec"`
	MatchedTotal  int64   `json:"matched_total"`
	MatchedPerSec float64 `json:"matched_per_sec"`
	FilterRatio   float64 `json:"filter_ratio"`
	Goroutines    int     `json:"goroutines"`
	HeapAllocMB   float64 `json:"heap_alloc_mb"`
	UptimeSeconds float64 `json:"uptime_seconds"`
}

// TODO(week 4):
//   - Collector: atomic counters for events in / matched, rate sampling.
//   - NewHandler(*Collector) http.Handler serving GET /stats (JSON) and
//     GET / (a //go:embed dashboard.html that polls /stats every second).
