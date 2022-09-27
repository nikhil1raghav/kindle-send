package memory

import (
	"io/fs"
	"path"
	"strings"
	"time"

	"github.com/bmaupin/go-epub/internal/storage"
)

type Memory struct {
	fs map[string]*file
}

func NewMemory() *Memory {
	return &Memory{
		fs: map[string]*file{
			"/": {
				name:    path.Base("/"),
				modTime: time.Now(),
				mode:    fs.ModeDir | (0666),
			},
		},
	}
}

// Open opens the named file.
//
// When Open returns an error, it should be of type *PathError
// with the Op field set to "open", the Path field set to name,
// and the Err field describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// ValidPath(name), returning a *PathError with Err set to
// ErrInvalid or ErrNotExist.
func (m *Memory) Open(name string) (fs.File, error) {
	var f fs.File
	var ok bool
	if f, ok = m.fs[name]; !ok {
		return nil, fs.ErrNotExist
	}
	return f, nil
}

// WriteFile writes data to the named file, creating it if necessary. If the file does not exist, WriteFile creates it with permissions perm (before umask); otherwise WriteFile truncates it before writing, without changing permissions.
func (m *Memory) WriteFile(name string, data []byte, perm fs.FileMode) error {
	if !fs.ValidPath(name) {
		return fs.ErrInvalid
	}
	f := &file{
		name:    path.Base(name),
		modTime: time.Now(),
		mode:    (perm),
		content: data,
	}
	m.fs[name] = f
	return nil
}

// Mkdir creates a new directory with the specified name and permission bits (before umask). If there is an error, it will be of type *PathError.
func (m *Memory) Mkdir(name string, perm fs.FileMode) error {
	if !fs.ValidPath(path.Base(name)) {
		return fs.ErrInvalid
	}
	f := &file{
		name:    path.Base(name),
		modTime: time.Now(),
		mode:    fs.ModeDir | (perm),
	}
	m.fs[name] = f
	return nil
}

// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters. If the path does not exist, RemoveAll returns nil (no error). If there is an error, it will be of type *PathError.
func (m *Memory) RemoveAll(name string) error {
	for k := range m.fs {
		if strings.HasPrefix(k, name) {
			delete(m.fs, k)
		}
	}
	return nil
}

// Create creates or truncates the named file. If the file already exists, it is truncated. If the file does not exist, it is created with mode 0666 (before umask). If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. If there is an error, it will be of type *PathError.
func (m *Memory) Create(name string) (storage.File, error) {
	if !fs.ValidPath(path.Base(name)) {
		return nil, fs.ErrInvalid
	}
	f := &file{
		name:    path.Base(name),
		modTime: time.Now(),
		mode:    0666,
	}
	m.fs[name] = f
	return f, nil
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
func (m *Memory) ReadDir(name string) ([]fs.DirEntry, error) {
	output := make([]fs.DirEntry, 0)
	for k, v := range m.fs {
		if path.Dir(k) == name {
			output = append(output, v)
		}
	}
	return output, nil
}

// Stat returns a FileInfo describing the file.
// If there is an error, it should be of type *PathError.
// This makes Memory compatible with the StatFS interface
func (m *Memory) Stat(name string) (fs.FileInfo, error) {
	f, ok := m.fs[name]
	if !ok {
		return nil, &fs.PathError{
			Op:   "Stat",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}
	return f.Stat()
}
