package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func FileExists(file string) (bool, os.FileInfo, error) {
	inf, err := os.Stat(file)
	return err == nil, inf, err
}

// ReadPackedFile is a function to unpack a tar.gz
func ReadPackedFile(filepath string) {
	if filepath == "" {
		panic("Empty input!")
	}

	processFile(filepath)
}

func processFile(srcFile string) {

	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarReader := tar.NewReader(gzf)
	// defer io.Copy(os.Stdout, tarReader)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := header.Name

		switch header.Typeflag {
		case tar.TypeDir: // = directory
			fmt.Println("Directory:", name)
			os.Mkdir(name, 0755)
		case tar.TypeReg: // = regular file
			fmt.Println("Regular file:", name)
			data := make([]byte, header.Size)
			_, err := tarReader.Read(data)
			if err != nil {
				panic("Error reading file!!! PANIC!!!!!!")
			}

			ioutil.WriteFile(name, data, 0755)
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
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
