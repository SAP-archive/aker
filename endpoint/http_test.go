package endpoint_test

import (
	"errors"

	"github.infra.hana.ondemand.com/I061150/aker/config"
	. "github.infra.hana.ondemand.com/I061150/aker/endpoint"
	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
	"github.infra.hana.ondemand.com/I061150/aker/plugin/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {

	var endpoint config.Endpoint
	var opener *fakes.FakeOpener
	var handler *Handler
	var err error

	JustBeforeEach(func() {
		handler, err = NewHandler(endpoint, opener)
	})

	BeforeEach(func() {
		logging.DefaultLogger = logging.NewNativeLogger(GinkgoWriter, GinkgoWriter)
		opener = new(fakes.FakeOpener)
	})

	Context("when created with empty config path", func() {
		BeforeEach(func() {
			endpoint = config.Endpoint{}
		})

		It("should have returned an error", func() {
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(InvalidPathError("")))
		})

		It("should have returned nil handler", func() {
			Ω(handler).Should(BeNil())
		})
	})

	Context("when created with no plugin configuration", func() {
		BeforeEach(func() {
			endpoint = config.Endpoint{
				Path: "/",
			}
		})

		It("should have returned an error", func() {
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(NoPluginsErr))
		})

		It("should have returned nil handler", func() {
			Ω(handler).Should(BeNil())
		})
	})

	Context("when created with valid configuration", func() {
		BeforeEach(func() {
			endpoint.Path = "/"
			endpoint.Plugins = []config.PluginReference{
				config.PluginReference{
					Name: "happy-unicorn",
				},
				config.PluginReference{
					Name: "mighty-grasshopper",
					Config: config.PluginConfig{
						"fly": "no",
					},
				},
			}

			opener.OpenStub = func(name string, _ []byte, next *plugin.Plugin) (*plugin.Plugin, error) {
				return &plugin.Plugin{}, nil
			}
		})

		Context("but there is problem with opening the plugin", func() {
			BeforeEach(func() {
				opener.OpenReturns(nil, errors.New("unable to open plugin"))
			})

			It("should have returned an error", func() {
				Ω(err).Should(HaveOccurred())
			})

			It("should have returned nil handler", func() {
				Ω(handler).Should(BeNil())
			})
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should have opened the plugins in reverse order", func() {
			Ω(opener.OpenCallCount()).Should(Equal(2))

			nameArg, configArg, nextArg := opener.OpenArgsForCall(0)
			Ω(nameArg).Should(Equal("mighty-grasshopper"))
			Ω(configArg).Should(Equal([]byte(`fly: "no"` + "\n")))
			Ω(nextArg).Should(BeNil())

			nameArg, configArg, nextArg = opener.OpenArgsForCall(1)
			Ω(nameArg).Should(Equal("happy-unicorn"))
			Ω(configArg).Should(Equal([]byte("{}\n")))
			Ω(nextArg).ShouldNot(BeNil())
		})

		It("should have not returned nil", func() {
			Ω(handler).ShouldNot(BeNil())
		})
	})
})
