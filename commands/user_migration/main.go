package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	fromVersion := flag.String("from", "v1", "source schema version")
	toVersion := flag.String("to", "v2", "target schema version")
	limit := flag.Int("limit", 1000, "max users per run")
	flag.Parse()

	start := time.Now()
	fmt.Printf("[user-migration] start from=%s to=%s limit=%d\n", *fromVersion, *toVersion, *limit)

	// TODO: load users, transform records, write migrated data.
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("[user-migration] done in %s\n", time.Since(start).String())
}
