package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/sfomuseum/runtimevar"
	"github.com/whosonfirst/go-whosonfirst-github/util"
)

func main() {

	token_uri := flag.String("api-token-uri", "", "...")
	flag.Parse()

	ctx := context.Background()

	token, err := runtimevar.StringVar(ctx, *token_uri)

	if err != nil {
		log.Fatalf("Failed to expand token URI, %v", err)
	}

	client, ctx, err := util.NewClientAndContext(token)

	if err != nil {
		log.Fatalf("Failed to create new client, %v", err)
	}

	limits, _, err := client.RateLimits(ctx)

	if err != nil {
		log.Fatalf("Failed to derive rate limits, %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(limits)

	if err != nil {
		log.Fatalf("Failed to encode limits, %v", err)
	}
}
