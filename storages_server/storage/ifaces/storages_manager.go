package storage_ifaces

import (
    "encoding/json"
)

type StoragesManagerOpts struct {
    DefaultStorageType  StorageType

    Metadata            Path

    // PlainFilesystemStorage && HashedFilesystemStorage parameters
    StoragesRoot        Path
    DirsMode            int
    TempDir             Path
    TempPattern         Path

    // HashedFilesystemStorage parameters
    VaultRoot           Path
    VaultMode           int
    VaultDepth          int

    // Buffers manager parameters
    BuffersRoot         Path
    BuffersMode         int
}


func (o StoragesManagerOpts) String() string {
    b, err := json.Marshal(o)
    if err != nil {
        panic(err)
    }
    return string(b)
}


type StoragesManager interface {
    Opts() StoragesManagerOpts
    Vault() Vault
}
