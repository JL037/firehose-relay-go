# Six-week Go + Postgres plan

- **Date:** 2026-09-10
- **Status:** active

## Goal

Close the two skill gaps every database/infra-adjacent role in this search
wanted — **Go** and **Kubernetes** — using a project that isn't a toy: a Go
rewrite of VoxPort's production Bluesky firehose filter, with real before/after
numbers against the Node version.

Source: the "Go + Postgres, Six Weeks" plan artifact.

## Approach

Six weeks, part-time, four phases:

1. **Weeks 1–2 — Go fundamentals.** A Tour of Go + Learn Go with Tests. One
   throwaway goroutines/channels program (a worker pool) to prove the
   concurrency model clicked.
2. **Week 3 — Go ↔ Postgres.** `pgx` + `pgxpool` connection pooling. Batch-
   insert 10k rows two ways (one-by-one vs. `COPY`) and time both. Lands in
   `internal/store`.
3. **Weeks 4–5 — the project: Firehose Relay.** Consume Jetstream, filter to
   the 7 production collection types, batch-write matches to Postgres, expose
   `/stats` + live dashboard, then the record-and-replay benchmark harness vs.
   the Node version.
4. **Week 6 — Kubernetes.** `kind` cluster, containerize, Deployment + Service,
   kill a pod and watch it reschedule. Stretch: Prometheus + Grafana.

## Benchmark methodology

Capture one real firehose slice once (`testdata/sample.jsonl`), replay that
identical file through both versions, save `results/{node,go}.json`, generate
the README table from those files.

| Metric             | Node (production) | Go (target)              |
| ------------------ | ----------------- | ------------------------ |
| Events handled     | 2,000+/sec        | measure & publish        |
| DB write reduction | 96%               | match or beat            |
| Memory footprint   | unmeasured        | goroutines vs. event loop|

## Open questions

- The exact 7 collection types the production filter uses (fill in
  `firehose.DefaultCollections`).
- Live source: Jetstream (JSON/WS) to start; raw relay later?
- `COPY` vs. batched `INSERT` — decide after the week-3 timing.
