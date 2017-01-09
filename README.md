# Java Buildpack Memory Calculator

### Getting started

[Install Go][] and then `get` the memory calculator (in the Go source tree).

We run our tests with [Ginkgo/Gomega][] and manage dependencies with [Godep][].
Ginkgo is one of the dependencies we manage, so get Godep before starting work:

```shell
go get -v github.com/cloudfoundry/java-buildpack-memory-calculator
cd src/github.com/cloudfoundry/java-buildpack-memory-calculator

go get -v github.com/tools/godep
```

(The `-v` options on `go get` are there so you can see what packages are compiled under the covers.)

The (bash) script `ci/test.sh` uses (the correct version of) Ginkgo to
run the tests (using the correct versions of the dependencies). `test.sh`
will recompile Ginkgo if necessary.

The parameters to `runTests` are passed directly to Ginkgo.  For example:

```shell
ci/test.sh -r=false memory
```

will run the tests in the memory subdirectory *without* recursion into lower
subdirectories (which is the default).

The current Go environment is not modified by `test.sh`.

### Development

To develop against the code, you should issue:

```shell
godep restore
```
in the project directory before building or running tests directly from the command line.

If you wish to develop against a particular tagged *version* then, in the
project directory, you need to checkout this version (using
`git checkout <tag>`) and re-issue `godep restore` before proceeding.

If `godep restore` fails, it is because one of the dependencies cannot be
obtained, or else it cannot be (re)set to the version this project depends on.
Normally `go get -u <project>` for the dependency in error will then allow
`godep restore` to complete normally.

### Release binaries

The executables are built for more than one platform, so the Go compiler must exist
for the target platforms we need (currently linux and darwin). The shell script (`ci/build.sh`)
will use the Go compiler with the `GOOS` environment variable to generate the executables.

This will not work if the Go installation doesn't support all these platforms, so you may have to
ensure Go is installed with cross-compiler support.

### Design

The document [Java Buildpack Memory Calculator v3](https://docs.google.com/document/d/1vlXBiwRIjwiVcbvUGYMrxx2Aw1RVAtxq3iuZ3UK2vXA/edit?usp=sharing)
provides some rationale for the memory calculator externals.

## License

The Spring Cloud Services CLI plugin is Open Source software released under the
[Apache 2.0 license][].

[Install Go]: http://golang.org/doc/install
[Godep]: http://github.com/tools/godep
[Ginkgo/Gomega]: http://github.com/onsi/ginkgo
[Apache 2.0 license]: http://www.apache.org/licenses/LICENSE-2.0.html