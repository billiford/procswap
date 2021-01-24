package procswap_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	. "github.com/billiford/procswap/pkg"
)

var _ = Describe("Loop", func() {
	const (
		fmtInfoLog  = `\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*INFO.*\| `
		fmtErrorLog = `\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*ERROR.*\| `
	)

	var (
		currentDir   string
		err          error
		loop         Loop
		buffer       *Buffer
		rescue, r, w *os.File
	)

	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)

		loop = NewLoop()
		loop.WithLimit(1)
		loop.WithPollInterval(0)

		// Output checks.
		rescue = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
		buffer = BufferReader(r)

		// Setup priority executables.
		currentDir, err = os.Getwd()
		Expect(err).To(BeNil())
		prioritiesPath := filepath.FromSlash(currentDir + "/test/priorities")

		execs, err := ProcessList(prioritiesPath)
		Expect(err).To(BeNil())

		loop.WithPriorities(execs)

		swaps := []string{
			swapFile(currentDir),
		}

		loop.WithSwaps(swaps)
	})

	JustBeforeEach(func() {
		loop.Run()
	})

	AfterEach(func() {
		w.Close()
		os.Stdout = rescue
	})

	Describe("#Run", func() {
		When("you pass in a swap file that doesn't exist", func() {
			var swap, guid string

			BeforeEach(func() {
				guid = uuid.New().String()
				swap = filepath.FromSlash(currentDir + "/" + guid)
				loop.WithSwaps([]string{swap})
			})

			It("prints some errors", func() {
				prioritiesPath := filepath.FromSlash(currentDir + "/test/priorities")
				Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `no priority processes running, starting all swap processes.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + guid + `.*`))
				Eventually(buffer).Should(Say(fmtErrorLog + `error starting swap process .*`))
			})
		})

		When("it runs", func() {
			It("runs", func() {
				prioritiesPath := filepath.FromSlash(currentDir + "/test/priorities")
				Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables`))
				Eventually(buffer).Should(Say(fmtInfoLog + `no priority processes running, starting all swap processes`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFile(currentDir) + `.*`))
			})
		})
	})
})

func swapFile(currentDir string) string {
	var swap string
	if runtime.GOOS == "windows" {
		swap = filepath.FromSlash(currentDir + "/test/swaps/swap.exe")
	} else {
		swap = filepath.FromSlash(currentDir + "/test/swaps/swap")
	}

	return swap
}
