package logging_test

import (
	"bytes"

	. "github.infra.hana.ondemand.com/cloudfoundry/aker/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NativeLogger", func() {

	var log Logger
	var stdout *bytes.Buffer
	var stderr *bytes.Buffer

	itShouldHaveLoggedToStdout := func(message string) {
		It("should have logged to stdout", func() {
			Ω(stdout.Bytes()).Should(ContainSubstring(message))
		})
	}

	itShouldHaveLoggedToStderr := func(message string) {
		It("should have logged to stderr", func() {
			Ω(stderr.Bytes()).Should(ContainSubstring(message))
		})
	}

	BeforeEach(func() {
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
		log = NewNativeLogger(stdout, stderr)
	})

	Describe("Debugf", func() {
		BeforeEach(func() {
			log.Debugf("debug message")
		})

		itShouldHaveLoggedToStdout("debug message")
	})

	Describe("Infof", func() {
		BeforeEach(func() {
			log.Infof("info message")
		})

		itShouldHaveLoggedToStdout("info message")
	})

	Describe("Warnf", func() {
		BeforeEach(func() {
			log.Warnf("warn message")
		})

		itShouldHaveLoggedToStdout("warn message")
	})

	Describe("Errorf", func() {
		BeforeEach(func() {
			log.Errorf("error message")
		})

		itShouldHaveLoggedToStderr("error message")
	})
})
