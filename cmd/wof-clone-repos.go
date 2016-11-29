package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// https://godoc.org/github.com/google/go-github/github#Repository

func Clone(dest string, repo *github.Repository, giturl bool, throttle chan bool, wg *sync.WaitGroup) error {

	defer func() {
		wg.Done()
		throttle <- true
	}()

	log.Println("waiting to clone", *repo.Name)

	<-throttle

	name := *repo.Name

	remote := *repo.CloneURL

	if giturl {
		remote = *repo.GitURL
	}

	local := filepath.Join(dest, name)

	_, err := os.Stat(local)

	var git_args []string

	if os.IsNotExist(err) {

		git_args = []string{"clone", remote, local}

	} else {

		dot_git := filepath.Join(local, ".git")

		git_dir := fmt.Sprintf("--git-dir=%s", dot_git)
		work_dir := fmt.Sprintf("--work-dir=%s", dot_git)

		git_args = []string{git_dir, work_dir, "fetch"}

	}

	log.Println("git", strings.Join(git_args, " "))

	t1 := time.Now()

	cmd := exec.Command("git", git_args...)

	_, err = cmd.Output()

	if err != nil {
		log.Println("failed to clone", local, err)
		return err
	}

	t2 := time.Since(t1)

	log.Printf("time to clone %s, %v\n", local, t2)
	return nil
}

func main() {

	procs := flag.Int("procs", 20, "The number of concurrent processes to clone with")
	dest := flag.String("destination", "/usr/local/data", "Where to clone repositories to")
	org := flag.String("org", "whosonfirst-data", "The name of the organization to clone repositories from")
	prefix := flag.String("prefix", "whosonfirst-data", "Limit repositories to only those with this prefix")
	giturl := flag.Bool("giturl", false, "Clone using Git URL (rather than default HTTPS)")

	flag.Parse()

	info, err := os.Stat(*dest)

	if os.IsNotExist(err) {
		log.Fatal(*dest, "does not exist")
	}

	if !info.IsDir() {
		log.Fatal(*dest, "is not a directory")
	}

	client := github.NewClient(nil)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	throttle := make(chan bool, *procs)

	for i := 0; i < *procs; i++ {
		throttle <- true
	}

	dest_abs, err := filepath.Abs(*dest)

	if err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)

	for {
		repos, resp, err := client.Repositories.ListByOrg(*org, opt)

		if err != nil {
			log.Fatal(err)
		}

		for _, r := range repos {

			if *prefix != "" && !strings.HasPrefix(*r.Name, *prefix) {
				continue
			}

			wg.Add(1)

			go Clone(dest_abs, r, *giturl, throttle, wg)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	t1 := time.Now()

	wg.Wait()

	t2 := time.Since(t1)

	log.Printf("finished cloning all the repos in %v\n", t2)
	os.Exit(0)
}
