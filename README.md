# go-whosonfirst-github

Go package for working with Who's On First GitHub repositories.

## Utilities

### wof-clone-repos

Clone (or update from `master`) Who's On First data repositories in parallel.

```
./bin/wof-clone-repos -h
Usage of ./bin/wof-clone-repos:
  -destination string
    	Where to clone repositories to (default "/usr/local/data")
  -dryrun
    	Go through the motions but don't actually clone (or update) anything
  -giturl
    	Clone using Git URL (rather than default HTTPS)
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
    	Limit repositories to only those with this prefix (default "whosonfirst-data")
  -procs int
    	The number of concurrent processes to clone with (default 20)
```

### wof-list-repos

Print (to STDOUT) the list of repository names for an organization.

```
./bin/wof-list-repos -h
Usage of ./bin/wof-list-repos:
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
    	Limit repositories to only those with this prefix (default "whosonfirst-data")
```

## Caveats

### Things this package doesn't deal with (yet)

* Anything that requires a GitHub API access token
* Anything other than the `master` branch of a repository
* The ability to exclude specific repositories

## See also

* https://github.com/whosonfirst-data/
* https://github.com/google/go-github
