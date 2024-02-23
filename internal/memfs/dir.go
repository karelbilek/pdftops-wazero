package memfs

import (
	wasys "github.com/tetratelabs/wazero/sys"

	"github.com/tetratelabs/wazero/experimental/sys"

	"github.com/blang/vfs/memfs"
)

type memoryFSDir struct {
	fs   *memfs.MemFS
	name string

	sys.UnimplementedFile
}

func (f *memoryFSDir) IsDir() (bool, sys.Errno) {
	return true, 0
}

func (f *memoryFSDir) Stat() (wasys.Stat_t, sys.Errno) {
	fst, err := f.fs.Stat(f.name)
	if err != nil {
		return wasys.Stat_t{}, sys.EBADF
	}
	st := wasys.NewStat_t(fst)

	return st, 0
}
