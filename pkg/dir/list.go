package dir

import (
	"os"
	"path/filepath"
	"regexp"
)

func ListForExecutables(directory string) ([]os.FileInfo, error) {
	files := []os.FileInfo{}

	libRegEx, err := regexp.Compile("^.*.exe$")
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			files = append(files, info)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
