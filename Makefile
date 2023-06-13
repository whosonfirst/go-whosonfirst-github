GOMOD=vendor

tools:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-add-files cmd/wof-add-files/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-clone-repos cmd/wof-clone-repos/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-create-repos cmd/wof-create-repo/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-create-hook cmd/wof-create-hook/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-update-hook cmd/wof-update-hook/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-list-repos cmd/wof-list-repos/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-list-hooks cmd/wof-list-hooks/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-rate-limits cmd/wof-rate-limits/main.go
