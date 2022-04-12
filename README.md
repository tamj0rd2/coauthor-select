# coauthor-select

## Usage

1. Install this package, somehow
2. Create a .hooks folder
3. Run `git config core.hooksPath .hooks`
4. In your project repo root, run `echo 'go run main.go' > .hooks/prepare-commit-msg && chmod +x ./hooks/prepare-commit-msg`
5. Add an authors.json to contain the names and email addresses of the people you commonly work with
6. Create a `pairs.json` file which you will need to maintain with the names of the people in your current pairing session
7. Commit. You'll be able to see the co-authors in the commit message

## Example

There's a working example in the /examples folder of this repo

## Contributing

- Run the tests using `make test`
