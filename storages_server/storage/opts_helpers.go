package storage

import (
    "./ifaces"
    "path/filepath"
)

func DefaultStoragesOpts() storage_ifaces.StoragesManagerOpts {
    return storage_ifaces.StoragesManagerOpts{
        DefaultStorageType  : storage_ifaces.StorageMemory,
        Metadata            : "storages.json",
        StoragesRoot        : ".storages",
        DirsMode            : 0o700,
        VaultRoot           : ".vault",
        VaultMode           : 0o600,
        BuffersMode         : 0o600,
        BuffersRoot         : ".buffers",
    }
}


func PrefixedStoragesOpts(prefix storage_ifaces.Path) storage_ifaces.StoragesManagerOpts {
    return storage_ifaces.StoragesManagerOpts{
        DefaultStorageType  : storage_ifaces.StorageMemory,
        DirsMode            : 0o700,
        VaultMode           : 0o600,
        BuffersMode         : 0o600,
        Metadata            : filepath.Join(prefix, "storages.json"),
        VaultRoot           : filepath.Join(prefix, "vault"),
        StoragesRoot        : filepath.Join(prefix, "storages"),
        TempDir             : filepath.Join(prefix, "temp"),
        BuffersRoot         : filepath.Join(prefix, "buffers"),
    }
}
