// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package solo

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLockUnlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "x")
	file := Open(path, LOCK_PROCESS)
	defer file.Close()

	file.Lock(nil)
	file.WriteString("TestLockUnlock")
	file.Unlock()
}

func TestRLockUnlock(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "x")
	file := Open(path, LOCK_PROCESS)
	defer file.Close()

	file.Lock(nil)
	file.WriteString("TestRLockUnlock")
	file.Unlock()

	b := make([]byte, 15)
	file.RLock(nil)
	file.ReadAt(b, 0)
	if !bytes.Equal([]byte("TestRLockUnlock"), b) {
		t.Errorf("TestRLockUnlock not equal: %v", b)
	}
	file.RUnlock()
}

func TestSimultaneousLock(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "x")
	file := Open(path, LOCK_PROCESS)
	defer file.Close()

	file.Lock(nil)

	state := "waiting"
	ch := make(chan struct{})
	go func() {
		file.Lock(nil)
		state = "acquired"
		ch <- struct{}{}

		<-ch
		file.Unlock()
		state = "released"
		ch <- struct{}{}
	}()

	if "waiting" != state {
		t.Errorf("TestSimultaneousLock state error: %v", state)
	}
	file.Unlock()

	<-ch
	if "acquired" != state {
		t.Errorf("TestSimultaneousLock state error: %v", state)
	}
	ch <- struct{}{}

	<-ch
	if "released" != state {
		t.Errorf("TestSimultaneousLock state error: %v", state)
	}
}
