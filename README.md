# deployKF - Command Line Interface (CLI)

[![Check Commit](https://github.com/deployKF/cli/actions/workflows/check-commit.yml/badge.svg)](https://github.com/deployKF/cli/actions/workflows/check-commit.yml)

This repo contains the command line interface (CLI) for [deployKF](https://github.com/deployKF/deployKF).

## Install

You can install the `deploykf` CLI by downloading the appropriate binary for your OS from the [releases page](https://github.com/deployKF/cli/releases).
 - You may wish to rename the binary to `deploykf` (or `deploykf.exe` on Windows).
 - On Unix-like systems, you may need to make the binary executable by running `chmod +x deploykf`.

## Usage

The simplest usage of the `deploykf` CLI is to run the following command:

```bash
deploykf \
  --source-version v0.1.0 \
  --values ./custom-values.yaml \
  --output-dir ./GENERATOR_OUTPUT
```

This command will generate deployKF manifests in the `./GENERATOR_OUTPUT` directory using the `v0.1.0` source version and the values specified in your `./custom-values.yaml` file. 
Note that the `--source-version` flag must correspond to a tag from a [deployKF release](https://github.com/deployKF/deployKF/releases).

## Development

Here are some helpful commands when developing the CLI:

- `make build`: Builds the binary for your local platform and outputs it to `./bin/deploykf`.
- `make install`: Installs the binary to `/usr/local/bin`.
- `make lint`: Runs `golangci-lint` against the codebase to check for errors.
- `make lint-fix`: Attempts to automatically fix any linting errors found.

## Releasing

To release a new version of the CLI, follow these steps:

1. For a new minor or major release, create a `release-*` branch first.
    - For example, for the `v0.2.0` release, create a new branch called `release-0.2`. 
    - This allows for the continued release of bug fixes to older CLI versions.
2. Create a new tag on the appropriate release branch for the version you are releasing.
    - Ensure you sign the tag with your GPG key. 
       - You can do this by running `git tag -s v0.1.1 -m "v0.1.1"`.
       - You can verify the tag by running `git verify-tag v0.1.1`.
    - For instance, you might create `v0.1.1` or `v0.1.1-alpha.1` on the `release-0.1` branch.
    - Remember to create tags only on the `release-*` branches, not on the `main` branch.
3. When a new semver tag is created, a workflow will automatically create a GitHub draft release.
    - The release will include binaries and corresponding SHA256 checksums for all supported platforms.
    - Don't forget to add the changelog to release notes.
4. Manually publish the draft release.
