# Releasing

A release is made from the GitHub UI: publish a release with a new `vX.Y.Z` tag, and the **Release** workflow builds that tag as version `X.Y.Z` and attaches the Linux binaries. No version-bump commit or pull request is needed first.

## Versions

The version lives only in the tag. `go.mod` carries no release version and the binary does not report one, so the files at a release tag look the same as on `main`.

Earlier releases used `vX.Y.Z` tags too (the 3.x line used `vX.Y.Z.RELEASE`). The workflow accepts plain `vX.Y.Z` only.

## Making a release

1. GitHub -> **Releases** -> **Draft a new release**.
2. **Choose a tag**: type a new tag `vX.Y.Z` (e.g. `v4.3.0`) and select **Create new tag on publish**, with target `main`.
3. **Generate release notes**, edit as needed.
4. **Publish release**.

The **Release** workflow then:

1. checks the tag has the form `vX.Y.Z`;
2. checks out the tag and runs `go vet ./...` and `go test -count=1 ./...`;
3. builds static binaries (`CGO_ENABLED=0`, `-trimpath`) for linux/amd64 and linux/arm64 with `.github/scripts/package.sh`;
4. attaches the archives and a `SHA256SUMS` file to the release.

| Asset | Contents |
| ----- | -------- |
| `memory-calculator-X.Y.Z.tgz` | linux/amd64 binary. Keeps the name the Java buildpack manifest has used since 4.1.0. |
| `memory-calculator-X.Y.Z-linux-arm64.tgz` | linux/arm64 binary. |
| `SHA256SUMS` | SHA-256 of both archives, for the buildpack manifest `sha256` field. |

Each archive holds one file, `java-buildpack-memory-calculator`. The archives are reproducible: building the same commit with the same Go version gives the same bytes.

Pre-releases are ignored by the workflow, because the rolling snapshot is a pre-release.

### When the workflow fails

Nothing is attached. Fix the cause, then either:

- delete the release **and** its tag (Releases -> the release -> Delete; Tags -> the tag -> Delete) and publish again, or
- if the tag is right and only the build failed for an unrelated reason, re-run it: **Actions -> Release -> Run workflow** with the tag (e.g. `v4.3.0`). This rebuilds the tag and replaces the attached artifacts.

## Snapshots

Every push to `main` that passes the **CI** workflow recreates the pre-release `snapshot` with the same two archives, named `memory-calculator-snapshot-<short sha>.tgz` and `...-linux-arm64.tgz`. It is for testing and not a version to depend on.

## Using a release in the Java buildpack

Set `uri` in the buildpack `manifest.yml` to the release asset (`https://github.com/cloudfoundry/java-buildpack-memory-calculator/releases/download/vX.Y.Z/memory-calculator-X.Y.Z.tgz`) and `sha256` to the matching line in `SHA256SUMS`.

## Dependency updates

Dependabot (`.github/dependabot.yml`) opens weekly pull requests for the Go modules and for the GitHub Actions used by the workflows. The **Dependency Review** workflow fails a pull request that adds a dependency with a known high or critical vulnerability.

## Previous process

The Concourse scripts in `ci/` are superseded by these workflows; see `ci/README.md`.
