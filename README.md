# firehose-relay-go

A Go rebuild of VoxPort's production AT Protocol firehose filter — the same job, the same 2,000+ events/sec input, the same 7 collection types — rewritten to put a real, measured number next to it.

> **Status: early build.** This is a from-scratch rewrite in progress. The benchmark numbers below are placeholders until `results/` has real data in it — see [Benchmarking](#benchmarking).

## What this is

VoxPort's ingestion pipeline filters the Bluesky firehose server-side, dropping ~96% of incoming events at the source before they ever touch the database. That pipeline runs in Node today. This repo is the same logic, same protocol, same filtering rules — written in Go — so the comparison between the two is an honest one: same input, same author, different language.

It's a controlled rewrite of something already running in production.

## Architecture

```
firehose-relay-go/
├── cmd/relay/main.go        entry point — wires everything together
├── internal/firehose/       connect + parse the AT Proto event stream
├── internal/store/          pgx writer, batching, the Postgres side
├── internal/stats/          the /stats endpoint + live dashboard
├── testdata/                recorded firehose samples for reproducible benchmarking
├── results/                 benchmark output (go.json, node.json)
├── deploy/                  Kubernetes manifests + Grafana dashboard
├── go.mod
└── .github/workflows/ci.yml
```

## Quickstart

Everything here runs locally — no cloud account, no hosted database.

```sh
# start local Postgres
docker compose up -d postgres

# run the relay against the live firehose
go run ./cmd/relay

# ...or replay a recorded sample instead (no network required)
go run ./cmd/relay -replay=testdata/sample.jsonl
```

Once it's running:

- `GET /stats` — JSON: events/sec, matched/sec, filter ratio, goroutine count, memory
- `GET /` — a live dashboard (polls `/stats` every second)

## Why Go

Three things pointed here rather than staying in Node:

1. **Go is the native language of this ecosystem.** Bluesky's own reference implementation — indigo, the PDS, the relay — is written in Go. This isn't an arbitrary language choice for the domain.
2. **Concurrency maps directly onto the problem.** Filtering a 2,000+ events/sec stream by collection type is a natural fit for goroutines + channels — one of the concrete things this rewrite is meant to test.
3. **It closes a real skill gap.** It's the language nearly every database/infra-adjacent role wanted that this project's stack didn't already cover.

## Benchmarking

The comparison against the Node version only means something if both sides see identical input. The methodology:

1. Capture one real slice of firehose traffic once, checked into `testdata/sample.jsonl`.
2. Replay that exact recording through both the Node version and this one.
3. Save each run's results to `results/node.json` and `results/go.json`.
4. Generate the table below from those files — not eyeballed.

| Metric            | Node (production) | Go (this repo) |
| ----------------- | ----------------- | -------------- |
| Events handled    | 2,000+/sec        | TBD            |
| DB write reduction| 96%               | TBD            |
| Memory footprint  | unmeasured        | TBD            |

## Monitoring

- **Local dev:** the built-in `/stats` dashboard above — no extra infrastructure.
- **Under Kubernetes (`deploy/`):** Prometheus scrapes `/metrics`, with the Grafana dashboard in `deploy/grafana-dashboard.json`.

## Roadmap

- [ ] Consume the firehose, filter to the 7 production collection types
- [ ] Batch-write matched events to Postgres via pgx
- [ ] `/stats` endpoint + live dashboard
- [ ] Recorded-sample benchmark harness vs. the Node version
- [ ] Containerize, deploy to a local kind cluster
- [ ] Prometheus + Grafana monitoring

## License

MIT — see [LICENSE](LICENSE).

## Author

Jared Lemler — founder/engineer, VoxPort · [@JL037](https://github.com/JL037)
