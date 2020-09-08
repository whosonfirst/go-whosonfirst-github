package main

// https://developer.github.com/v3/repos/hooks/#create-a-hook

// https://godoc.org/github.com/google/go-github/github#RepositoriesService.CreateHook
// https://godoc.org/github.com/google/go-github/github#OrganizationsService.CreateHook
// https://godoc.org/github.com/google/go-github/github#Hook

import (
	"flag"
	"fmt"
	"github.com/google/go-github/v27/github"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
)

func main() {

	org := flag.String("org", "", "The GitHub organization to create webhookd in.")
	token := flag.String("token", "", "A valid GitHub API access token.")

	url := flag.String("hook-url", "", "A valid webhook URL.")
	content_type := flag.String("hook-content-type", "json", "The content type for your webhook.")
	secret := flag.String("hook-secret", "", "The secret for your webhook.")

	var prefix multi.MultiString
	flag.Var(&prefix, "prefix", "Limit repositories to only those with this prefix")

	var exclude multi.MultiString
	flag.Var(&exclude, "exclude", "Exclude repositories with this prefix")

	dryrun := flag.Bool("dryrun", false, "Go through the motions but don't create webhooks.")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing OAuth2 token")
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	hook_config := make(map[string]interface{})

	hook_config["url"] = *url
	hook_config["content_type"] = *content_type
	hook_config["secret"] = *secret

	hook := github.Hook{
		// Name:   name,
		Config: hook_config,
	}

	opts := organizations.NewDefaultListOptions()

	opts.Prefix = prefix
	opts.Exclude = exclude
	// opts.Forked = *forked
	// opts.NotForked = *not_forked
	opts.AccessToken = *token

	repos, err := organizations.ListRepos(*org, opts)

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {

		has_hook := false

		hooks_opts := github.ListOptions{PerPage: 100}

		hooks, _, err := client.Repositories.ListHooks(ctx, *org, r, &hooks_opts)

		if err != nil {
			log.Fatal(err)
		}

		for _, h := range hooks {

			if h.Config["url"] == *url {
				has_hook = true
				break
			}
		}

		if has_hook {
			log.Println(fmt.Sprintf("webhook already configured for %s, skipping", r))
			continue
		}

		if *dryrun {
			log.Printf("Create Webhook for %s\n", r)
			continue
		}

		_, _, err = client.Repositories.CreateHook(ctx, *org, r, &hook)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(fmt.Sprintf("created webhook for %s", r))
	}

}
