// freshnesscheck verifies that the checked-in public release has not exceeded
// its declared operational freshness window. It is intentionally independent
// of live adapters so it can run in CI and alert before stale data is served.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type manifest struct {
	DatasetVersion string `json:"dataset_version"`
	GeneratedAt    string `json:"generated_at"`
}

func main() {
	manifestPath := flag.String("manifest", "data/public/manifest.json", "path to public release manifest")
	maxAge := flag.Duration("max-age", 30*24*time.Hour, "maximum permitted age of a public release")
	now := flag.String("now", "", "RFC3339 time used for deterministic verification (optional)")
	flag.Parse()

	contents, err := os.ReadFile(*manifestPath)
	if err != nil {
		fail("read manifest", err)
	}
	var release manifest
	if err := json.Unmarshal(contents, &release); err != nil {
		fail("parse manifest", err)
	}
	generatedAt, err := time.Parse(time.RFC3339, release.GeneratedAt)
	if err != nil {
		fail("parse generated_at", err)
	}
	checkedAt := time.Now().UTC()
	if *now != "" {
		checkedAt, err = time.Parse(time.RFC3339, *now)
		if err != nil {
			fail("parse now", err)
		}
	}
	age := checkedAt.Sub(generatedAt)
	if age < 0 {
		fail("validate generated_at", fmt.Errorf("release timestamp %s is in the future", generatedAt.Format(time.RFC3339)))
	}
	if age > *maxAge {
		fail("validate freshness", fmt.Errorf("release %s is %s old; maximum is %s", release.DatasetVersion, age.Round(time.Hour), maxAge.String()))
	}
	fmt.Printf("release %s is fresh: generated_at=%s age=%s max_age=%s\n", release.DatasetVersion, generatedAt.Format(time.RFC3339), age.Round(time.Hour), maxAge.String())
}

func fail(operation string, err error) {
	fmt.Fprintf(os.Stderr, "freshnesscheck: %s: %v\n", operation, err)
	os.Exit(1)
}
