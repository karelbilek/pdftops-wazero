package memfs

import (
	wasys "github.com/tetratelabs/wazero/sys"

	"github.com/tetratelabs/wazero/experimental/sys"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
)

type memoryFSFile struct {
	fs   *memfs.MemFS
	fl   vfs.File
	name string

	sys.UnimplementedFile
}

// used
func (f *memoryFSFile) Stat() (wasys.Stat_t, sys.Errno) {
	fst, err := f.fs.Stat(f.name)
	if err != nil {
		return wasys.Stat_t{}, sys.EBADF
	}
	st := wasys.NewStat_t(fst)

	return st, 0
}

// used
func (f *memoryFSFile) Close() sys.Errno {
	err := f.fl.Close()
	if err != nil {
		return sys.EIO
	}
	return 0
}

// used
func (f *memoryFSFile) IsDir() (bool, sys.Errno) {
	return false, 0
}

// used
func (f *memoryFSFile) Read(buf []byte) (n int, errno sys.Errno) {
	n, err := f.fl.Read(buf)
	if err != nil {
		return 0, sys.EBADF
	}
	return
}

// used
func (f *memoryFSFile) Seek(offset int64, whence int) (newOffset int64, errno sys.Errno) {
	r, err := f.fl.Seek(offset, whence)
	if err != nil {
		return 0, sys.EBADF
	}
	return r, 0
}

// used
func (f *memoryFSFile) Write(buf []byte) (n int, errno sys.Errno) {
	n, err := f.fl.Write(buf)
	if err != nil {
		return 0, sys.EBADF
	}
	return
}
