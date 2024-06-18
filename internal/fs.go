package internal

import (
	"io/fs"
	"os"
	"strings"
	"time"
)

type osFS struct {
	// filename -> content
	overlay map[string]string
}

func NewOSFS() fs.FS {
	return &osFS{}
}

func NewOverlayFS(overlay map[string]string) fs.FS {
	return &osFS{
		overlay: overlay,
	}
}

func (f *osFS) Open(name string) (fs.File, error) {
	if content, present := f.overlay[name]; present {
		return NewStringFile(content), nil
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewStringFile(string(data)), err
}

type StringFile struct {
	data *strings.Reader
}

func NewStringFile(content string) *StringFile {
	return &StringFile{
		data: strings.NewReader(content),
	}
}

// Read implements the io.Reader interface
func (f *StringFile) Read(p []byte) (int, error) {
	return f.data.Read(p)
}

// Stat implements the fs.File interface
func (f *StringFile) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		size: int64(f.data.Len()),
	}, nil
}

// Close implements the fs.File interface
func (f *StringFile) Close() error {
	return nil
}

// fileInfo implements fs.FileInfo
type fileInfo struct {
	size int64
}

func (fi *fileInfo) Name() string       { return "string file" }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() fs.FileMode  { return 0 }
func (fi *fileInfo) ModTime() time.Time { return time.Time{} }
func (fi *fileInfo) IsDir() bool        { return false }
func (fi *fileInfo) Sys() interface{}   { return nil }
