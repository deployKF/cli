# deployKF - cli

The CLI for [deployKF](https://github.com/deployKF/deployKF).

## Usage

The simplest usage of the `deploykf` CLI is to run the following command:

```bash
deploykf \
  --source-version=v0.1.0 \
  --values "./custom-values.yaml" \
  --output-dir "./GENERATOR_OUTPUT"
```

Which will generate deployKF manifests in the `./GENERATOR_OUTPUT` directory based on the `v0.1.0` source version and your `./values.yaml` file.

The `--source-version` flag must be a tag associated with a [deployKF release](https://github.com/deployKF/deployKF/releases).

## Development

Running `make build` will build the binary for your local platform and output it to `./bin/deploykf`.

Running `make install` will install the binary to `/usr/local/bin`.

Running `make lint` will run `golangci-lint` against the codebase.

Running `make lint-fix` will attempt to automatically fix linting errors.