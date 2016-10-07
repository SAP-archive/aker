package plugin_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"time"

	. "github.com/SAP/aker/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Open Plugin", func() {
	const buildedPluginName = "aker-test-plugin"

	var pluginName string
	var config []byte

	var plugin *Plugin
	var err error

	BeforeEach(func() {
		buildCmd := exec.Command("go", "build", "-o", buildedPluginName, "test_plugin/main.go")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		err := buildCmd.Run()
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		os.Remove("./" + pluginName)
	})

	JustBeforeEach(func() {
		opener := &Opener{
			PluginStdout: GinkgoWriter,
			PluginStderr: GinkgoWriter,
		}
		plugin, err = opener.Open("./"+pluginName, config, nil)
	})

	Context("when the plugin does not exist", func() {
		BeforeEach(func() {
			pluginName = "not_existing"
		})

		It("should return an error", func() {
			Ω(err).Should(HaveOccurred())
			Ω(plugin).Should(BeNil())
		})
	})

	Context("when the plugin exists", func() {
		BeforeEach(func() {
			pluginName = buildedPluginName
			config = []byte("lalalalala")
		})

		AfterEach(func() {
			// signal plugin's process to exit
			Ω(plugin.Close()).Should(Succeed())
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(plugin).ShouldNot(BeNil())
			// give the process chance to start
			time.Sleep(time.Millisecond * 20)
		})

		It("should receive the correct configuration", func() {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://does.not.matter.com", nil)
			Ω(err).ShouldNot(HaveOccurred())
			plugin.ServeHTTP(rr, req)
			Ω(rr.Body.Bytes()).Should(Equal(config))
		})

	})

})
