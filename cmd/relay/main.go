// Command relay is the entry point: it will wire a firehose Source to the
// filter, the stats collector, and the store, and serve the /stats endpoint
// plus the live dashboard.
//
// SCAFFOLD: flags and shape are sketched; the wiring is TODO (weeks 4–5).
//
// Planned usage:
//
//	relay                                  connect to the live firehose
//	relay -replay testdata/sample.jsonl    replay a recorded sample (offline)
//	relay -addr :9090                       serve stats on a different port
package main

import (
	"flag"
	"log"
)

func main() {
	var (
		replay string
		addr   string
	)
	flag.StringVar(&replay, "replay", "", "path to a recorded JSONL sample; if empty, connect to the live firehose")
	flag.StringVar(&addr, "addr", ":8080", "listen address for the /stats HTTP server")
	flag.Parse()

	log.SetPrefix("relay: ")
	log.SetFlags(0)

	// TODO(weeks 4–5): wire up
	//   source  := firehose source (live, or ReplaySource when -replay set)
	//   filter  := firehose.NewFilter(firehose.DefaultCollections)
	//   collector, stats HTTP server on addr
	//   store   := batching store.Store
	// then source.Run(ctx, handle) with graceful shutdown on SIGINT/SIGTERM.
	log.Printf("scaffold only — not yet implemented (replay=%q addr=%q)", replay, addr)
}
