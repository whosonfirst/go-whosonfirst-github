// wof-update-hook is a command line tool to update one or more GitHub webhooks. For example.
//
//	$> ./bin/wof-update-hook -org sfomuseum-data -prefix sfomuseum-data -token {TOKEN} -hook-url https://{HOST}/webhookd/dist -hook-active=false
package main

// https://developer.github.com/v3/repos/hooks/#edit-a-hook
// https://godoc.org/github.com/google/go-github/github#RepositoriesService.EditHook
// https://godoc.org/github.com/google/go-github/github#OrganizationsService.EditHook
// https://godoc.org/github.com/google/go-github/github#Hook

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"	
)

type update struct {
	Repo string
	Hook *github.Hook
}

func main() {

	org := flag.String("org", "", "")
	token := flag.String("token", "", "...")
	prefix := flag.String("prefix", "", "Limit repositories to only those with this prefix")

	url := flag.String("hook-url", "", "")
	content_type := flag.String("hook-content-type", "json", "")
	secret := flag.String("hook-secret", "", "")

	active := flag.Bool("hook-active", true, "")
	delete := flag.Bool("delete", false, "")

	dryrun := flag.Bool("dryrun", false, "")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing OAuth2 token")
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	updates := make([]update, 0)

	repos_opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {

		repos_list, repos_rsp, err := client.Repositories.ListByOrg(ctx, *org, repos_opts)

		if err != nil {
			log.Fatal(err)
		}

		for _, r := range repos_list {

			if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
				continue
			}

			hooks_opts := github.ListOptions{PerPage: 100}

			hooks, _, err := client.Repositories.ListHooks(ctx, *org, *r.Name, &hooks_opts)

			if err != nil {
				log.Fatal(err)
			}

			for _, h := range hooks {

				hook_url := h.Config["url"].(string)

				if strings.HasPrefix(hook_url, *url) {

					u := update{
						Repo: *r.Name,
						Hook: h,
					}

					updates = append(updates, u)
				}
			}

		}

		if repos_rsp.NextPage == 0 {
			break
		}

		repos_opts.ListOptions.Page = repos_rsp.NextPage
	}

	for _, u := range updates {

		r := u.Repo
		hook := u.Hook

		if *secret != "" {
			hook.Config["secret"] = *secret
		}

		if *content_type != "" {
			hook.Config["content_type"] = *content_type
		}

		hook.Active = active

		if *dryrun {
			log.Printf("DRYRUN %d (%s) %v\n", *hook.ID, r, hook)
			continue
		}

		if *delete {

			_, err := client.Repositories.DeleteHook(ctx, *org, r, *hook.ID)

			if err != nil {
				log.Fatal(fmt.Sprintf("failed to edit webhook for %s, because %s", r, err.Error()))
			}

		} else {

			_, _, err := client.Repositories.EditHook(ctx, *org, r, *hook.ID, hook)

			if err != nil {
				log.Fatal(fmt.Sprintf("failed to edit webhook for %s, because %s", r, err.Error()))
			}

		}

		log.Println(fmt.Sprintf("edited webhook for %s", r))
	}
}
