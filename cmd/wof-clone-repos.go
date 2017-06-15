package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/whosonfirst/go-whosonfirst-github/util"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// https://godoc.org/github.com/google/go-github/github#Repository

// please make me a struct-thingy or something (20161129/thisisaaronland)

func Clone(dest string, repo *github.Repository, giturl bool, throttle chan bool, wg *sync.WaitGroup, dryrun bool, strict bool) error {

	defer func() {
		wg.Done()
		throttle <- true
	}()

	<-throttle

	name := *repo.Name

	remote := *repo.CloneURL

	if giturl {
		remote = *repo.GitURL
		remote = strings.Replace(remote, "git://github.com/", "git@github.com:", -1) // why do I even need to do this...???
	}

	local := filepath.Join(dest, name)

	_, err := os.Stat(local)

	var git_args []string

	if os.IsNotExist(err) {

		git_args = []string{"clone", "-v", remote, local}

	} else {

		dot_git := filepath.Join(local, ".git")

		git_dir := fmt.Sprintf("--git-dir=%s", dot_git)
		work_tree := fmt.Sprintf("--work-tree=%s", dot_git)

		git_args = []string{git_dir, work_tree, "pull", "origin", "master"}
	}

	log.Println("git", strings.Join(git_args, " "))

	if dryrun {
		return nil
	}

	t1 := time.Now()

	cmd := exec.Command("git", git_args...)

	_, err = cmd.Output()

	if err != nil {

		if strict {

			msg := fmt.Sprintf("failed to clone %s because : %s", local, err)

			euid := os.Geteuid()
			str_euid := strconv.Itoa(euid)

			eff_user, eff_err := user.LookupId(str_euid)

			if eff_err == nil {
				msg = fmt.Sprintf("failed to clone %s because : %s (running as %s (%d))", local, err, eff_user.Username, euid)
			}

			log.Fatal(msg)
		}

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
	exclude := flag.String("exclude", "", "Exclude repositories with this prefix")
	giturl := flag.Bool("giturl", false, "Clone using Git URL (rather than default HTTPS)")
	dryrun := flag.Bool("dryrun", false, "Go through the motions but don't actually clone (or update) anything")
	strict := flag.Bool("strict", false, "If any attempt to clone a repo fails trigger a fatal error")
	token := flag.String("token", "", "A valid GitHub API access token")

	flag.Parse()

	info, err := os.Stat(*dest)

	if os.IsNotExist(err) {
		log.Fatal(*dest, "does not exist")
	}

	if !info.IsDir() {
		log.Fatal(*dest, "is not a directory")
	}

	client, ctx, err := util.NewClientAndContext(*token)

	if err != nil {
		log.Fatal(err)
	}

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

			wg.Add(1)

			go Clone(dest_abs, r, *giturl, throttle, wg, *dryrun, *strict)
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
