// Command capture will record a fixed slice of live firehose traffic to a
// JSONL file — the frozen input the Node-vs-Go benchmark replays through both
// versions.
//
// SCAFFOLD: flags are sketched; implementation is TODO (week 4), blocked on
// the live Jetstream source in internal/firehose.
//
// Planned usage:
//
//	capture -duration 10m -out testdata/sample.jsonl
package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	var (
		duration time.Duration
		out      string
	)
	flag.DurationVar(&duration, "duration", 10*time.Minute, "how long to record")
	flag.StringVar(&out, "out", "testdata/sample.jsonl", "output JSONL path")
	flag.Parse()

	log.SetPrefix("capture: ")
	log.SetFlags(0)
	log.Printf("scaffold only — not yet implemented (duration=%s out=%q)", duration, out)
}
