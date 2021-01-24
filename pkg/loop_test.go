package procswap_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	. "github.com/billiford/procswap/pkg"
)

var _ = Describe("Loop", func() {
	const (
		fmtInfoLog  = `\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*INFO.*\| `
		fmtWarnLog  = `\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*WARN.*\| `
		fmtErrorLog = `\[PROCSWAP\] \d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*ERROR.*\| `
	)

	var (
		currentDir     string
		prioritiesPath string
		err            error
		loop           Loop
		buffer         *Buffer
		rescue, r, w   *os.File
	)

	BeforeEach(func() {
		// log.SetOutput(ioutil.Discard)

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
		prioritiesPath = filepath.FromSlash(currentDir + "/test/priorities")

		execs, err := ProcessList(prioritiesPath)
		Expect(err).To(BeNil())

		loop.WithPriorities(execs)

		swaps := []string{
			swapFilePath(currentDir),
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
		Context("when it is the first loop and a priority process is already running", func() {
			var cmd *exec.Cmd

			BeforeEach(func() {
				file := priorityFilePath(currentDir)
				cmd = exec.Command(file)
				err := cmd.Start()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				err := cmd.Process.Kill()
				Expect(err).To(BeNil())
				cmd.Wait()
			})

			It("lets us know that priorities are already running", func() {
				Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables`))
				Eventually(buffer).Should(Say(fmtWarnLog + `not starting swap processes, priority processes already running: .*`))
			})
		})

		Context("when the swaps have started and then a priority is started", func() {
			var cmd *exec.Cmd

			BeforeEach(func() {
				loop.WithLimit(2)
				loop.WithPollInterval(1)

				go func() {
					time.Sleep(500 * time.Millisecond)
					// start a priority process
					file := priorityFilePath(currentDir)
					cmd = exec.Command(file)
					err := cmd.Start()
					Expect(err).To(BeNil())
				}()
			})

			AfterEach(func() {
				err := cmd.Process.Kill()
				Expect(err).To(BeNil())
				cmd.Wait()
			})

			It("lets us know it is stopping the running processes", func() {
				Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath(currentDir) + `.*\.\.\. .*OK.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + priorityFile() + `.*`))
				Eventually(buffer).Should(Say(fmtInfoLog + `.*stop.* .*` + swapFilePath(currentDir) + `.*\.\.\. .*OK.*`))
			})
		})

		Context("when there are no running priorities and swap processes have not been started", func() {
			When("you pass in a swap file that doesn't exist", func() {
				var swap, guid string

				BeforeEach(func() {
					guid = uuid.New().String()
					swap = filepath.FromSlash(currentDir + "/" + guid)
					loop.WithSwaps([]string{swap})
				})

				It("prints some errors", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swap + `.*\.\.\. .*FAILED.*`))
					// The error will likely change cross platform, so don't test too much.
					Eventually(buffer).Should(Say(fmtErrorLog + `error starting swap process ` + swap + ".*"))
				})
			})

			When("it runs", func() {
				It("runs", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `searching ` + prioritiesPath + ` for executables`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath(currentDir) + `.*\.\.\. .*OK.*`))
				})
			})
		})
	})
})

func swapFilePath(currentDir string) string {
	var swap string
	if runtime.GOOS == "windows" {
		swap = filepath.FromSlash(currentDir + "/test/swaps/swap.exe")
	} else {
		swap = filepath.FromSlash(currentDir + "/test/swaps/swap")
	}

	return swap
}

func priorityFilePath(currentDir string) string {
	var priority string
	if runtime.GOOS == "windows" {
		priority = filepath.FromSlash(currentDir + "/test/priorities/wait.exe")
	} else {
		priority = filepath.FromSlash(currentDir + "/test/priorities/wait_linux.exe")
	}

	return priority
}

func priorityFile() string {
	var priority string
	if runtime.GOOS == "windows" {
		priority = "wait.exe"
	} else {
		priority = "wait_linux.exe"
	}

	return priority
}
