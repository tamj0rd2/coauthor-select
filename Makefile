setup:
	git config core.hooksPath ./hooks

test: setup
	go test ./...
