// Package firehose connects to and parses the Bluesky AT Protocol event
// stream (Jetstream), and filters it down to the collection types this relay
// cares about.
//
// SCAFFOLD: types and seams are sketched here; the implementations land during
// weeks 4–5 (see docs/planning). Nothing in this package is wired up yet.
package firehose

import (
	"context"
	"encoding/json"
)

// Event is one message from the Jetstream firehose. Fields will be filled in
// as the parser is built; kept minimal for now.
type Event struct {
	Did    string  `json:"did"`
	TimeUS int64   `json:"time_us"`
	Kind   string  `json:"kind"`
	Commit *Commit `json:"commit,omitempty"`
}

// Commit is the repository mutation carried by a "commit" event.
type Commit struct {
	Operation  string          `json:"operation"`
	Collection string          `json:"collection"`
	Rkey       string          `json:"rkey"`
	Record     json.RawMessage `json:"record,omitempty"`
}

// DefaultCollections will mirror the seven collection types VoxPort's
// production pipeline filters for. TODO: pin down the exact set.
var DefaultCollections = []string{
	// "app.bsky.feed.post",
	// ... (7 total)
}

// Source is anything that produces a stream of firehose Events: the live
// Jetstream socket, or a recorded replay file. TODO: implement both.
type Source interface {
	Run(ctx context.Context, handle func(Event) error) error
}
