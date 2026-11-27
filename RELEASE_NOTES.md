# Release notes

The newest release comes first. Before tagging a release, add a `## vX.Y.Z` section at the top. The Release workflow copies that section into the GitHub release text. Write each paragraph and list item on one line, because GitHub keeps single line breaks. See [docs/RELEASE.md](docs/RELEASE.md).

## v4.3.0

- New experimental `--xms-equals-xmx` flag sets the initial heap (`-Xms`) equal to the calculated or explicitly configured maximum heap (`-Xmx`). It is off by default, an explicit `-Xms` takes precedence, and an `-Xms` larger than `-Xmx` is an error that names the excess.
- Memory sizes that overflow int64 (for example `-Xms8388608T`) are now rejected with `memory size overflows: ...`. Before, they wrapped silently. This applies to all size options, with or without the new flag.
- JVM options separated by tabs or newlines are now recognized, not only spaces.
- Errors are printed with a trailing newline, and a `%` in an error message is no longer mangled.
- New `--version` flag prints the version. Release binaries report their release version, for example `4.3.0`.
- Release archives are named by platform: `memory-calculator-4.3.0-linux-amd64.tgz` and `memory-calculator-4.3.0-linux-arm64.tgz`, with a `SHA256SUMS` file. The earlier unsuffixed `memory-calculator-X.Y.Z.tgz` is gone. The binary is statically linked and runs on cflinuxfs4 and cflinuxfs5.
- Installing with `go install github.com/cloudfoundry/java-buildpack-memory-calculator/v4@latest` is documented. Building requires Go 1.25 or later.

## v4.2.0

Tagged without GitHub release notes. See the [commits](https://github.com/cloudfoundry/java-buildpack-memory-calculator/compare/v4.1.0...v4.2.0).
