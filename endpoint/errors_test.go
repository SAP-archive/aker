package endpoint_test

import (
	. "github.infra.hana.ondemand.com/I061150/aker/endpoint"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	Describe("InvalidPathError", func() {
		It("should return proper error message", func() {
			err := InvalidPathError("path")
			Ω(err.Error()).Should(Equal(`invalid endpoint path: "path"`))
		})
	})
})
