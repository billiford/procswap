package procswap_test

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"

	. "github.com/billiford/procswap/pkg"
	. "github.com/onsi/gomega"
)

var _ = Describe("Swap", func() {
	var (
		err  error
		path string
		swap Swap
	)

	BeforeEach(func() {
		// This file just prints out hello world.
		path = swapFilePath()
	})

	JustBeforeEach(func() {
		swap = NewSwap(path)
		err = swap.Start()
	})

	Describe("#Path", func() {
		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(swap.Path()).To(Equal(path))
			})
		})
	})

	Describe("#PID", func() {
		When("there is no underlying command", func() {
			JustBeforeEach(func() {
				swap.Kill()
				swap = NewSwap(path)
			})

			It("returns -1", func() {
				Expect(swap.PID()).To(Equal(-1))
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(swap.PID() > 0).To(BeTrue())
			})
		})
	})

	Describe("#Start", func() {
		When("the path does not exist", func() {
			BeforeEach(func() {
				path = filepath.FromSlash(currentDir() + "/" + uuid.New().String())
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("#Kill", func() {
		BeforeEach(func() {
			path = waitFilePath()
		})

		When("there is no underlying command", func() {
			JustBeforeEach(func() {
				err = swap.Kill()
				Expect(err).To(BeNil())
				swap = NewSwap(path)
				err = swap.Kill()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("no command to kill"))
			})
		})

		When("it succeeds", func() {
			JustBeforeEach(func() {
				err = swap.Kill()
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})

func currentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

func swapFilePath() string {
	var swap string

	if runtime.GOOS == "windows" {
		swap = filepath.FromSlash(currentDir() + "/test/swaps/swap.exe")
	} else {
		swap = filepath.FromSlash(currentDir() + "/test/swaps/swap")
	}

	return swap
}

func priorityFileDir() string {
	return filepath.FromSlash(currentDir() + "/test/priorities")
}

func priorityFilePath() string {
	var priority string

	if runtime.GOOS == "windows" {
		priority = filepath.FromSlash(currentDir() + "/test/priorities/wait.exe")
	} else {
		priority = filepath.FromSlash(currentDir() + "/test/priorities/wait_linux.exe")
	}

	return priority
}

func priorityScriptPath() string {
	var priorityScript string

	if runtime.GOOS == "windows" {
		priorityScript = filepath.FromSlash(currentDir() + "/test/scripts/priority-script.exe")
	} else {
		priorityScript = filepath.FromSlash(currentDir() + "/test/scripts/priority-script")
	}

	return priorityScript
}

func waitFilePath() string {
	var wait string

	if runtime.GOOS == "windows" {
		wait = filepath.FromSlash(currentDir() + "/test/wait/wait.exe")
	} else {
		wait = filepath.FromSlash(currentDir() + "/test/wait/wait_linux.exe")
	}

	return wait
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
