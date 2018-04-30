package cmd

import "os"

func validatePath(path string) bool {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return f.IsDir()
}

func validateFile(path string) bool {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !f.IsDir()
}
