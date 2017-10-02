package organizations

import (
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"strings"
)

type ListOptions struct {
     Prefix string
     Exclude string
     Forked bool
     NotForked bool
     AccessToken string
}

func NewDefaultListOptions() *ListOptions {

     opts := ListOptions{
     Prefix: "",
     Exclude: "",
     Forked: false,
     NotForked: false,
     AccessToken: "",
     }

     return &opts
}

func ListRepos(org string, opts *ListOptions) ([]string, error) {

     	repos := make([]string, 0)
	
	client, ctx, err := util.NewClientAndContext(opts.AccessToken)

	if err != nil {
		return nil, err
	}

	gh_opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
	
		possible, resp, err := client.Repositories.ListByOrg(ctx, org, gh_opts)

		if err != nil {
		   	return nil, err
		}

		for _, r := range possible {

			if opts.Prefix != "" && !strings.HasPrefix(*r.Name, opts.Prefix) {
				continue
			}

			if opts.Exclude != "" && strings.HasPrefix(*r.Name, opts.Exclude) {
				continue
			}

			if opts.Forked && !*r.Fork {
				continue
			}

			if opts.NotForked && *r.Fork {
				continue
			}

			repos = append(repos, *r.Name)
		}

		if resp.NextPage == 0 {
			break
		}

		gh_opts.ListOptions.Page = resp.NextPage
	}

	return repos, nil
}
