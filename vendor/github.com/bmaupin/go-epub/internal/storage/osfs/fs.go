// Package osfs implements the Storage interface for os' filesystems

package osfs

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bmaupin/go-epub/internal/storage"
)

type OSFS struct {
	rootDir string
	fs.FS
}

func NewOSFS(rootDir string) *OSFS {
	return &OSFS{
		rootDir: rootDir,
		FS:      os.DirFS(rootDir),
	}
}

func (o *OSFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(o.rootDir, name), data, perm)
}

func (o *OSFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(filepath.Join(o.rootDir, name), perm)
}

func (o *OSFS) RemoveAll(name string) error {
	return os.RemoveAll(filepath.Join(o.rootDir, name))
}

func (o *OSFS) Create(name string) (storage.File, error) {
	return os.Create(filepath.Join(o.rootDir, name))
}

func (o *OSFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(o.rootDir, name))
}

func (o *OSFS) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(o.rootDir, name))
}
