package repositories

import (
	"context"
	"fmt"
	_ "log"
	"log/slog"
	"sync"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
)

type ListCommitFileCallback func(context.Context, *github.CommitFile) error

type ListCommitFilesOptions struct {
	AccessToken string
	Org         string
	Repo        string
	Since       *time.Time
	MaxCommits  int
}

func ListCommitFiles(ctx context.Context, opts *ListCommitFilesOptions) ([]string, error) {

	lookup := new(sync.Map)

	cb := func(ctx context.Context, f *github.CommitFile) error {
		lookup.Store(*f.Filename, true)
		return nil
	}

	err := ListCommitFilesWithCallback(ctx, opts, cb)

	if err != nil {
		return nil, err
	}

	files := make([]string, 0)

	lookup.Range(func(k interface{}, v interface{}) bool {
		f := k.(string)
		files = append(files, f)
		return true
	})

	return files, nil
}

func ListCommitFilesWithCallback(ctx context.Context, opts *ListCommitFilesOptions, cb ListCommitFileCallback) error {

	client, _, err := util.NewClientAndContext(opts.AccessToken)

	if err != nil {
		return fmt.Errorf("Failed to create new client, %w", err)
	}

	commits_opts := &github.CommitsListOptions{}

	if opts.Since != nil {
		commits_opts.Since = *opts.Since
	}

	commits, _, err := client.Repositories.ListCommits(ctx, opts.Org, opts.Repo, commits_opts)

	if err != nil {
		return fmt.Errorf("Failed to list commits for %s, %w", opts.Repo, err)
	}

	for i, rc := range commits {

		slog.Debug("Commit", "sha", *rc.SHA)

		// https://pkg.go.dev/github.com/google/go-github/v48/github#RepositoriesService.GetCommit
		// https://pkg.go.dev/github.com/google/go-github/v48/github#Response
		// https://pkg.go.dev/github.com/google/go-github/v48/github#ListOptions

		page := 1

		list_opts := &github.ListOptions{
			Page: page,
		}

		for {
			c, rsp, err := client.Repositories.GetCommit(ctx, opts.Org, opts.Repo, *rc.SHA, list_opts)

			if err != nil {
				return fmt.Errorf("Failed to get commit %s, %w", *rc.SHA, err)
			}

			for _, f := range c.Files {

				err := cb(ctx, f)

				if err != nil {
					return fmt.Errorf("Failed to execute callback for %s", f.Filename)
				}
			}

			slog.Error("Finished processing page", "current", page, "next", rsp.NextPage, "last", rsp.LastPage)

			if page < rsp.LastPage {
				page = rsp.NextPage
			} else {
				break
			}
		}

		if opts.MaxCommits > 0 && (i+1) >= opts.MaxCommits {
			break
		}
	}

	return nil
}
