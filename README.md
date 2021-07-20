# go-whosonfirst-github

Go package for working with Who's On First GitHub repositories.

## Important

Here's a weird thing. If you try to use this code, either via the CLI tools or in your own packages, using `go run ... ` it will fail with incomprehensible "bad file descriptor" errors whenever the `google/go-github` code tries to make an HTTP request. If you compile things and run the same code as a binary everything works fine. Computers, amirite?

## Tools

### Webhooks

_All of the webhook tools need some documentation loving..._

#### wof-create-hook

```
> ./bin/wof-create-hook -h
Usage of ./bin/wof-create-hook:
  -dryrun
    	Go through the motions but don't create webhooks.
  -exclude value
    	Exclude repositories with this prefix
  -hook-content-type string
    	The content type for your webhook. (default "json")
  -hook-secret string
    	The secret for your webhook.
  -hook-url string
    	A valid webhook URL.
  -org string
    	The GitHub organization to create webhookd in.
  -prefix value
    	Limit repositories to only those with this prefix
  -token string
    	A valid GitHub API access token.
```

For example:

```
$> ./bin/wof-create-hook -org sfomuseum-data -hook-secret {SECRET} -hook-url {WEBHOOK_URL} -prefix sfomuseum-data -exclude sfomuseum-data-flights -exclude sfomuseum-data-faa -token {GITHUB_TOKEN}
2020/09/08 16:46:17 created webhook for sfomuseum-data-whosonfirst
2020/09/08 16:46:18 created webhook for sfomuseum-data-publicart
2020/09/08 16:46:19 created webhook for sfomuseum-data-architecture
2020/09/08 16:46:20 created webhook for sfomuseum-data-maps
2020/09/08 16:46:21 created webhook for sfomuseum-data-exhibition
2020/09/08 16:46:22 created webhook for sfomuseum-data-enterprise
2020/09/08 16:46:23 created webhook for sfomuseum-data-aircraft
2020/09/08 16:46:24 created webhook for sfomuseum-data-media
2020/09/08 16:46:25 created webhook for sfomuseum-data-media-flickr
2020/09/08 16:46:26 created webhook for sfomuseum-data-testing
2020/09/08 16:46:28 created webhook for sfomuseum-data-collection-classifications
2020/09/08 16:46:29 created webhook for sfomuseum-data-media-collection
2020/09/08 16:46:29 created webhook for sfomuseum-data-collection
```

#### wof-list-hooks

_Please write me_

#### wof-update-hook

_Please write me_

```
./bin/wof-update-hook -token {TOKEN} -hook-url {URL} -hook-secret {NEW_SECRET} -org whosonfirst-data -repo whosonfirst-data-venue-us-il
```

You can also update webhooks for all of the repositories in an organization by passing the `-repo '*'` flag. You can still filter the list of repos by setting the `-prefix` flag.

```
./bin/wof-update-hook -token {TOKEN} -hook-url {URL} -hook-secret {NEW_SECRET} -org whosonfirst-data -repo '*' -prefix whosonfirst-data
fetching repo list...🕓 
2017/04/05 15:42:24 edited webhook for whosonfirst-data
2017/04/05 15:42:24 edited webhook for whosonfirst-data-venue-us-wv
2017/04/05 15:42:25 edited webhook for whosonfirst-data-venue-us-ne
2017/04/05 15:42:25 edited webhook for whosonfirst-data-venue-us-wi

...and so on
```

### Repos

#### wof-clone-repos

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

#### wof-list-repos

Print (to STDOUT) the list of repository names for an organization.

```
./bin/wof-list-repos -h
Usage of ./bin/wof-list-repos:
  -exclude string
    	Exclude repositories with this prefix
  -forked
    	Only include repositories that have been forked
  -not-forked
    	Only include repositories that have not been forked
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix string
    	Limit repositories to only those with this prefix (default "whosonfirst-data")
  -token string
    	A valid GitHub API access token
  -updated-since string
    	A valid ISO8601 duration string (months are currently not supported)
```

For example:

```
./bin/wof-list-repos -org whosonfirst -prefix '' -forked | sort
Clustr
emoji-search
flamework
flamework-api
flamework-artisanal-integers
flamework-aws
flamework-geo
flamework-invitecodes
flamework-multifactor-auth
flamework-storage
flamework-tools
go-pubsocketd
go-ucd
now
privatesquare
py-flamework-api
py-machinetag
py-slack-api
python-edtf
reachable
redis-tools
slackcat
suncalc-go
walk
watchman
whereonearth-metropolitan-area
youarehere-www
```

## Caveats

### Things this package doesn't deal with (yet)

* Anything that requires a GitHub API access token
* Anything other than the `master` branch of a repository
* The ability to exclude specific repositories

## See also

* https://github.com/whosonfirst-data/
* https://github.com/google/go-github
* https://en.wikipedia.org/wiki/ISO_8601#Durations