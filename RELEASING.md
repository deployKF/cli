# Releasing Guide

This guide is intended for maintainers who want to make a new release of the deployKF CLI.

1. For a new minor or major release, create a `release-*` branch first.
    - For example, for the `v0.2.0` release, create a new branch called `release-0.2`.
    - This allows for the continued release of bug fixes to older versions.
2. Create a new tag on the appropriate release branch for the version you are releasing.
    - For instance, you might create `v0.1.1` or `v0.1.1-alpha.1` on the `release-0.1` branch.
    - Ensure you ONLY create tags on the `release-*` branches, not on the `main` branch.
    - Remember to sign the tag with your GPG key.
       - You can do this by running `git tag -s v0.1.1 -m "v0.1.1"`.
       - You can verify the tag signature by running `git verify-tag v0.1.1`.
    - Ensure you ONLY push the specific tag you want to release. 
       - For example, if you want to release `v0.1.1`, you should run `git push origin v0.1.1`.
       - Do NOT run `git push origin --tags` or `git push origin main`.
    - When a new semver tag is created, a workflow will automatically create a GitHub draft release.
       - The release will include the binaries and corresponding SHA256 checksums for all supported platforms.
3. Generate the changelog using the "generate release notes" feature of GitHub and set it as the release description.
4. When ready to ship, manually publish the draft release.
    - This will trigger a workflow to build and push the Docker images for the release.
5. Update the changelog on the website using the [`update_changelogs.sh`](https://github.com/deployKF/website/blob/main/update_changelogs.sh) script.
