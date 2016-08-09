package plugin_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"

	"testing"
)

func TestPlugin(t *testing.T) {
	gologger.DefaultLogger = gologger.NewNativeLogger(GinkgoWriter, GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Suite")
}
