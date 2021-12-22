tools:
	go build -mod vendor -o bin/wof-add-files cmd/wof-add-files/main.go
	go build -mod vendor -o bin/wof-clone-repos cmd/wof-clone-repos/main.go
	go build -mod vendor -o bin/wof-create-repos cmd/wof-create-repo/main.go
	go build -mod vendor -o bin/wof-create-hook cmd/wof-create-hook/main.go
	go build -mod vendor -o bin/wof-update-hook cmd/wof-update-hook/main.go
	go build -mod vendor -o bin/wof-list-repos cmd/wof-list-repos/main.go
	go build -mod vendor -o bin/wof-list-hooks cmd/wof-list-hooks/main.go
