package repositories

import (
	"context"
	"fmt"
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
}

func ListCommitFiles(ctx context.Context, opts *ListCommitFilesOptions) ([]string, error) {

	files := make([]string, 0)

	cb := func(ctx context.Context, f *github.CommitFile) error {
		files = append(files, *f.Filename)
		return nil
	}

	err := ListCommitFilesWithCallback(ctx, opts, cb)

	if err != nil {
		return nil, err
	}

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

	for _, c := range commits {

		for _, f := range c.Files {

			err := cb(ctx, f)

			if err != nil {
				return fmt.Errorf("Failed to execute callback for %s", f.Filename)
			}
		}
	}

	return nil
}
