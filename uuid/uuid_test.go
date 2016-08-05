package uuid_test

import (
	"sync"

	. "github.infra.hana.ondemand.com/cloudfoundry/aker/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("UUID", func() {
	DescribeTable("String", func(uuid UUID, expected string) {
		Ω(uuid.String()).Should(Equal(expected))
	},
		Entry("Zero", UUID{},
			"00000000-0000-0000-0000-000000000000"),
		Entry("Random",
			UUID{0x44, 0x3, 0x60, 0x82, 0xc7, 0x9b, 0x42, 0x78,
				0x37, 0xa5, 0x32, 0x5c, 0x1d, 0xbb, 0xb0, 0xd6},
			"44036082-c79b-4278-37a5-325c1dbbb0d6"),
		Entry("Random-2",
			UUID{0xa5, 0xa8, 0x9c, 0x62, 0x7a, 0xe4, 0x2d, 0xdc,
				0x13, 0x86, 0xd7, 0xa3, 0xc8, 0xe, 0x24, 0x85},
			"a5a89c62-7ae4-2ddc-1386-d7a3c80e2485"),
		Entry("Last",
			UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			"ffffffff-ffff-ffff-ffff-ffffffffffff"),
	)

	Describe("Random", func() {
		It("should generate a UUID", func() {
			uid, err := Random()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(uid.String()).ShouldNot(BeEmpty())
			Ω(uid.String()).Should(HaveLen(36))
		})

		It("should return random UUID", func() {
			first, err := Random()
			Ω(err).ShouldNot(HaveOccurred())
			second, err := Random()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(first.String()).ShouldNot(Equal(second.String()))
		})
	})

	Measure("Performance", func(b Benchmarker) {
		b.Time("runtime", func() {
			wait := &sync.WaitGroup{}
			for i := 0; i < 50000; i++ {
				wait.Add(1)
				go func() {
					Random()
					wait.Done()
				}()
			}
			wait.Wait()
		})
	}, 100)
})
