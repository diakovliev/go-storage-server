package storage

import (
    "fmt"

    "./ifaces"
    "./memory"
    "./filesystem/plain"
    "./filesystem/hashed"
)

func (sm *StoragesManager) createOps(storageType storage_ifaces.StorageType, opts storage_ifaces.StoragesManagerOpts) (storage_ifaces.StorageOps, storage_ifaces.StorageType) {
    switch storageType {
    case storage_ifaces.StorageDefault:
        if sm.opts.DefaultStorageType == storage_ifaces.StorageDefault {
            panic("The opts.defaultStorageType is set to StorageDefault! Prevent infinite recursion.")
        }
        o, t := sm.createOps(sm.opts.DefaultStorageType, opts)
        return o, t
    case storage_ifaces.StorageMemory:
        return memory_storage.NewStorageOps(opts), storageType
    case storage_ifaces.StoragePlainFilesystem:
        return plain_filesystem_storage.NewStorageOps(opts), storageType
    case storage_ifaces.StorageHashedFilesystem:
       return hashed_filesystem_storage.NewStorageOps(opts), storageType
    default:
        panic(fmt.Errorf("Unknown storage type: %v", storageType))
    }
    return nil, storage_ifaces.StorageDefault
}
