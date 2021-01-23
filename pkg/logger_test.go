package procswap

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Logger", func() {
	var (
		message      string
		buffer       *Buffer
		rescue, r, w *os.File
	)

	BeforeEach(func() {
		message = "test"
		rescue = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
		buffer = BufferReader(r)
	})

	AfterEach(func() {
		w.Close()
		os.Stdout = rescue
	})

	When("there is no newline", func() {
		JustBeforeEach(func() {
			logDebug(message, false)
		})

		It("logs the message without a newline on the end", func() {
			Eventually(buffer).Should(Say(`\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*DEBUG.*\| test$`))
		})
	})

	Describe("#logDebug", func() {
		JustBeforeEach(func() {
			logDebug(message)
		})

		It("logs the message", func() {
			Eventually(buffer).Should(Say(`\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*DEBUG.*\| test\n$`))
		})
	})

	Describe("#logInfo", func() {
		JustBeforeEach(func() {
			logInfo(message)
		})

		It("logs the message", func() {
			Eventually(buffer).Should(Say(`\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*INFO.*\| test\n$`))
		})
	})

	Describe("#logWarn", func() {
		JustBeforeEach(func() {
			logWarn(message)
		})

		It("logs the message", func() {
			Eventually(buffer).Should(Say(`\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*WARN.*\| test\n$`))
		})
	})

	Describe("#logError", func() {
		JustBeforeEach(func() {
			logError(message)
		})

		It("logs the message", func() {
			Eventually(buffer).Should(Say(`\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*ERROR.*\| test\n$`))
		})
	})
})
