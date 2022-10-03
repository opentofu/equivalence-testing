package files

import (
	"io"
	"io/fs"
	"os"
	"strings"
)

// CopyDir should be used in conjunction with filepath.WalkDir to recursively
// copy all the files within sourceDirectory into targetDirectory.
//
// Any file names in skipFiles will be skipped.
func CopyDir(sourceDirectory, targetDirectory string, skipFiles []string) fs.WalkDirFunc {
	return func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, skip := range skipFiles {
			if skip == entry.Name() {
				return nil
			}
		}

		if path == sourceDirectory {
			return nil
		}

		targetFile := strings.ReplaceAll(path, sourceDirectory, targetDirectory)

		if entry.IsDir() {
			return os.MkdirAll(targetFile, os.ModePerm)
		}

		source, err := os.Open(path)
		if err != nil {
			return err
		}
		defer source.Close()

		target, err := os.Create(targetFile)
		if err != nil {
			return err
		}
		defer target.Close()

		_, err = io.Copy(target, source)
		return err
	}
}
