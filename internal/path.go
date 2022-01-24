package procswap

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/logrusorgru/aurora"
)

// ProcessList lists all .exe files in a given directory.
func ProcessList(path string, ignored []string) ([]*godirwalk.Dirent, error) {
	files := []*godirwalk.Dirent{}

	// Check to make sure it exists first.
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}

		return nil, fmt.Errorf("error checking if %s exists: %w", path, err)
	}
	// Don't check if it's a .exe if the user has passed in a
	// single file as priority.
	if !info.IsDir() {
		de, err := godirwalk.NewDirent(path)
		if err != nil {
			return nil, fmt.Errorf("error creating new dirent for priority file %s: %w", path, err)
		}

		files = append(files, de)

		return files, nil
	}

	logInfo(fmt.Sprintf("%s searching %s for executables", aurora.Cyan("setup"), path))

	// Only list files that end in '.exe'.
	libRegEx := regexp.MustCompile("^.*.exe$")

	err = godirwalk.Walk(path, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if libRegEx.MatchString(de.Name()) {
				if contains(ignored, de.Name()) {
					logInfo(fmt.Sprintf("%s ignoring priority %s", aurora.Cyan("setup"), aurora.Bold(de.Name())))
				} else {
					files = append(files, de)
				}
			}

			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error walking %s searching for .exes: %w", path, err)
	}

	return files, nil
}

// contains returns true if slice s contains element e, ignoring case.
func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}

	return false
}
