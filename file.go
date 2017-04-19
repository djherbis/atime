package atime

import (
	"io"
	"os"
	"time"
)

// FileReadCloser satisfies io.Reader, Closer, and Seeker. Stats the
// file before opening it and restores the mtime and atime when
// closing it.
type FileReadCloser struct {
	f            *os.File
	mtime, atime time.Time
}

func NewFileReadCloser(path string) (*FileReadCloser, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FileReadCloser{
		f:     f,
		mtime: fi.ModTime(),
		atime: Get(fi),
	}, nil
}

func (a FileReadCloser) Read(p []byte) (int, error) {
	return a.f.Read(p)
}

func (a FileReadCloser) Seek(offset int64, whence int) (int64, error) {
	return a.f.Seek(offset, whence)
}

func (a FileReadCloser) Close() error {
	path := a.f.Name()
	err := a.f.Close()
	if err != nil {
		return err
	}
	return os.Chtimes(path, a.atime, a.mtime)
}

// WithTimesRestored opens the named file, passes it to a callback,
// and closes it afterward, restoring its atime and mtime.
func WithTimesRestored(path string, fn func(io.ReadSeeker) error) error {
	r, err := NewFileReadCloser(path)
	if err != nil {
		return err
	}
	defer r.Close()
	return fn(r)
}
