# deployKF - Command Line Interface (CLI)

[![Check Commit](https://github.com/deployKF/cli/actions/workflows/check-commit.yml/badge.svg)](https://github.com/deployKF/cli/actions/workflows/check-commit.yml)

This repo contains the command line interface (CLI) for [deployKF](https://github.com/deployKF/deployKF).

## Install

You can install the `deploykf` CLI by downloading the appropriate binary for your OS from the [releases page](https://github.com/deployKF/cli/releases).
 - You may wish to rename the binary to `deploykf` (or `deploykf.exe` on Windows).
 - You may need to make the binary executable on Linux/macOS by running `chmod +x deploykf`.
 - You may also wish to move the binary to a directory that is in your `PATH` environment variable. 
    - On Linux/macOS you might run `sudo mv deploykf /usr/local/bin/`.
    - On Windows you might run `move deploykf.exe C:\Windows\System32\`.

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

