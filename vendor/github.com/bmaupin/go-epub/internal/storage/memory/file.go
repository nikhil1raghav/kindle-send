package memory

import (
	"io"
	"io/fs"
	"time"
)

type file struct {
	name    string
	modTime time.Time
	offset  int
	content []byte
	mode    fs.FileMode
}

func (f *file) Info() (fs.FileInfo, error) {
	return f, nil
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *file) Read(b []byte) (int, error) {
	length := len(f.content)
	start := f.offset
	if start == length {
		return 0, io.EOF
	}
	end := start + len(b)
	if end > length {
		end = length
	}
	f.offset = end
	count := copy(b, f.content[start:end])
	return count, nil
}

func (f *file) Close() error {
	f.offset = 0
	return nil
}

func (f *file) Write(p []byte) (n int, err error) {
	f.content = append(f.content, p...)
	return len(p), nil
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Size() int64 {
	return int64(len(f.content))
}

func (f *file) Type() fs.FileMode {
	return f.mode & fs.ModeType
}

func (f *file) Mode() fs.FileMode {
	return f.mode
}

func (f *file) ModTime() time.Time {
	return f.modTime
}

func (f *file) IsDir() bool {
	return f.mode&fs.ModeDir != 0
}

func (f *file) Sys() interface{} {
	return nil
}
