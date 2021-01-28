package procswap_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/billiford/procswap/pkg"
)

var _ = Describe("Swap", func() {
	var (
		swap       Swap
		currentDir string
		err        error
		path       string
	)

	BeforeEach(func() {
		// Setup priority executables.
		currentDir, err = os.Getwd()
		Expect(err).To(BeNil())

		path = swapFilePath(currentDir)
	})

	JustBeforeEach(func() {
		swap = NewSwap(path)
	})

	Describe("#PID", func() {
		JustBeforeEach(func() {
			err = swap.Start()
			Expect(err).To(BeNil())
		})

		It("returns the PID", func() {
			Expect(swap.PID()).ToNot(BeZero())
		})
	})

	Describe("#Path", func() {
		It("returns the path", func() {
			Expect(swap.Path()).To(Equal(swapFilePath(currentDir)))
		})
	})

	Describe("#Cmd", func() {
		JustBeforeEach(func() {
			err = swap.Start()
			Expect(err).To(BeNil())
		})

		It("returns the underlying cmd", func() {
			Expect(swap.Cmd()).ToNot(BeNil())
		})
	})
})
