package procswap_test

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"

	. "github.com/billiford/procswap/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	var (
		infos  []os.FileInfo
		path   string
		err    error
		rescue *os.File
	)

	BeforeEach(func() {
		// Just set the path to the current directory since we know that exists.
		path, err = os.Getwd()
		Expect(err).To(BeNil())
		rescue = os.Stdout
		os.Stdout = os.NewFile(0, os.DevNull)
	})

	AfterEach(func() {
		os.Stdout = rescue
	})

	JustBeforeEach(func() {
		infos, err = ProcessList(path)
	})

	Describe("#ProcessList", func() {
		When("the path does not exist", func() {
			BeforeEach(func() {
				path = uuid.New().String()
			})

			It("returns an error", func() {
				Expect(err).ToNot(BeNil())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})

		Context("when the path is a file", func() {
			BeforeEach(func() {
				path = filepath.FromSlash(path + "/test/file")
			})

			It("succeeds", func() {
				Expect(err).To(BeNil())
				Expect(infos).To(HaveLen(1))
			})
		})

		Context("when the path is a directory", func() {
			BeforeEach(func() {
				path = filepath.FromSlash(path + "/test")
			})

			It("walks the directory returning all .exe files", func() {
				Expect(err).To(BeNil())
				Expect(infos).To(HaveLen(2))
			})
		})
	})
})
