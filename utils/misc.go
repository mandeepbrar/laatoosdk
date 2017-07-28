package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func FileExists(file string) (bool, os.FileInfo, error) {
	inf, err := os.Stat(file)
	return err == nil, inf, err
}

func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return
}

func CopyDir(source string, dest string, prefix string) (err error) {
	// get properties of source dir
	_, fi, err := FileExists(source)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return
	}

	for _, entry := range entries {
		sfp := path.Join(source, entry.Name())
		dfp := path.Join(dest, fmt.Sprintf("%s%s", prefix, entry.Name()))
		if entry.IsDir() {
			err = CopyDir(sfp, dfp, prefix)
			if err != nil {
				return
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err != nil {
				return
			}
		}
	}
	return
}
