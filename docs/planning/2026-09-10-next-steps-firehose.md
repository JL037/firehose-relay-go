# Next steps: the firehose ingestion slice (replay-first)

- **Date:** 2026-09-10
- **Status:** active
- **Phase:** first real code — start of the plan's weeks 4–5 project work,
  used as the vehicle for the weeks 1–2 Go fundamentals
- **Depends on:** the scaffold (`internal/firehose` stubs, `cmd/relay`)
- **Unblocks:** the store layer (week 3), the `/stats` collector, and the
  whole benchmark harness

## Why this phase, and why in this order

The scaffold defines the seams; this phase fills the first one end to end.
`internal/firehose` is the right place to start because it is the smallest
vertical slice that produces something runnable, and because it exercises
exactly the Go concepts the fundamentals phase is meant to build intuition on —
structs and JSON tags, interfaces, error handling, `bufio` scanning, and
(at the end) goroutines, channels, and `context` cancellation.

We build it **replay-first**: a file-backed `Source` before the live socket.
That ordering matters for three reasons:

1. **It is offline and deterministic.** You can run, test, and re-run the whole
   pipeline with no network and no rate limits.
2. **It is the benchmark's foundation.** The entire Node-vs-Go comparison rests
   on replaying one identical recording through both versions. The replay path
   is not a testing convenience bolted on later — it is a load-bearing part of
   the product, so it gets built first.
3. **It de-risks the live source.** By the time we dial the WebSocket, the
   parsing, filtering, and event-handling are already proven against known
   input; the only new variable is the transport.

The definition of done for the whole phase: `go run ./cmd/relay -replay=<file>`
reads a JSONL capture, filters it to the configured collections, prints matched
events, and exits cleanly — with tests covering the filter and the parser.

## Prerequisite: fundamentals checkpoints (weeks 1–2)

Don't treat these as a gate to clear before touching the repo — fold them in.
Each maps onto a step below. Consider the fundamentals "done" when you can:

- [ ] Read and write a Go struct with JSON tags without reaching for docs
      (→ Step 1, `Event`/`Commit`).
- [ ] Explain what an interface is and why `Source` is one (→ Step 3).
- [ ] Write a table-driven test with `t.Run` subtests (→ Step 2).
- [ ] Use `context.Context` to cancel a loop, and a `select` with `ctx.Done()`
      (→ Steps 3 and 5).
- [ ] Write a worker pool: a producer goroutine, N workers over a channel, and
      a clean shutdown (→ Step 5; the plan's "one throwaway concurrency program"
      is really this, done for real here).

Recommended reading, in order: [A Tour of Go](https://go.dev/tour),
[Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/) (the
"Concurrency", "Select", and "Context" chapters map directly onto Step 5).

## Decisions to lock before coding

1. **The 7 collection types.** `firehose.DefaultCollections` is empty. This is
   the one true blocker for Step 1. Recommended default until VoxPort's exact
   production set is confirmed:

   ```
   app.bsky.feed.post
   app.bsky.feed.repost
   app.bsky.feed.like
   app.bsky.graph.follow
   app.bsky.graph.block
   app.bsky.graph.listitem
   app.bsky.actor.profile
   ```

   > **OPEN:** confirm this matches production. The benchmark's "matched" number
   > is meaningless if the filter set differs from the Node pipeline's.

2. **Jetstream vs. raw relay for the live source.** Recommend Jetstream
   (JSON over WebSocket) — no CBOR decoding, and the server does server-side
   collection filtering via `wantedCollections` query params. The raw relay is
   a later optimization, not a starting point.

3. **WebSocket library.** Recommend `github.com/coder/websocket` (formerly
   `nhooyr.io/websocket`) — small, context-native API. This is the first entry
   in `go.mod`'s require block. (Alternative: `gorilla/websocket`, more common
   in older code but a heavier API.)

4. **Cursor semantics.** Jetstream lets you resume from a `cursor` (a
   `time_us`). Decide now that the live source tracks the last-seen `time_us` in
   memory so a reconnect resumes roughly where it left off. Durable cursor
   storage is out of scope for this phase.

## Build steps

Each step is a small, self-contained commit that builds and passes tests before
the next begins. This keeps the git history readable and each change reviewable.

### Step 1 — Pin the collection set and finish the types

**Files:** `internal/firehose/firehose.go`

- Fill in `DefaultCollections` with the 7 types from the decision above.
- Confirm `Event`/`Commit` cover what the filter and store need. At minimum the
  store will want `Did`, `Commit.Collection`, `Commit.Rkey`, `Commit.Operation`,
  and the raw `Commit.Record`. Add fields only when a consumer needs them —
  don't model the whole AT Proto schema speculatively.

**Done when:** the package compiles and `len(DefaultCollections) == 7`.

### Step 2 — Implement and test `Filter.Match`

**Files:** `internal/firehose/firehose.go`, `internal/firehose/filter_test.go` (new)

- `Match` returns true only for a `"commit"` event whose `Commit` is non-nil and
  whose `Collection` is in the configured set. Build the set into a
  `map[string]struct{}` once in `NewFilter`, not per call.
- Table-driven test covering: a matched collection, an unmatched collection, a
  non-`commit` kind, and a `"commit"` kind with a nil `Commit`. Add a guard test
  asserting the default set has exactly 7 entries — cheap insurance against an
  accidental edit silently changing what the benchmark measures.

**Done when:** `go test ./internal/firehose/` passes with those cases.

### Step 3 — Implement `ReplaySource`

**Files:** `internal/firehose/replay.go` (new)

- A `ReplaySource{ Path string }` implementing `Source`.
- `Run` opens the file, scans it line by line with `bufio.Scanner`, and — this
  is the real gotcha — **raises the scanner's buffer** with
  `scanner.Buffer(make([]byte, 0, 64<<10), 4<<20)`. The default 64 KB token
  limit will silently fail on large post records; a ~4 MB cap is safe.
- For each non-blank line: `json.Unmarshal` into an `Event`, then call `handle`.
  Skip malformed lines rather than aborting the whole replay (one corrupt record
  shouldn't kill a benchmark run) — but count skips and log the total at the end
  so corruption isn't invisible.
- Check `ctx` on each iteration via a non-blocking `select` and return
  `ctx.Err()` when cancelled.
- Test with a tiny in-repo fixture (a handful of lines, a couple matching): a
  known input yields a known matched count. This fixture is a test asset, not
  the benchmark capture — keep it small and synthetic.

**Done when:** a unit test drives a fixture file through `ReplaySource` + `Filter`
and asserts the matched count.

### Step 4 — Wire `cmd/relay` for the replay path

**Files:** `cmd/relay/main.go`

- Replace the "scaffold only" log with real wiring **for the replay path only**:
  - `ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` for
    clean Ctrl-C shutdown.
  - When `-replay` is set, construct a `ReplaySource`; otherwise keep returning
    the JetstreamSource's "not implemented" error (Step 5 lights it up).
  - `filter := firehose.NewFilter(firehose.DefaultCollections)`.
  - A `handle` closure that, for now, prints matched events as JSON lines to
    stdout. (This is a stand-in for the store; keep it dead simple until week 3.)
  - Run the source, and on return flush/shutdown cleanly.

**Done when:** `go run ./cmd/relay -replay=testdata/<fixture>.jsonl` prints the
matched events and exits 0.

### Step 5 — Implement the live `JetstreamSource`

**Files:** `internal/firehose/jetstream.go`, `go.mod`/`go.sum`

This is the concurrency step — the one worth slowing down on.

- Add the WebSocket dependency (`go get github.com/coder/websocket`).
- Build the subscribe URL from `Endpoint` + `wantedCollections[]=…` for each
  configured collection, so the server filters before bytes hit the wire.
- Read loop: read a frame, `json.Unmarshal` into `Event`, hand to `handle`.
  Record `time_us` as the resume cursor as you go.
- **Concurrency shape:** the read loop is the producer; push events onto a
  buffered channel; a small pool of worker goroutines drains it into `handle`.
  This is where a slow consumer (the DB, later) must not block the socket read —
  size the channel and decide the drop-vs-block policy deliberately, and write
  that decision down.
- **Reconnect with backoff:** on a dropped connection, reconnect with
  exponential backoff (cap ~30s), resuming from the last `time_us`. Respect
  `ctx` throughout — cancellation must unblock the read and stop the workers.

**Done when:** `go run ./cmd/relay` (no `-replay`) connects to the live firehose
and prints matched events; Ctrl-C shuts it down cleanly with no goroutine leak.

## Testing & verification

- **Unit:** filter cases (Step 2) and the replay-fixture count (Step 3) run in
  CI via the existing `go test ./...`. No network in unit tests — the live
  source is exercised manually, not in CI.
- **Manual smoke (replay):** run the relay against the fixture; eyeball that
  matched events print and non-matching collections are dropped.
- **Manual smoke (live):** run against Jetstream for ~30s; confirm events flow,
  then Ctrl-C and confirm a clean exit. Watch for goroutine leaks with a
  `-race` build during development (`go run -race ./cmd/relay`).
- **Race detector:** develop Step 5 with `-race` on; the producer/worker/channel
  code is exactly where data races hide.

## Risks & gotchas

- **`bufio.Scanner` token limit** (Step 3) — the single most likely silent bug;
  covered above.
- **Backpressure** (Step 5) — a slow `handle` blocking the socket read is the
  classic firehose failure. Decide the channel size and drop policy on purpose,
  not by accident.
- **Goroutine leaks on shutdown** (Step 5) — every goroutine must key off `ctx`;
  verify with `-race` and a manual Ctrl-C.
- **Filter-set drift** (Step 1) — if `DefaultCollections` ever diverges from the
  Node pipeline's set, every benchmark number is invalid. The count guard test
  is a tripwire, not a substitute for confirming the real list.

## What this phase deliberately excludes

- Postgres / `pgx` — that's week 3, behind the `store.Store` interface.
- The `/stats` collector and dashboard — the `handle` closure prints to stdout
  as a placeholder until then.
- `cmd/capture` — the real benchmark capture tool; it reuses the live source
  built here, so it comes after Step 5.
- Kubernetes, Prometheus, Grafana — week 6.

## Suggested commit sequence

1. `firehose: pin the 7 production collection types`
2. `firehose: implement and test Filter.Match`
3. `firehose: implement ReplaySource with tests`
4. `relay: wire up the replay path end to end`
5. `firehose: implement live Jetstream source with reconnect`

When this phase is done, log it in `docs/worklogs/` and record any wall you hit
in `docs/lessons/`.
