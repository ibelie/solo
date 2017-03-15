// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package solo

import (
	"os"
	"syscall"
)

func Open(filename string, level LockLevel) *File {
	if level == LOCK_PROCESS {
		fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		return &File{File: os.NewFile(uintptr(fd), filename), Level: level}
	} else {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		return &File{File: file, Level: level}
	}
}

func (f *File) fileLock() {
	fd := int(f.Fd())
	if fd != -1 {
		if err := syscall.Flock(fd, syscall.LOCK_EX); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileUnlock() {
	fd := int(f.Fd())
	if fd != -1 {
		if err := syscall.Flock(fd, syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileRLock() {
	fd := int(f.Fd())
	if fd != -1 {
		if err := syscall.Flock(fd, syscall.LOCK_SH); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileRUnlock() {
	fd := int(f.Fd())
	if fd != -1 {
		if err := syscall.Flock(fd, syscall.LOCK_UN); err != nil {
			panic(err)
		}
	}
}
