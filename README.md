# Java Buildpack Memory Calculator
[![Build Status](https://travis-ci.org/cloudfoundry/java-buildpack-memory-calculator.svg)](https://travis-ci.org/cloudfoundry/java-buildpack-memory-calculator)

### Getting started

[Install Go][] and then get the memory calculator (in the Go source tree).
Then get [Godep][], restore the dependencies in the project directory, and you
should be able to run the tests successfully: 

```shell
go get github.com/cloudfoundry/java-buildpack-memory-calculator
cd src/github.com/cloudfoundry/java-buildpack-memory-calculator

go get github.com/tools/godep
godep restore

go test -a -v ./...
```

If you wish to work on a particular version then, in the project directory, you need to checkout this version (using `git checkout <tag>`) and repeat the `godep restore` before proceeding.

If `godep restore` fails, it is because one of the dependencies cannot be obtained, or else it cannot be (re)set to the version this project depends on.

[Install Go]: http://golang.org/doc/install
[Godep]: http://github.com/tools/godep