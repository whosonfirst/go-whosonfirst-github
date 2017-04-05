package main

// THIS HAS NOT BEEN TESTED YET
// https://developer.github.com/v3/repos/hooks/#create-a-hook

// https://godoc.org/github.com/google/go-github/github#RepositoriesService.CreateHook
// https://godoc.org/github.com/google/go-github/github#OrganizationsService.CreateHook
// https://godoc.org/github.com/google/go-github/github#Hook

import (
	"flag"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
)

func main() {

	org := flag.String("org", "whosonfirst", "")
	repo := flag.String("repo", "", "")
	token := flag.String("oauth2-token", "", "...")

	name := flag.String("hook-name", "web", "")
	url := flag.String("hook-url", "", "")
	content_type := flag.String("hook-content-type", "json", "")
	secret := flag.String("hook-secret", "", "")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing OAuth2 token")
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	// Check to see if *url is already registered for this org/repo...

	// https://developer.github.com/v3/repos/hooks/#example

	config := make(map[string]interface{})

	config["url"] = *url
	config["content_type"] = *content_type
	config["secret"] = *secret

	hook := github.Hook{
		Name:   name,
		Config: config,
	}

	if *repo == "" {

		_, _, err = client.Organizations.CreateHook(ctx, *org, &hook)

	} else {

		_, _, err = client.Repositories.CreateHook(ctx, *org, *repo, &hook)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
}
