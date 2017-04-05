package main

// https://godoc.org/github.com/google/go-github/github#Client
// https://developer.github.com/v3/orgs/hooks/#list-hooks

import (
        "context"
	"flag"
	"github.com/google/go-github/github"
	"log"
	"os"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	// prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	// forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	// not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")

	flag.Parse()

	client := github.NewClient(nil)

	ctx := context.TODO()

	opts := github.ListOptions{PerPage: 100}

	hooks, rsp, err := client.Organizations.ListHooks(ctx, *org, &opts)

	if err != nil {
	   log.Fatal(err)
	}
	
	log.Println(rsp)
	
	for _, h := range hooks {

	    log.Println(h)
	}

	os.Exit(0)
}
