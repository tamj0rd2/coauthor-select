# coauthor-select

These docs are WIP :)

- **coauthor-select select <args>** - allows you to select coauthors from a list. Add this to your prepare-commit-msg hook
- **coauthor-select validate <args>** - validates your coauthors to make sure 4 eyes are on your code. Add this to your commit-msg hook

## How to use these tools as git hooks (if you're using go modules):

1. `go install github.com/tamj0rd2/coauthor-select@latest`
2. Create a hooks folder `mkdir .hooks` in your project and enable it as the git hooks folder `git config core.hooksPath .hooks`
3. Copy /examples/.hooks to your repo and make all files executable
4. Create a `.coauthors` file like [this](./example/.coauthors)

## Specifying pairs via the command line

1. Commit as you usually do
2. You'll be prompted to select pairs from the list or be given the option to choose the people you were last pairing with. This is enabled by the prepare-commit-msg hook.
3. You'll be warned if you're trying to commit to the trunk without specifying a pair

## Configuration

### select command

Check [here](src/select/main) for defaults and the latest documentation

- `--authorsFile` - the path to your authors.json file
- `--pairsFile` - the path to your pairs.json file
- `--interactive` - set this to false if you're using a non-interactive console

### validate command

Check [here](src/validate/main) for defaults and the latest documentation

- `--trunkName` - the name of the trunk branch
- `--protectTrunk` - if true, people will not be able to commit to the trunk branch without specifying a pair
