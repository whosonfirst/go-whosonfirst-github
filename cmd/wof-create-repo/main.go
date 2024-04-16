package main

import (
	"flag"
	"log"

	"github.com/whosonfirst/go-whosonfirst-github/organizations"
)

func main() {

	org := flag.String("org", "", "The GitHub organization to create webhookd in.")
	token := flag.String("token", "", "A valid GitHub API access token.")

	name := flag.String("name", "", "The name of the repo to create")
	description := flag.String("description", "", "The description of the repo to create")
	private := flag.Bool("private", true, "A flag to indicate whether the repo is public or private.")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing OAuth2 token")
	}

	opts := &organizations.CreateOptions{
		AccessToken: *token,
		Name:        *name,
		Description: *description,
		Private:     *private,
	}

	err := organizations.CreateRepo(*org, opts)

	if err != nil {
		log.Fatalf("Failed to create new repo, %v", err)
	}
}
