package storage_ifaces

import (
    "io"
    "fmt"
)

type Path = string


// Asset options.
type StorageAssetOpts struct {
    Mode int    `json:"mode"`
}


func (o StorageAssetOpts) String() string {
    return fmt.Sprintf("{ mode: 0%03o(%d) }", o.Mode, o.Mode)
}


// Helper for reading storage asset. Wraps implementation
// specific 'io.Reader' and 'io.Closer' instances. Also
// provides assess to asset options.
type StorageAssetReader struct {
    io.Reader
    io.Closer
    Opts StorageAssetOpts
}


func (sar *StorageAssetReader) Close() error {
    if sar.Closer == nil {
        return nil
    }
    return sar.Closer.Close()
}


// Callback for the SrorageOps.Range method. If the callback
// returns 'false', enumeration will be stopped.
type StorageOpsCallback = func(Path, StorageAssetOpts) bool


// Storage operations.
type StorageOps interface {

    // Method should initialize Ops instance. The 'nil' must be returned if
    // initialization is sucessed, non 'nil' value means what initializadtion
    // faiiled and Ops instance can't be used.
    Initialize(*Storage) error

    // Destroy Ops instance and release all releated resources. Non 'nil'
    // result means what we have to 'panic'.
    Destroy(*Storage) error

    // Must create storage asset with given options and fill them by data
    // provided by reader. The Path must be an unique per storage. So, if
    // the method will be called several times with the same Path value,
    // all calls, except first, must return errors. Non 'nil' result means
    // what asset is not created.
    CreateAsset(*Storage, Path, *StorageAssetReader) error

    // Open and get reader for existing asset.
    ReadAsset(*Storage, Path) (*StorageAssetReader, error)

    // Enumerates assets existing in storage.
    Range(*Storage, StorageOpsCallback)
}
