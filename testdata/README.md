# testdata

Recorded firehose samples used as **frozen benchmark input** — both the Node
and Go versions replay the *identical* file so the comparison means something.

## Files (planned)

- **`sample.jsonl`** — the real benchmark capture: a ~10-minute slice of live
  traffic, produced once with `cmd/capture` and checked in for reproducibility.
  Not present yet (`cmd/capture` is week 4).

## Format

One [Jetstream](https://github.com/bluesky-social/jetstream) event per line —
the JSON-over-WebSocket view of the AT Protocol firehose. See
`internal/firehose/firehose.go` for the decoded shape.
