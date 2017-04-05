package main

// https://developer.github.com/v3/repos/hooks/#create-a-hook

// https://godoc.org/github.com/google/go-github/github#RepositoriesService.CreateHook
// https://godoc.org/github.com/google/go-github/github#OrganizationsService.CreateHook
// https://godoc.org/github.com/google/go-github/github#Hook

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
	"strings"
)

func main() {

	org := flag.String("org", "whosonfirst", "")
	repo := flag.String("repo", "", "")
	token := flag.String("oauth2-token", "", "...")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")

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

		has_hook := make(map[string]bool)

		repos := make([]string, 0)

		if *repo == "all" {

			repos_opts := &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{PerPage: 100},
			}

			for {

				repos_list, repos_rsp, err := client.Repositories.ListByOrg(ctx, *org, repos_opts)

				if err != nil {
					log.Fatal(err)
				}

				for _, r := range repos_list {

					log.Println(r)

					if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
						continue
					}

					repos = append(repos, *r.Name)

					hooks_opts := github.ListOptions{PerPage: 100}

					hooks, _, err := client.Repositories.ListHooks(ctx, *org, *r.Name, &hooks_opts)

					if err != nil {
						log.Fatal(err)
					}

					for _, h := range hooks {

						if h.Config["url"] == *url {
							has_hook[*r.Name] = true
							break
						}
					}

				}

				if repos_rsp.NextPage == 0 {
					break
				}

				repos_opts.ListOptions.Page = repos_rsp.NextPage
			}

		} else {
			repos = append(repos, *repo)
		}

		for _, r := range repos {

			_, ok := has_hook[r]

			if ok {
				log.Println("hook already configured", r)
				continue
			}

			log.Println(fmt.Sprintf("CREATE HOOK FOR '%s'", r))
			// _, _, err = client.Repositories.CreateHook(ctx, *org, r, &hook)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
}
