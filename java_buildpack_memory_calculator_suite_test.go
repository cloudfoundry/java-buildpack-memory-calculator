package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJavaBuildpackMemoryCalculator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JavaBuildpackMemoryCalculator Suite")
}
