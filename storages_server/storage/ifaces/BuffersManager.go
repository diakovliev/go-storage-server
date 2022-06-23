package storage_ifaces

import (
    "io"
    "fmt"
)


type BuffersManagerOpts struct {
    StorageRootMode int
    StorageRoot     Path

    FilesMode       int
}


func (opts BuffersManagerOpts) String() string {
    return fmt.Sprintf("{ StorageRootMode: %d, StorageRoot: %s, FilesMode: %d }",
        opts.StorageRootMode, opts.StorageRoot, opts.FilesMode)
}


type BuffersManager interface {
    Abspath(bid string) string
    Create() (string, error)
    Discard(bid string) error
    EnsureBuffer(bid string) error
    Append(bid string, source io.Reader) (int, error)
}
