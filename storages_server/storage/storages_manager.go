package storage

import (
    "os"
    "fmt"

    "./ifaces"
    "./filesystem"
    "./vault"
    "./buffers"
)


type StoragesManager struct {
    opts        storage_ifaces.StoragesManagerOpts
    storages    storagesMap
    vault       *vault.Vault
    buffers     *buffers.BuffersManager
}


func NewStoragesManager(opts storage_ifaces.StoragesManagerOpts) *StoragesManager {

    storagesLog.Printf("Create storages manager. Opts: %s", opts.String())

    sm := &StoragesManager{
        opts        : opts,
        storages    : makeStoragesMap(),
        vault       : vault.NewVault(opts),
        buffers     : buffers.NewBuffersManager(storage_ifaces.BuffersManagerOpts{
            StorageRoot     : opts.BuffersRoot,
            StorageRootMode : opts.DirsMode,
            FilesMode       : opts.BuffersMode,
        }),
    }

    if len(sm.opts.TempDir) > 0 {
        if err := filesystem_utils.EnsureDir(sm.opts.TempDir, os.FileMode(sm.opts.DirsMode)); err != nil {
            storagesLog.Panicf("Ensure temp dir error: %s", err)
        }
    }

    sm.reattachToStorages()

    return sm
}


func (sm *StoragesManager) Opts() storage_ifaces.StoragesManagerOpts {
    return sm.opts
}


func (sm *StoragesManager) Vault() storage_ifaces.Vault {
    return sm.vault
}


func (sm *StoragesManager) Buffers() storage_ifaces.BuffersManager {
    return sm.buffers
}


func (sm *StoragesManager) Get(id storage_ifaces.StorageId) *storage_ifaces.Storage {
    s, ok := sm.storages.Load(id)
    if !ok {
        storagesLog.Printf("Can't find storage by id: %s", id.Id)
        return nil
    }

    return s
}


func (sm *StoragesManager) reattachToStorages() {

    // Init storages from metadata
    sm.loadMetadata()

    toRemove := make([]storage_ifaces.StorageId, 0, 100)

    // Create runtime objects
    sm.storages.Range(func(id storage_ifaces.StorageId, s *storage_ifaces.Storage) bool {

        if s.Type == storage_ifaces.StorageMemory {
            storagesLog.Printf("No sence in reattaching to memory storage: %s", id.Id)
            toRemove = append(toRemove, id)
            return true
        }

        storagesLog.Printf("Attach to existing storage: %s", id.Id)

        s.Parent        = sm
        s.Ops, s.Type   = sm.createOps(s.Type, sm.Opts())

        if err := s.Ops.Initialize(s); err != nil {
            storagesLog.Panicf("Can't initialize storage! Error: %s", err)
        }

        return true
    })

    for _, id := range toRemove {
        storagesLog.Printf("Delete memory storage: %s", id.Id)
        sm.storages.Delete(id)
    }

    if len(toRemove) > 0 {
        sm.storeMetadata()
    }
}


func (sm *StoragesManager) Create(storageType storage_ifaces.StorageType) *storage_ifaces.Storage {

    ops, sType := sm.createOps(storageType, sm.Opts())

    if ops == nil {
        storagesLog.Panicf("Can't create storage ops for storage type: %d.", int(storageType))
    }

    s := &storage_ifaces.Storage{
        Parent  : sm,
        Type    : sType,
        Ops     : ops,
        Id      : storage_ifaces.MakeNewStorageId(),
    }

    storagesLog.Printf("Created new storage: %s", s.Name())

    if err := s.Ops.Initialize(s); err != nil {
        storagesLog.Panicf("Storage initialization error: %s", err)
    }

    sm.storages.Store(s.Id, s)

    sm.storeMetadata()

    return s
}


func (sm *StoragesManager) Destroy(storageId storage_ifaces.StorageId) error {

    storage, ok := sm.storages.Load(storageId)
    if !ok {
        return fmt.Errorf("Attempt to destroy non existing storage: %s!", storageId.Id)
    }

    storagesLog.Printf("Destroy storage: %s", storage.Name())

    err := storage.Ops.Destroy(storage)
    if err != nil {
        return err
    }

    sm.storages.Delete(storageId)

    sm.storeMetadata()

    return nil
}


func (sm *StoragesManager) CreateStorageAssetFromBuffer(
    storageId storage_ifaces.StorageId,
    path storage_ifaces.Path,
    bufferId string,
    opts storage_ifaces.StorageAssetOpts) error {

    storage, ok := sm.storages.Load(storageId)
    if !ok {
        return fmt.Errorf("Attempt to use non existing storage: %s!", storageId.Id)
    }

    if err := sm.Buffers().EnsureBuffer(bufferId); err != nil {
        return fmt.Errorf("Attempt to use non existing buffer { id: %s }!", bufferId)
    }

    f, err := os.OpenFile(sm.Buffers().Abspath(bufferId), os.O_RDONLY, os.FileMode(0))
    if err != nil {
        return fmt.Errorf("Can't open buffer { id: %s }", bufferId)
    }
    defer f.Close()

    return storage.CreateAsset(path, &storage_ifaces.StorageAssetReader{Reader: f, Opts: opts})
}
