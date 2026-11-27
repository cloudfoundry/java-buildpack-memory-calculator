/*
 * Copyright 2026 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

var binary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "memory-calculator")
	if err != nil {
		panic(err)
	}

	binary = filepath.Join(dir, "java-buildpack-memory-calculator")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		_ = os.RemoveAll(dir)
		panic(string(out))
	}

	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func run(args ...string) (string, string, int) {
	cmd := exec.Command(binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return stdout.String(), stderr.String(), 0
	}

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		panic(err)
	}

	return stdout.String(), stderr.String(), exitError.ExitCode()
}

func buildWithVersion(t *testing.T, version string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "java-buildpack-memory-calculator")
	out, err := exec.Command("go", "build", "-ldflags", "-X main.version="+version, "-o", path, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build: %s", out)
	}

	return path
}

func TestVersion(t *testing.T) {
	t.Run("prints the stamped version without other flags", func(t *testing.T) {
		g := NewGomegaWithT(t)

		out, err := exec.Command(buildWithVersion(t, "1.2.3"), "--version").Output()

		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(string(out)).To(Equal("1.2.3\n"))
	})

	t.Run("prints a fallback when no version is stamped", func(t *testing.T) {
		g := NewGomegaWithT(t)

		stdout, stderr, code := run("--version")

		g.Expect(code).To(Equal(0))
		g.Expect(stderr).To(BeEmpty())
		g.Expect(stdout).To(MatchRegexp(`^\S+\n$`))
	})

	t.Run("is listed in the usage", func(t *testing.T) {
		g := NewGomegaWithT(t)

		_, stderr, _ := run("--help")

		g.Expect(stderr).To(ContainSubstring("--version"))
	})
}

func TestCLI(t *testing.T) {
	spec.Run(t, "CLI", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		base := []string{"--total-memory=1G", "--loaded-class-count=10000", "--thread-count=250"}

		it("does not emit an initial heap by default", func() {
			stdout, _, code := run(base...)

			g.Expect(code).To(Equal(0))
			g.Expect(stdout).To(Equal("-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=70312K -XX:ReservedCodeCacheSize=240M -Xss1M -Xmx466263K\n"))
		})

		it("emits an initial heap equal to the calculated maximum heap", func() {
			stdout, _, code := run(append(base, "--xms-equals-xmx")...)

			g.Expect(code).To(Equal(0))
			g.Expect(stdout).To(Equal("-XX:MaxDirectMemorySize=10M -XX:MaxMetaspaceSize=70312K -XX:ReservedCodeCacheSize=240M -Xss1M -Xmx466263K -Xms466263K\n"))
		})

		it("emits the exact bytes for an unaligned configured maximum heap", func() {
			stdout, _, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Xmx231858240")...)

			g.Expect(code).To(Equal(0))
			g.Expect(stdout).To(HaveSuffix(" -Xms231858240\n"))
			g.Expect(stdout).NotTo(ContainSubstring("-Xmx"))
		})

		it("keeps an explicit initial heap separated by a tab", func() {
			stdout, _, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Xms64M\t-Dfoo=bar")...)

			g.Expect(code).To(Equal(0))
			g.Expect(stdout).NotTo(ContainSubstring("-Xms"))
		})

		it("keeps an explicit initial heap separated by a newline", func() {
			stdout, _, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Dfoo=bar\n-Xms64M")...)

			g.Expect(code).To(Equal(0))
			g.Expect(stdout).NotTo(ContainSubstring("-Xms"))
		})

		it("rejects an initial heap larger than a configured maximum heap with a readable error", func() {
			stdout, stderr, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Xmx128M -Xms256M")...)

			g.Expect(code).To(Equal(2))
			g.Expect(stdout).To(BeEmpty())
			g.Expect(stderr).To(Equal("initial heap -Xms256M is greater than configured maximum heap -Xmx128M by 128M\n"))
		})

		it("rejects an initial heap larger than the calculated maximum heap with a readable error", func() {
			_, stderr, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Xms1G")...)

			g.Expect(code).To(Equal(2))
			g.Expect(stderr).To(Equal("initial heap -Xms1G is greater than calculated maximum heap -Xmx466263K by about 569M\n"))
		})

		it("ends an allocation error with a newline", func() {
			_, stderr, code := run("--total-memory=10M", "--loaded-class-count=10000", "--thread-count=250")

			g.Expect(code).To(Equal(2))
			g.Expect(stderr).To(HaveSuffix("\n"))
		})

		it("rejects a size that overflows", func() {
			_, stderr, code := run(append(base, "--xms-equals-xmx", "--jvm-options=-Xmx128M -Xms18014398509481985K")...)

			g.Expect(code).NotTo(Equal(0))
			g.Expect(stderr).To(ContainSubstring("memory size overflows: 18014398509481985K"))
		})

		it("rejects an overflowing size without the flag", func() {
			_, stderr, code := run(append(base, "--jvm-options=-Xmx8388608T")...)

			g.Expect(code).NotTo(Equal(0))
			g.Expect(stderr).To(ContainSubstring("memory size overflows: 8388608T"))
		})
	})
}
