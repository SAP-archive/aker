package config_test

import (
	. "github.infra.hana.ondemand.com/cloudfoundry/aker/config"

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
				Ω(config.Server.ReadTimeout).Should(Equal(5))
				Ω(config.Server.WriteTimeout).Should(Equal(10))
			})

			It("should have proper endpoint section", func() {
				Ω(len(config.Endpoints)).Should(Equal(2))
				Ω(config.Endpoints[0]).Should(Equal(Endpoint{
					Path:    "/",
					Audit:   false,
					Plugins: []PluginReference{},
				}))
				Ω(config.Endpoints[1]).Should(Equal(Endpoint{
					Path:  "/proxy",
					Audit: true,
					Plugins: []PluginReference{
						PluginReference{
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
				_, err = LoadFromFile("./invalid_config.yml")
			})

			It("should have returned an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when there is problem opening the config file", func() {
			BeforeEach(func() {
				_, err = LoadFromFile("/not/existing")
			})

			It("should have returned an error", func() {
				Ω(err).Should(HaveOccurred())
			})
		})
	})
})
