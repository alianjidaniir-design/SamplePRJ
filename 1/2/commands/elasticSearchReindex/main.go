package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	indexName := flag.String("index", "studentss", "target Elasticsearch index")
	batchSize := flag.Int("batch", 500, "batch size for reindex")
	flag.Parse()

	start := time.Now()
	fmt.Printf("[es-reindex] start index=%s batch=%d\n", *indexName, *batchSize)

	// TODO: connect to Elasticsearch and run real reindex logic.
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("[es-reindex] done in %s\n", time.Since(start).String())
}
