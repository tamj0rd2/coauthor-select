include .bingo/Variables.mk

.PHONY: test
.DEFAULT_GOAL := test

setup:
	git config core.hooksPath ./hooks
	git config --global url."git@github.com:saltpay".insteadOf https://github.com/saltpay
	go install github.com/bwplotka/bingo@latest
	bingo get

t: test
test: lint
	$(GOTEST) -count=1 ./...

lint:
	$(GOLANGCI_LINT) run --timeout=5m ./...

lf: lintfix
lintfix:
	@$(GOLANGCI_LINT) run ./... --fix

ci:
	git pull -r && make test && git push
