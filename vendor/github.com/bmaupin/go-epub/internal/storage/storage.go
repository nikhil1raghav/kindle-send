// Package storage hold and abstraction of the filesystem

package storage

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Storage is an abstraction of the filesystem
type Storage interface {
	fs.FS
	// WriteFile writes data to the named file, creating it if necessary. If the file does not exist, WriteFile creates it with permissions perm (before umask); otherwise WriteFile truncates it before writing, without changing permissions.
	WriteFile(name string, data []byte, perm fs.FileMode) error
	// Mkdir creates a new directory with the specified name and permission bits (before umask). If there is an error, it will be of type *PathError.
	Mkdir(name string, perm fs.FileMode) error
	// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters. If the path does not exist, RemoveAll returns nil (no error). If there is an error, it will be of type *PathError.
	RemoveAll(name string) error
	// Create creates or truncates the named file. If the file already exists, it is truncated. If the file does not exist, it is created with mode 0666 (before umask). If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. If there is an error, it will be of type *PathError.
	Create(name string) (File, error)
}

type File interface {
	fs.File
	io.Writer
}

// ReadFile returns the content of name in the filesystem
func ReadFile(fs Storage, name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
func MkdirAll(fs Storage, dir string, perm fs.FileMode) error {
	list := make([]string, 0)
	stop := ""
	for dir := filepath.Dir(dir); dir != stop; dir = filepath.Dir(dir) {
		list = append(list, dir)
		stop = dir
	}
	for i := len(list); i > 0; i-- {
		err := fs.Mkdir(list[i-1], perm)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}
	return nil

}
