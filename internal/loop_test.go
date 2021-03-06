package procswap_test

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	"github.com/billiford/go-ps"
	gopsfakes "github.com/billiford/go-ps/go-psfakes"
	. "github.com/billiford/procswap/internal"
	"github.com/billiford/procswap/internal/internalfakes"
)

var _ = Describe("Loop", func() {
	const (
		fmtInfoLog  = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*INFO.*\| `
		fmtWarnLog  = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*WARN.*\| `
		fmtErrorLog = `\d{4}\/\d{2}\/\d{2} - \d{2}:\d{2}:\d{2} \|.*ERROR.*\| `
	)

	var (
		fakePs         *gopsfakes.FakePs
		fakeProcess    *gopsfakes.FakeProcess
		fakeSwap       *internalfakes.FakeSwap
		prioritiesPath string
		ignored        []string
		err            error
		loop           Loop
		buffer         *Buffer
		rescue, r, w   *os.File
	)

	BeforeEach(func() {
		loop = NewLoop()
		loop.WithLimit(1)
		loop.WithPollInterval(0)

		// Output checks.
		rescue = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
		buffer = BufferReader(r)

		// Setup priority executables.
		Expect(err).To(BeNil())
		prioritiesPath = filepath.FromSlash(currentDir() + "/test/priorities")

		ignored = []string{}
		execs, err := ProcessList(prioritiesPath, ignored)
		Expect(err).To(BeNil())

		loop.WithPriorities(execs)

		fakeSwap = &internalfakes.FakeSwap{}
		fakeSwap.PathReturns(swapFilePath())
		swaps := []Swap{
			fakeSwap,
		}

		loop.WithSwaps(swaps)

		fakePs = &gopsfakes.FakePs{}
		fakeProcess = &gopsfakes.FakeProcess{}
		fakePs.ProcessesReturns([]ps.Process{fakeProcess}, nil)
		loop.WithPs(fakePs)
		loop.WithActionsEnabled(false)
	})

	JustBeforeEach(func() {
		loop.Run()
	})

	AfterEach(func() {
		w.Close()
		os.Stdout = rescue
	})

	Describe("#Run", func() {
		When("listing processes returns an error", func() {
			BeforeEach(func() {
				fakePs.ProcessesReturns(nil, errors.New("error listing processes"))
			})

			It("logs the error", func() {
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
				Eventually(buffer).Should(Say(fmtErrorLog + `error listing currently running processes: error listing processes`))
			})
		})

		Context("when it is the first loop and a priority process is already running", func() {
			BeforeEach(func() {
				fakeProcess.ExecutableReturns(priorityFile())
				processes := []ps.Process{
					fakeProcess,
				}
				fakePs.ProcessesReturns(processes, nil)
			})

			It("lets us know that priorities are already running", func() {
				Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
				Eventually(buffer).Should(Say(fmtWarnLog + `not starting swap processes, priority processes already running: .*`))
			})
		})

		Context("when the swaps have started and then a priority is started", func() {
			BeforeEach(func() {
				loop.WithLimit(2)
				loop.WithPollInterval(1)

				go func() {
					time.Sleep(500 * time.Millisecond)
					fakeProcess.ExecutableReturns(priorityFile())
					processes := []ps.Process{
						fakeProcess,
					}
					fakePs.ProcessesReturns(processes, nil)
				}()
			})

			Context("when there is a priority script set", func() {
				var priorityScript string

				BeforeEach(func() {
					priorityScript = priorityScriptPath()
					loop.WithPriorityScript(priorityScript)
				})

				When("the priority script fails to run", func() {
					BeforeEach(func() {
						// Just pass in a path that doesn't exist to generate a failure.
						priorityScript = filepath.FromSlash(currentDir() + "/" + uuid.New().String())
						loop.WithPriorityScript(priorityScript)
					})

					It("lets us know and logs the error", func() {
						Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*priority.* .*` + priorityFile() + `.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*stop.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*priority script.* .*` + priorityScript + `.*\.\.\. .*FAILED.*`))
						Eventually(buffer).Should(Say(fmtErrorLog + `.*` + priorityScript + `.*`))
					})
				})

				When("it succeeds", func() {
					It("lets us know the priority script has started successfully", func() {
						Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*priority.* .*` + priorityFile() + `.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*stop.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
						Eventually(buffer).Should(Say(fmtInfoLog + `.*priority script.* .*` + priorityScriptPath() + `.*\.\.\. .*OK.*`))
					})
				})
			})

			When("stopping a swap process fails", func() {
				BeforeEach(func() {
					fakeSwap.KillReturns(errors.New("error stopping swap"))
				})

				It("logs the error", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*priority.* .*` + priorityFile() + `.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*stop.* .*` + swapFilePath() + `.*\.\.\. .*FAILED.*`))
					Eventually(buffer).Should(Say(fmtErrorLog + `error stopping swap`))
				})
			})

			When("it succeeds", func() {
				It("lets us know it is stopping the running processes", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*priority.* .*` + priorityFile() + `.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*stop.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
				})
			})
		})

		Context("when there are no running priorities and swap processes have not been started", func() {
			When("you pass in a swap file that doesn't exist", func() {
				BeforeEach(func() {
					fakeSwap.StartReturns(errors.New("exec error"))
				})

				It("prints some errors", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables.*`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + fakeSwap.Path() + `.*\.\.\. .*FAILED.*`))
					// The error will likely change cross platform, so don't test too much.
					Eventually(buffer).Should(Say(fmtErrorLog + `error starting swap process ` + fakeSwap.Path() + ": exec error"))
				})
			})

			When("it runs", func() {
				It("runs", func() {
					Eventually(buffer).Should(Say(fmtInfoLog + `.*setup.* searching ` + prioritiesPath + ` for executables`))
					Eventually(buffer).Should(Say(fmtInfoLog + `.*start.* .*` + swapFilePath() + `.*\.\.\. .*OK.*`))
				})
			})
		})
	})
})
