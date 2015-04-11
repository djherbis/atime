package atime

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	et := time.Now().Add(-time.Second)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(f.Name())
	defer f.Close()

	at, err := Stat(f.Name())
	if err != nil {
		t.Error(err.Error())
	}
	if at.Before(et) {
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at, et)
	}
}

func TestGet(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	et := time.Now().Add(-time.Second)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fi, err := os.Stat(f.Name())
	if err != nil {
		t.Error(err.Error())
	}
	at := Get(fi)
	if at.Before(et) {
		t.Errorf("expected atime to be recent: got %v instead of ~%v", at, et)
	}
}

func TestStatErr(t *testing.T) {
	_, err := Stat("badfile?")
	if err == nil {
		t.Error("expected an error")
	}
}
