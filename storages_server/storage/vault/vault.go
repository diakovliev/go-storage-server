package vault

import (
    "os"
    "fmt"
    "strings"
    "sync"
    "path/filepath"

    "../ifaces"
    "../filesystem"
)

type Asset = storage_ifaces.VaultAsset


type Vault struct {
    sync.Mutex

    opts    storage_ifaces.StoragesManagerOpts
    Root    storage_ifaces.Path
    Depth   int
    Mode    int

    refs    *Refs

    opened  *Opened
}


func NewVault(opts storage_ifaces.StoragesManagerOpts) *Vault {
    v := &Vault{
        opts    : opts,
        Root    : opts.VaultRoot,
        Depth   : opts.VaultDepth,
        Mode    : opts.VaultMode,
        opened  : NewOpened(),
    }

    vaultLog.Printf("Ensure vault root: %s", v.Root)

    err := filesystem_utils.EnsureDir(v.Root, os.FileMode(opts.DirsMode))
    if err != nil {
        vaultLog.Printf("Ensure vault root error: %s", err)
        return nil
    }

    v.refs = NewRefs(filepath.Join(opts.VaultRoot, "refs.db"), opts)

    return v
}


func (v *Vault) path(h string) string {
    var result strings.Builder

    reader := strings.NewReader(h)

    if v.Depth == 0 {
        v.Depth = 2
    }

    for i := 0; i < v.Depth; i++ {
        ch0, _, err := reader.ReadRune()
        if err != nil {
            panic(err)
        }
        ch1, _, err := reader.ReadRune()
        if err != nil {
            panic(err)
        }

        result.WriteRune(ch0)
        result.WriteRune(ch1)
        result.WriteRune(filepath.Separator)
    }

    for {
        ch, _, err := reader.ReadRune()
        if err != nil {
            break
        }
        result.WriteRune(ch)
    }

    return result.String()
}


func (v *Vault) objectPath(h string) string {
    return filepath.Join(v.Root, v.path(h))
}


func (v *Vault) OpenObject(asset Asset) (*storage_ifaces.VaultFile, error) {
    v.Lock()
    defer v.Unlock()

    vaultLog.Printf("Open object '%s' reader.", asset.Object)

    objectPath := v.objectPath(asset.Object)

    if _, ok := os.Stat(objectPath); os.IsNotExist(ok) {
        return nil, fmt.Errorf("Attempt to open non existing file: %s", objectPath)
    }

    f, err := os.Open(objectPath)
    if err != nil {
        vaultLog.Printf("Open file error: %s", err)
        return nil, err
    }

    v.opened.Open(asset.Object)

    return &storage_ifaces.VaultFile{ File: f, Owner: v, Asset: &asset }, err
}


func (v *Vault) Put(s *storage_ifaces.Storage, asset Asset, filePath storage_ifaces.Path) error {
    v.Lock()
    defer v.Unlock()

    vaultLog.Printf("Put object '%s' to vault.", asset.Object)

    objectPath := v.objectPath(asset.Object)

    if _, ok := os.Stat(objectPath); os.IsNotExist(ok) {

        if err := filesystem_utils.EnsureDir(filepath.Dir(objectPath), os.FileMode(v.opts.DirsMode)); err != nil {
            vaultLog.Printf("Ensure vault root error: %s", err)
            return err
        }

        vaultLog.Printf("Create new object: %s from file: %s", objectPath, filePath)

        if err := os.Rename(filePath, objectPath); err != nil {
            vaultLog.Printf("Rename file error: %s", err)
            return nil
        }

    } else {

        vaultLog.Printf("Remove object source file: %s", filePath)
        os.Remove(filePath)
    }

    // Cancel remove for the object if sheduled
    v.opened.Cancel(asset.Object)

    // Add reference to object
    v.refs.Add(asset.Object, s.Id, asset.Path)

    return nil
}


func (v *Vault) CloseObject(asset *Asset, f *storage_ifaces.VaultFile) {
    v.Lock()
    defer v.Unlock()

    v.opened.Close(asset.Object)

    vaultLog.Printf("Close object '%s' reader.", asset.Object)
}


func (v *Vault) removeObject(h string) {
    removeProc := func() {
        // Lock is not needed because removeProc always called from
        // locked context.

        objectPath := v.objectPath(h)

        vaultLog.Printf("Remove unreferenced object: %s", h)

        os.Remove(objectPath)
    }

    if v.opened.IsOpen(h) {
        v.opened.OnClose(h, removeProc)
    } else {
        removeProc()
    }
}


func (v *Vault) Unref(s *storage_ifaces.Storage, asset Asset) {
    v.Lock()
    defer v.Unlock()

    if refsCount, _ := v.refs.Remove(asset.Object, s.Id, asset.Path); refsCount == 0 {
        // Remove unreferenced object
        v.removeObject(asset.Object)
    }
}
