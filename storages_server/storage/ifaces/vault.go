package storage_ifaces

import "os"


type VaultAsset struct {
    Object  string
    Path    Path
}


type VaultFile struct {
    File *os.File
    Asset *VaultAsset
    Owner Vault
}


func (v* VaultFile) Close() error {
    err := v.File.Close()
    v.Owner.CloseObject(v.Asset, v)
    return err
}


type Vault interface {
    Unref(*Storage, VaultAsset)
    Put(*Storage, VaultAsset, Path) error
    OpenObject(VaultAsset) (*VaultFile, error)
    CloseObject(*VaultAsset, *VaultFile)
}
