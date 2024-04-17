package main

// ./bin/wof-list-updates -max-commits 4 -org whosonfirst-data -repo whosonfirst-data-admin-us

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/sfomuseum/iso8601duration"
	"github.com/whosonfirst/go-whosonfirst-github/repositories"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to query")
	repo := flag.String("repo", "", "The name of the repository to query")

	updated_since := flag.String("updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")
	token := flag.String("token", "", "A valid GitHub API access token")

	ensure_geojson := flag.Bool("ensure-geojson", true, "Ensure that commits are for files ending in .geojson")
	max_commits := flag.Int("max-commits", 1, "...")
	flag.Parse()

	ctx := context.Background()

	opts := &repositories.ListCommitFilesOptions{
		AccessToken: *token,
		Org:         *org,
		Repo:        *repo,
		MaxCommits:  *max_commits,
	}

	if *ensure_geojson {

		re, err := regexp.Compile(`.*\.geojson$`)

		if err != nil {
			log.Fatalf("Failed to compile ensure geojson pattern, %v", err)
		}

		opts.Matching = re
	}

	if *updated_since != "" {

		var since time.Time

		is_timestamp, err := regexp.MatchString("^\\d+$", *updated_since)

		if err != nil {
			log.Fatal(err)
		}

		if is_timestamp {

			ts, err := strconv.Atoi(*updated_since)

			if err != nil {
				log.Fatal(err)
			}

			now := time.Now()

			tm := time.Unix(int64(ts), 0)
			since = now.Add(-time.Since(tm))

		} else {

			// maybe also this https://github.com/araddon/dateparse ?

			d, err := duration.FromString(*updated_since)

			if err != nil {
				log.Fatal(err)
			}

			now := time.Now()
			since = now.Add(-d.ToDuration())
		}

		// log.Printf("SINCE %v\n", since)
		// os.Exit(0)

		opts.Since = &since
	}

	files, err := repositories.ListCommitFiles(ctx, opts)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f)
	}

	os.Exit(0)
}
