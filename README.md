# go-whosonfirst-github

Go package for working with Who's On First GitHub repositories.

## Important

This is still wet paint.

## Utilities

### wof-clone-repos

Clone Who's On First data repositories in parallel.

```
./bin/wof-clone-repos -h
Usage of ./bin/wof-clone-repos:
  -destination string
        Where to clone repositories to (default "/usr/local/data")
  -org string
        The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
        Limit repositories to only those with this prefix (default "whosonfirst-data")
  -procs int
        The number of concurrent processes to clone with (default 20)
```

## See also

* https://github.com/whosonfirst-data/
* https://github.com/google/go-github
