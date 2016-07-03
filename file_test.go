package atime

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Expected time.
var et = time.Now().Add(-time.Second)

func TestFileReadCloser(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	defer os.Remove(name)

	closeAndSetTimes(t, f, et)

	r, err := NewFileReadCloser(name)
	if err != nil {
		t.Fatal(err)
	}

	expectATimeUpdate(t, r, name, et)
	r.Close()
	expectATimeReset(t, name, et)
}

func TestWithTimesRestored(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	defer os.Remove(name)

	closeAndSetTimes(t, f, et)

	WithTimesRestored(name, func(r io.ReadSeeker) error {
		expectATimeUpdate(t, r, name, et)
		return nil
	})
	expectATimeReset(t, name, et)
}

func closeAndSetTimes(t *testing.T, f *os.File, et time.Time) {
	name := f.Name()
	err := f.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Chtimes(name, et, et)
	if err != nil {
		t.Fatal(err)
	}
}

func expectATimeUpdate(t *testing.T, r io.Reader, name string, et time.Time) error {
	_, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	// Expect updated access time after reading.
	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	at := Get(fi)
	if !at.After(et) {
		t.Errorf("expected atime later than %v, got %v", et, at)
	}
	return nil
}

func expectATimeReset(t *testing.T, name string, et time.Time) {
	fi, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	at := Get(fi)
	if !at.Equal(et) {
		t.Errorf("expected atime %v, got %v", et, at)
	}
}
