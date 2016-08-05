package plugin_test

import (
	. "github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testConfig struct {
	A int    `yaml:"a"`
	B string `yaml:"b"`
}

var _ = Describe("Config", func() {

	Describe("MarshalConfig", func() {

		var config testConfig
		var data []byte
		var err error

		BeforeEach(func() {
			config = testConfig{
				A: 42,
				B: "hooray!",
			}
			data, err = MarshalConfig(config)
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should have returned data", func() {
			Ω(data).ShouldNot(BeEmpty())
		})

		It("should be possible to marshal it back and get the same struct", func() {
			var unmarshaled testConfig
			err = UnmarshalConfig(data, &unmarshaled)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(unmarshaled).Should(Equal(config))
		})
	})

	Describe("UnmarshalConfig", func() {

		var config testConfig
		var data []byte
		var err error

		BeforeEach(func() {
			data = []byte("a: 42\nb: hooray\n")
			err = UnmarshalConfig(data, &config)
		})

		It("should have not returned an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should set proper struct values", func() {
			Ω(config.A).Should(Equal(42))
			Ω(config.B).Should(Equal("hooray"))
		})
	})

})
