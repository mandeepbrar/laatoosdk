package utils

import "os"

func FileExists(file string) (bool, os.FileInfo, error) {
	inf, err := os.Stat(file)
	return err == nil, inf, err
}
