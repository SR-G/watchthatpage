package core

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Exists allow to check if a file or directory is available on filesystem
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// RemoveAllContentFromDirectory will remove all files from a given directory
func RemoveAllContentFromDirectory(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadGzFile loads a gzipped file content
func ReadGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateDirectoryIfNeeded creates the expected directory if not already available
func CreateDirectoryIfNeeded(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}
