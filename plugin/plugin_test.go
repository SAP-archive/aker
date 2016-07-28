package plugin_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"

	. "github.infra.hana.ondemand.com/I061150/aker/plugin"

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
		plugin, err = Open("./"+pluginName, config, nil)
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

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(plugin).ShouldNot(BeNil())
		})

		It("should have received the correct configuration", func() {
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "http://does.not.matter.com", nil)
			plugin.ServeHTTP(rr, req)
			Ω(rr.Body.Bytes()).Should(Equal(config))
		})

	})

})
