package socket_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const testServerName = "test-server"

func TestSocket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Socket Suite")
}

var _ = BeforeSuite(func() {
	buildCmd := exec.Command("go", "build", "-o", testServerName, "test_server/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err := buildCmd.Run()
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.Remove(testServerName)
	Ω(err).ShouldNot(HaveOccurred())
})
