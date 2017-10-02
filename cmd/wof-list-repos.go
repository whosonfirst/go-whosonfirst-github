package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"log"
	"os"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	flag.Parse()

	opts := organizations.NewDefaultListOptions()

	opts.Prefix = *prefix
	opts.Exclude = *exclude
	opts.Forked = *forked
	opts.NotForked = *not_forked
	opts.AccessToken = *token
	
	repos, err := organizations.ListRepos(*org, opts)

	if err != nil {
	   log.Fatal(err)
	}

	for _, name := range repos {
		fmt.Println(name)
	}

	os.Exit(0)
}
