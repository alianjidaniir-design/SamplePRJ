package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	period := flag.String("period", "daily", "stats update period: hourly|daily|weekly")
	dryRun := flag.Bool("dry-run", false, "run without writing updates")
	flag.Parse()

	start := time.Now()
	fmt.Printf("[stats-update] start period=%s dryRun=%t\n", *period, *dryRun)

	// TODO: aggregate statistics and persist results.
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("[stats-update] done in %s\n", time.Since(start).String())
}
