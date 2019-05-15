package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {

	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	forked := flag.Bool("forked", false, "Only include repositories that have been forked")
	not_forked := flag.Bool("not-forked", false, "Only include repositories that have not been forked")
	token := flag.String("token", "", "A valid GitHub API access token")

	// mmmmmaybe? seems like overkill right now (20190515/thisisaaronland)
	// updated_since := flag.String("updated-since", "", "A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)")

	flag.Parse()

	if *token == "" {
		log.Fatal("Missing Github API access token")
	}

	files := make([]string, 0)

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(path, err)
		}

		_, err = os.Stat(abs_path)

		if err != nil {
			log.Fatal(abs_path, err)
		}

		files = append(files, abs_path)
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

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

	for _, repo := range repos {

		err := add_files(ctx, client, *org, repo, files...)

		if err != nil {
			log.Fatal(org, repo, err)
		}
	}

}

func add_files(ctx context.Context, client *github.Client, org string, repo string, files ...string) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	// something something something
	// channels and/or failing hard on any single error condition
	// or something something something, right?
	// (20190515/thisisaaronland)

	wg := new(sync.WaitGroup)

	for _, path := range files {

		wg.Add(1)

		go func(path string) {

			defer wg.Done()
			err := add_file(ctx, client, org, repo, path)

			if err != nil {
				log.Println("ERROR", org, repo, path, err)
			}

		}(path)
	}

	wg.Wait()

	return nil
}

// https://godoc.org/github.com/google/go-github/github#ex-RepositoriesService-CreateFile

func add_file(ctx context.Context, client *github.Client, org string, repo string, abs_path string) error {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return nil
	}

	fname := filepath.Base(abs_path)

	commit_name := fmt.Sprintf("%s bot", org)
	commit_email := fmt.Sprintf("%s@localhost", org)

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("initial commit"),
		Content:   body,
		Branch:    github.String("master"), // PLEASE MAKE ME AN OPTION...
		Committer: &github.CommitAuthor{Name: github.String(commit_name), Email: github.String(commit_email)},
	}

	log.Println("ADD", abs_path, org, repo, commit_name)

	_, _, err = client.Repositories.CreateFile(ctx, org, repo, fname, opts)

	if err != nil {
		return err
	}

	return nil
}
