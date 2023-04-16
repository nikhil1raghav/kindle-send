package types

type FileType string

var (
	TypeUrl     FileType = "url"
	TypeUrlFile FileType = "urlfile"
	TypeFile    FileType = "file"
)

type Request struct {
	Path    string
	Type    FileType
	Options map[string]string
}

func NewRequest(path string, fileType FileType, opts map[string]string) Request {
	return Request{path, fileType, opts}
}
