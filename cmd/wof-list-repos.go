package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
	"os"
	"strings"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")
	
	flag.Parse()

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, *org, opt)

		if err != nil {
			log.Fatal(err)
		}

		for _, r := range repos {

			if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
				continue
			}

			if *exclude != "" && strings.HasPrefix(*r.Name, *exclude) {
			   	continue
			}

			if *forked && !*r.Fork {
				continue
			}

			if *not_forked && *r.Fork {
				continue
			}

			fmt.Println(*r.Name)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	os.Exit(0)
}
