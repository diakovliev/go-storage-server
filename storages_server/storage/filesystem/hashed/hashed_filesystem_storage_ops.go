package hashed_filesystem_storage

import (
    "../../ifaces"
)

func NewStorageOps(opts storage_ifaces.StoragesManagerOpts) storage_ifaces.StorageOps {
    return &HashedFilesystemStorage{
        assets: makeAssetsMap(),
    }
}
