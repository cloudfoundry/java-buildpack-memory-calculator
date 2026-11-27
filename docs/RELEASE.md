# Releasing

A release is made from the GitHub UI: publish a release with a new `vX.Y.Z` tag, and the **Release** workflow builds that tag and attaches the Linux binaries. The version lives only in the tag, so no version-bump commit is needed.

## Making a release

1. Add a `## vX.Y.Z` section at the top of [RELEASE_NOTES.md](../RELEASE_NOTES.md) in a normal pull request and merge it. Write each paragraph and list item on one line, because GitHub keeps single line breaks.
2. GitHub -> **Releases** -> **Draft a new release**.
3. **Choose a tag**: type a new tag `vX.Y.Z` (e.g. `v4.3.0`) and select **Create new tag on publish**, with target `main`. The notes must be on `main` before this, because the workflow reads them from the tag.
4. Leave the description empty. The workflow replaces it with the notes section.
5. **Publish release**.

The workflow checks the tag format, extracts the notes (failing early if the section is missing or empty), runs `go vet` and `go test`, builds static binaries, attaches them, and sets the release text. Pre-releases are ignored.

| Asset | Contents |
| ----- | -------- |
| `memory-calculator-X.Y.Z-linux-amd64.tgz` | linux/amd64 binary |
| `memory-calculator-X.Y.Z-linux-arm64.tgz` | linux/arm64 binary |
| `SHA256SUMS` | SHA-256 of both archives, for the `sha256` field in the Java buildpack `manifest.yml` |

Each archive holds one file, `java-buildpack-memory-calculator`. The binary is static (`CGO_ENABLED=0`) and runs on cflinuxfs4 and cflinuxfs5, so there is one archive per CPU architecture and none per stack. Releases up to 4.1.0 used the name `memory-calculator-X.Y.Z.tgz`, without a platform. Tooling that picks up the assets, such as the buildpacks-ci dependency-build, must use the platform names.

## When the workflow fails

If the tag check, notes extraction, vet, tests or packaging fail, nothing is attached. If the upload itself fails, the release can hold some of the three assets and not others, because the upload is not transactional. Check that all three assets are there before using a release.

Fix the cause, then either:

- delete the release **and** its tag (Releases -> the release -> Delete; Tags -> the tag -> Delete) and publish again. This is also the fix for a missing notes section: merge it to `main` first, because the workflow reads the file from the tag, or
- if the tag is right and only the build or upload failed for an unrelated reason, re-run it: **Actions -> Release -> Run workflow** with the tag (e.g. `v4.3.0`). This rebuilds the tag, replaces all attached artifacts and sets the release text again.

The Concourse scripts in `ci/` are superseded, see `ci/README.md`.

## Snapshots

Every push to `main` that passes the **CI** workflow recreates the pre-release `snapshot` with the two archives, named `memory-calculator-snapshot-<short sha>-linux-amd64.tgz` and `...-linux-arm64.tgz`, and `--version` prints `snapshot-<short sha>`. It is a convenience for trying out the latest `main`, not a version to depend on.

## Dependency updates

Dependabot (`.github/dependabot.yml`) opens weekly pull requests for the Go modules and the GitHub Actions. The **Dependency Review** workflow fails a pull request that adds a dependency with a known high or critical vulnerability.
