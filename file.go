// Copyright 2017 ibelie, Chen Jie, Joungtao. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package solo

import (
	"os"
	"sync"
)

type LockLevel int

const (
	LOCK_NOTHING LockLevel = iota
	LOCK_GOROUTINE
	LOCK_PROCESS
)

type File struct {
	*os.File
	Level LockLevel
	mutex sync.RWMutex
	once  sync.Once
}

func (f *File) Lock(verify func()) *File {
	if f.Level == LOCK_NOTHING && verify != nil {
		f.once.Do(verify)
	} else if f.Level == LOCK_GOROUTINE {
		f.mutex.Lock()
		if verify != nil {
			verify()
		}
	} else if f.Level == LOCK_PROCESS {
		f.mutex.Lock()
		f.fileLock()
		if verify != nil {
			verify()
		}
	}
	return f
}

func (f *File) Unlock() {
	if f.Level == LOCK_GOROUTINE {
		f.mutex.Unlock()
	} else if f.Level == LOCK_PROCESS {
		f.fileUnlock()
		f.mutex.Unlock()
	}
}

func (f *File) RLock(verify func()) *File {
	if f.Level == LOCK_NOTHING && verify != nil {
		f.once.Do(verify)
	} else if f.Level == LOCK_GOROUTINE {
		f.mutex.RLock()
		if verify != nil {
			verify()
		}
	} else if f.Level == LOCK_PROCESS {
		f.mutex.RLock()
		f.fileRLock()
		if verify != nil {
			verify()
		}
	}
	return f
}

func (f *File) RUnlock() {
	if f.Level == LOCK_GOROUTINE {
		f.mutex.RUnlock()
	} else if f.Level == LOCK_PROCESS {
		f.fileRUnlock()
		f.mutex.RUnlock()
	}
}
