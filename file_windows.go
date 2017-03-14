// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package solo

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

func Open(filename string, level LockLevel) *File {
	if level == LOCK_PROCESS {
		fd, err := syscall.CreateFile(&(syscall.StringToUTF16(filename)[0]), syscall.GENERIC_READ|syscall.GENERIC_WRITE,
			syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
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

func lockFileEx(h uintptr, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procLockFileEx.Addr(), 6, h, uintptr(flags), 0, 1, 0,
		uintptr(unsafe.Pointer(new(syscall.Overlapped))))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func unlockFileEx(h uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procUnlockFileEx.Addr(), 5, h, 0, 1, 0,
		uintptr(unsafe.Pointer(new(syscall.Overlapped))), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func (f *File) fileLock() {
	fd := f.Fd()
	if fd != 0 {
		if err := lockFileEx(fd, 2); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileUnlock() {
	fd := f.Fd()
	if fd != 0 {
		if err := unlockFileEx(fd); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileRLock() {
	fd := f.Fd()
	if fd != 0 {
		if err := lockFileEx(fd, 0); err != nil {
			panic(err)
		}
	}
}

func (f *File) fileRUnlock() {
	fd := f.Fd()
	if fd != 0 {
		if err := unlockFileEx(fd); err != nil {
			panic(err)
		}
	}
}
