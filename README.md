# deployKF - Command Line Interface (CLI)

<a href="https://github.com/deployKF/cli/actions/workflows/check-commit.yml">
  <img alt="Check Commit" src="https://github.com/deployKF/cli/actions/workflows/check-commit.yml/badge.svg">
</a>
<a href="https://github.com/deployKF/cli/releases">
  <img alt="Downloads" src="https://img.shields.io/github/downloads/deployKF/cli/total?color=28a745">
</a>
<a href="https://github.com/deployKF/cli/releases">
  <img alt="Latest Release" src="https://img.shields.io/github/v/release/deployKF/cli?color=6f42c1&label=latest%20release">
</a>

This repo contains the command line interface (CLI) for [deployKF](https://github.com/deployKF/deployKF).

## Install

To install the `deploykf` CLI, please follow the [installation instructions](https://www.deploykf.org/guides/install-deploykf-cli/) on the deployKF website.

## Usage

The simplest usage of the `deploykf` CLI is to run the following command:

```bash
deploykf \
  --source-version 0.1.1 \
  --values ./custom-values.yaml \
  --output-dir ./GENERATOR_OUTPUT
```

This command will generate deployKF manifests in the `./GENERATOR_OUTPUT` directory using the `v0.1.1` source version and the values specified in your `./custom-values.yaml` file. 
Note that the `--source-version` flag must correspond to a tag from a [deployKF release](https://github.com/deployKF/deployKF/releases).

> __TIP:__ 
> 
> The version of the CLI does NOT need to match the `--source-version` you are generating manifests for.
> If a breaking change is ever needed, the CLI will fail to generate with newer source versions, and will print message telling you to upgrade the CLI.

## Container Image

We publish the `deploykf` CLI as a container image on the following registries:

| Registry                  | Image                                                  | Pull Command                                |
|---------------------------|--------------------------------------------------------|---------------------------------------------|
| GitHub Container Registry | [`ghcr.io/deploykf/cli`](https://ghcr.io/deploykf/cli) | `docker pull ghcr.io/deploykf/cli:TAG_NAME` |

To use the container image, you need to mount your local filesystem into the container:

```bash
CONTAINER_IMAGE="ghcr.io/deploykf/cli:0.1.2"

docker run \
  --rm \
  --volume "$(pwd):/home/deploykf" \
  --volume "${HOME}/.deploykf:/home/deploykf/.deploykf" \
  "${CONTAINER_IMAGE}" \
    generate \
    --source-version "0.1.1" \
    --values ./sample-values.yaml \
    --output-dir ./GENERATOR_OUTPUT
```

## Development

Here are some helpful commands when developing the CLI:

- `make build`: Builds the binary for your local platform and outputs it to `./bin/deploykf`.
- `make install`: Installs the binary to `/usr/local/bin`.
- `make lint`: Runs `golangci-lint` against the codebase to check for errors.
- `make lint-fix`: Attempts to automatically fix any linting errors found.
