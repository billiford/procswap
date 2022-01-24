package procswap_test

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/karrick/godirwalk"

	. "github.com/billiford/procswap/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path", func() {
	var (
		infos   []*godirwalk.Dirent
		ignored []string
		path    string
		err     error
		rescue  *os.File
	)

	BeforeEach(func() {
		// Just set the path to the current directory since we know that exists.
		path, err = os.Getwd()
		Expect(err).To(BeNil())
		rescue = os.Stdout
		os.Stdout = os.NewFile(0, os.DevNull)
		ignored = []string{}
	})

	AfterEach(func() {
		os.Stdout = rescue
	})

	JustBeforeEach(func() {
		infos, err = ProcessList(path, ignored)
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
				Expect(infos).ToNot(HaveLen(0))
			})
		})

		Context("when the path is a directory", func() {
			BeforeEach(func() {
				path = filepath.FromSlash(path + "/test/priorities")
			})

			It("walks the directory returning all .exe files", func() {
				Expect(err).To(BeNil())
				Expect(infos).ToNot(HaveLen(0))
			})
		})
	})
})
