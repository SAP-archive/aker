package config_test

import (
	. "github.infra.hana.ondemand.com/I061150/aker/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("LoadFromFile", func() {

		var config Config
		var err error

		Context("when the config file is valid", func() {
			BeforeEach(func() {
				config, err = LoadFromFile("valid_config.yml")
			})

			It("should have not returned an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have proper server config section", func() {
				Ω(config.Server.Host).Should(Equal("localhost"))
				Ω(config.Server.Port).Should(Equal(8080))
			})

			It("should have proper endpoint section", func() {
				Ω(len(config.Endpoints)).Should(Equal(2))
				Ω(config.Endpoints[0]).Should(Equal(EndpointConfig{
					Path:    "/",
					Audit:   false,
					Plugins: []PluginReferenceConfig{},
				}))
				Ω(config.Endpoints[1]).Should(Equal(EndpointConfig{
					Path:  "/proxy",
					Audit: true,
					Plugins: []PluginReferenceConfig{
						PluginReferenceConfig{
							Name: "aker-proxy",
							Config: map[string]interface{}{
								"url": "http://location.com",
							},
						}},
				}))
			})
		})

		Context("when the config file is invalid", func() {
			BeforeEach(func() {
				config, err = LoadFromFile("./invalid_config.yml")
			})

			It("should have returned an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
	})
})
