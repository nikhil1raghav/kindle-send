package epub

import (
	"os"

	"github.com/bmaupin/go-epub/internal/storage"
	"github.com/bmaupin/go-epub/internal/storage/memory"
	"github.com/bmaupin/go-epub/internal/storage/osfs"
)

type FSType int

// filesystem is the current filesytem used as the underlying layer to manage the files.
// See the storage.Use method to change it.
var filesystem storage.Storage = osfs.NewOSFS(os.TempDir())

const (
	// This defines the local filesystem
	OsFS FSType = iota
	// This defines the memory filesystem
	MemoryFS
)

// Use s as default storage/ This is typically used in an init function.
// Default to local filesystem
func Use(s FSType) {
	switch s {
	case OsFS:
		filesystem = osfs.NewOSFS(os.TempDir())
	case MemoryFS:
		//TODO
		filesystem = memory.NewMemory()
	default:
		panic("unexpected FSType")
	}
}
