package procswap_test

import (
	"os"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/urfave/cli/v2"

	. "github.com/billiford/procswap/internal"
)

var _ = Describe("App", func() {
	const (
		fmtInfoLog  = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*INFO.*\| `
		fmtWarnLog  = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*WARN.*\| `
		fmtErrorLog = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*ERROR.*\| `
	)

	var (
		err          error
		app          *cli.App
		args         []string
		buffer       *Buffer
		rescue, r, w *os.File
	)

	Describe("#NewApp", func() {
		BeforeEach(func() {
			app = NewApp()
		})

		It("returns an app", func() {
			Expect(app.Name).To(Equal("procswap"))
		})
	})

	Describe("#Run", func() {
		BeforeEach(func() {
			app = NewApp()
			args = []string{procswapFilename(),
				"-p", priorityFileDir(),
				"-s", swapFilePath(),
				"-ps", priorityScriptPath(),
				"--limit", "1",
				"--poll-interval", "1",
				"--ignore", "ignore_me.exe",
				"--disable-actions",
			}

			// Output checks.
			rescue = os.Stdout
			r, w, _ = os.Pipe()
			os.Stdout = w
			buffer = BufferReader(r)
		})

		AfterEach(func() {
			w.Close()
			os.Stdout = rescue
		})

		JustBeforeEach(func() {
			err = app.Run(args)
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + priorityFileDir() + ` for executables`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* ignoring priority .*ignore_me.exe.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* found .*\d.* priority executables`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* registered .*\d.* swap processes`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* registered priority script .*` + priorityScriptPath() + `.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
			})
		})
	})
})

func procswapFilename() string {
	if runtime.GOOS == "windows" {
		return "procswap.exe"
	}

	return "procswap"
}
