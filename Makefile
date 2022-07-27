setup:
	git config core.hooksPath ./hooks

t: test
test: setup
	go test -count=1 ./...
