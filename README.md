# Java Buildpack Memory Calculator
[![Build Status](https://travis-ci.org/cloudfoundry/java-buildpack-memory-calculator.svg)](https://travis-ci.org/cloudfoundry/java-buildpack-memory-calculator)

### Getting started

[Install Go][] and then `get` the memory calculator (in the Go source tree).

We run our tests with [Ginkgo/Gomega][] and manage dependencies with [Godep][].
Ginkgo is one of the dependencies we manage, so get Godep, too.

```shell
go get -v github.com/cloudfoundry/java-buildpack-memory-calculator
cd src/github.com/cloudfoundry/java-buildpack-memory-calculator

go get -v github.com/tools/godep
```
(The `-v` options on `go get` are there so you can see what packages are compiled under the covers.)

The (bash) script `scripts/runTests` uses (the correct version of) Ginkgo to
run the tests (using the correct versions of the dependencies). `runTests`
will recompile Ginkgo if necessary.

The parameters to runTests are passed directly to Ginkgo.  For example:

```shell
scripts/runTests -r=false memory
```

will run the tests in the memory subdirectory *without* recursion into lower
subdirectories (which is the default).

The current Go environment is not modified by `runTests`.

### Development

To develop against the code, you should issue:

```shell
godep restore
```
in the project directory before running tests or building from the command line.

If you wish to develop against a particular tagged *version* then, in the
project directory, you need to checkout this version (using 
`git checkout <tag>`) and re-issue `godep restore` before proceeding.

If `godep restore` fails, it is because one of the dependencies cannot be
obtained, or else it cannot be (re)set to the version this project depends on.
Normally `go get -u <project>` for the dependency in error will then allow
`godep restore` to complete normally.

[Install Go]: http://golang.org/doc/install
[Godep]: http://github.com/tools/godep
[Ginkgo/Gomega]: http://github.com/onsi/ginkgo