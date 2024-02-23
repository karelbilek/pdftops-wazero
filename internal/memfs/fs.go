// memfs implements a necessary subset of wazero FS to be able to run the conversion.
// It is NOT intended as a general fake memory FS, as it implements only few operations.
package memfs

import (
	"errors"
	"io/fs"

	wasys "github.com/tetratelabs/wazero/sys"

	"os"

	"github.com/tetratelabs/wazero/experimental/sys"

	"github.com/blang/vfs/memfs"
)

func New(inFile []byte) (sys.FS, *memfs.MemFS, error) {
	mfs := memfs.Create()

	f, err := mfs.OpenFile("in.pdf", os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		return nil, nil, err
	}

	_, err = f.Write(inFile)
	if err != nil {
		return nil, nil, err
	}
	mmfs := &memoryFS{fs: mfs}
	return mmfs, mfs, nil
}

type memoryFS struct {
	fs *memfs.MemFS

	sys.UnimplementedFS
}

func toOsOpenFlag(oflag sys.Oflag) (flag int) {
	// First flags are exclusive
	switch oflag & (sys.O_RDONLY | sys.O_RDWR | sys.O_WRONLY) {
	case sys.O_RDONLY:
		flag |= os.O_RDONLY
	case sys.O_RDWR:
		flag |= os.O_RDWR
	case sys.O_WRONLY:
		flag |= os.O_WRONLY
	}

	// Run down the flags defined in the os package
	if oflag&sys.O_APPEND != 0 {
		flag |= os.O_APPEND
	}
	if oflag&sys.O_CREAT != 0 {
		flag |= os.O_CREATE
	}
	if oflag&sys.O_EXCL != 0 {
		flag |= os.O_EXCL
	}
	if oflag&sys.O_SYNC != 0 {
		flag |= os.O_SYNC
	}
	if oflag&sys.O_TRUNC != 0 {
		flag |= os.O_TRUNC
	}
	return flag
}

func (m *memoryFS) OpenFile(path string, flag sys.Oflag, perm fs.FileMode) (sys.File, sys.Errno) {
	f, err := m.fs.OpenFile(path, toOsOpenFlag(flag), perm)
	if err != nil {
		if errors.Is(err, memfs.ErrIsDirectory) {
			dir := &memoryFSDir{fs: m.fs, name: path}
			return dir, 0
		}
		return nil, sys.EINVAL // just general IO error, not that important
	}
	fl := &memoryFSFile{fl: f, name: path, fs: m.fs}
	return fl, 0
}

func (m *memoryFS) Stat(path string) (wasys.Stat_t, sys.Errno) {
	f, errno := m.OpenFile(path, sys.O_RDONLY, 0)
	if errno != 0 {
		return wasys.Stat_t{}, errno
	}
	defer f.Close()
	return f.Stat()
}
