package main

import (
	"flag"
	"github.com/google/go-github/github"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Clone(dest string, repo *github.Repository, throttle chan bool, wg *sync.WaitGroup) error {

	defer func() {
		wg.Done()
		throttle <- true
	}()

	log.Println("waiting to clone", *repo.Name)

	<-throttle

	name := *repo.Name

	remote := *repo.CloneURL
	local := filepath.Join(dest, name)

	_, err := os.Stat(local)

	if os.IsExist(err) {
		log.Println(local, "already exists")
		return err
	}

	git := "git"
	args := []string{"clone", remote, local}

	t1 := time.Now()

	cmd := exec.Command(git, args...)

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

			go Clone(*dest, r, throttle, wg)
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
