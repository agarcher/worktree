# Contributing to wt

Thanks for your interest in contributing to `wt`!

## Getting Started

### Prerequisites

- Go 1.21+
- Git 2.5+ (2.13+ recommended)
- Make

### Building

```bash
# Clone the repository
git clone https://github.com/agarcher/wt.git
cd wt

# Build for current platform
make build

# Run tests
make test

# Run linter
make lint
```

### Testing Locally

This repository includes a `.wt.yaml` config, so you can test the tool against itself:

```bash
# Build and add to PATH
make build
export PATH="$(pwd)/build:$PATH"

# Set up shell integration
eval "$(./build/wt init zsh)"  # or bash/fish

# Test commands
wt list
wt create test-feature
wt cd test-feature
wt exit
wt delete test-feature
```

## Architecture

For details on the codebase structure and design patterns, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Git Workflow

- **Never squash or rebase on merge.** Always use `gh pr merge --merge`.
- Create feature branches for your work
- Write clear commit messages

## Releasing

Version is tracked in the `VERSION` file. To release:

```bash
make release patch "Fix bug in cleanup"
make release minor "Add new feature"
make release major "Breaking change"
```

This bumps `VERSION`, commits, tags, and pushes. GitHub Actions then builds binaries and updates the Homebrew tap.

## Code Style

- Run `make lint` before submitting PRs
- Follow existing code patterns in the codebase
- Keep changes focused and minimal
