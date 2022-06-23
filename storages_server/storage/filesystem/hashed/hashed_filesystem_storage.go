package hashed_filesystem_storage

import (
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "fmt"
    "path/filepath"

    filesystem_utils ".."
    "../../ifaces"
)


type HashedFilesystemStorage struct {
    root        storage_ifaces.Path
    metadata    storage_ifaces.Path
    assets      assetsMap
}


func (hfs *HashedFilesystemStorage) Initialize(s *storage_ifaces.Storage) error {

    hfsLog.Printf("%s: Initialize storage.", s.Name())

    hfs.root = filepath.Join(s.Parent.Opts().StoragesRoot, s.Id.String())

    hfsLog.Printf("%s: Ensure root: %s", s.Name(), hfs.root)

    err := filesystem_utils.EnsureDir(hfs.root, os.FileMode(s.Parent.Opts().DirsMode))
    if err != nil {
        hfsLog.Printf("%s: Ensure root error: %s", s.Name(), err)
        return err
    }

    hfs.metadata = filepath.Join(hfs.root, "metadata.json")

    hfsLog.Printf("%s: Metadata will be in: %s", s.Name(), hfs.metadata)

    hfs.loadMetadata(s)

    return nil
}


func (hfs *HashedFilesystemStorage) Destroy(s *storage_ifaces.Storage) error {

    hfsLog.Printf("%s: Destroy storage.", s.Name())

    hfs.assets.Range(func(path storage_ifaces.Path, asset *asset) bool {

        hfsLog.Printf("%s: Unreference vault object: %s referenced by: %s", s.Name(), asset.Object, asset.Path)

        // Unreference vault object
        s.Parent.Vault().Unref(s, asset.VaultAsset())

        return true
    })

    hfsLog.Printf("%s: Remove root: %s", s.Name(), hfs.root)

    return os.RemoveAll(hfs.root)
}


func (hfs *HashedFilesystemStorage) CreateAsset(s *storage_ifaces.Storage, path storage_ifaces.Path, r *storage_ifaces.StorageAssetReader) error {

    hfsLog.Printf("%s: Create asset: %s opts: %s", s.Name(), path, r.Opts.String())

    if _, ok := hfs.assets.Load(path); ok {

        err := fmt.Errorf("Asset: '%s' already exist!", path)

        hfsLog.Printf("%s: Create asset error: %s", s.Name(), err)

        return err
    }


    f, err := ioutil.TempFile(s.Parent.Opts().TempDir, s.Parent.Opts().TempPattern)
    if err != nil {
        hfsLog.Printf("%s: Can't create temp file! Error: %s", s.Name(), err)
        return err
    }

    fw      := bufio.NewWriter(f)

    writer  := NewCalcChecksumsWriter(fw)

    writeProc := func() error {
        defer f.Close()

        _, err  = io.Copy(writer, r)
        if err != nil {
            hfsLog.Printf("%s: Create asset copy error: %s", s.Name(), err)
            return err
        }

        err = fw.Flush()
        if err != nil {
            hfsLog.Printf("%s: Create asset flush error: %s", s.Name(), err)
            return err
        }

        err = f.Chmod(os.FileMode(s.Parent.Opts().VaultMode))
        if err != nil {
            hfsLog.Printf("%s: Create asset chmod error: %s", s.Name(), err)
            return err
        }

        return nil
    }

    err = writeProc()
    if err != nil {
        hfsLog.Printf("%s: Remove temp file: %s", s.Name(), f.Name())
        if rmErr := os.Remove(f.Name()); rmErr != nil {
            hfsLog.Panicf("%s: Remove temp file error: %s", s.Name(), rmErr)
        }

        hfsLog.Panicf("%s: Write proc error: %s", s.Name(), err)
    }

    asset := &asset{
        Opts:   r.Opts,
        Path:   path,
        Object: writer.String(),
    }

    hfsLog.Printf("%s: Put object to vault as: %s referenced by asset: %s", s.Name(), asset.Object, asset.Path)

    err = s.Parent.Vault().Put(s, asset.VaultAsset(), f.Name())
    if err != nil {
        hfsLog.Printf("%s: Put object to vault error: %s", s.Name())
        return err
    }

    hfsLog.Printf("%s: Register asset: %s", s.Name(), path)
    hfs.assets.Store(path, asset)

    hfs.storeMetadata(s)

    return err
}


func (hfs *HashedFilesystemStorage) ReadAsset(s *storage_ifaces.Storage, path storage_ifaces.Path) (*storage_ifaces.StorageAssetReader, error) {

    hfsLog.Printf("%s: Read asset: %s", s.Name(), path)

    asset, ok := hfs.assets.Load(path)
    if !ok {
        return nil, fmt.Errorf("Attempt to read non existing asset: %s", path)
    }

    f, err := s.Parent.Vault().OpenObject(asset.VaultAsset())
    if err != nil {
        hfsLog.Printf("%s: Open object error: %s", s.Name(), err)
        return nil, err
    }

    return &storage_ifaces.StorageAssetReader{
        Reader: bufio.NewReader(f.File),
        Closer: f,
        Opts: asset.Opts}, nil
}


func (hfs *HashedFilesystemStorage) Range(s *storage_ifaces.Storage, callback storage_ifaces.StorageOpsCallback) {
    hfs.assets.Range(func(path storage_ifaces.Path, asset *asset) bool {
        return callback(path, asset.Opts)
    })
}
