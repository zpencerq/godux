package godux_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGodux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Godux Suite")
}
