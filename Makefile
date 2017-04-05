CWD=$(shell pwd)
GOPATH := $(CWD)

build:	rmdeps deps fmt bin

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test ! -d src; then mkdir src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-github/util
	cp util/*.go src/github.com/whosonfirst/go-whosonfirst-github/util/
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/google/go-github/github"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/oauth2"
	@GOPATH=$(GOPATH) go get -u "github.com/briandowns/spinner"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt util/*.go

bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-clone-repos cmd/wof-clone-repos.go
	@GOPATH=$(GOPATH) go build -o bin/wof-create-hook cmd/wof-create-hook.go
	@GOPATH=$(GOPATH) go build -o bin/wof-list-repos cmd/wof-list-repos.go
	@GOPATH=$(GOPATH) go build -o bin/wof-list-hooks cmd/wof-list-hooks.go
