package procswap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/logrusorgru/aurora"
)

// ProcessList lists all .exe files in a
// given directory.
func ProcessList(path string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}

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
		files = append(files, info)

		return files, nil
	}

	logInfo(fmt.Sprintf("%s searching %s for executables", aurora.Cyan("setup"), path))

	// Only list files that end in '.exe'.
	libRegEx := regexp.MustCompile("^.*.exe$")

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			files = append(files, info)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking %s searching for .exes: %w", path, err)
	}

	return files, nil
}
