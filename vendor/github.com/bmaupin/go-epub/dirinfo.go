package epub

import "io/fs"

// This is a transitioning function for go < 1.17

// dirInfo is a DirEntry based on a FileInfo.
type dirInfo struct {
	fileInfo fs.FileInfo
}

func (di dirInfo) IsDir() bool {
	return di.fileInfo.IsDir()
}

func (di dirInfo) Type() fs.FileMode {
	return di.fileInfo.Mode().Type()
}

func (di dirInfo) Info() (fs.FileInfo, error) {
	return di.fileInfo, nil
}

func (di dirInfo) Name() string {
	return di.fileInfo.Name()
}

// fileInfoToDirEntry returns a DirEntry that returns information from info.
// If info is nil, FileInfoToDirEntry returns nil.
func fileInfoToDirEntry(info fs.FileInfo) fs.DirEntry {
	if info == nil {
		return nil
	}
	return dirInfo{fileInfo: info}
}
