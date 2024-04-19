package filestorage

import "os"

func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}

	return os.MkdirAll(fp, os.ModePerm)
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(fp string) bool {
	_, err := os.Stat(fp)

	return err == nil || os.IsExist(err)
}

// Remove one file
func Remove(name string) error {
	return os.Remove(name)
}
