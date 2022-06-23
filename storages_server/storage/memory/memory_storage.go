package memory_storage

import (
    "io"
    "bytes"
    "fmt"
    "sync"

    "../ifaces"
)


type MemoryStorage struct {
    assets sync.Map
}


type asset struct {
    path    storage_ifaces.Path
    opts    storage_ifaces.StorageAssetOpts
    payload []byte
}


func (ms *MemoryStorage) Initialize(s *storage_ifaces.Storage) error {

    memoryLog.Printf("%s: Initialize storage.", s.Name())

    return nil
}


func (ms *MemoryStorage) Destroy(s *storage_ifaces.Storage) error {

    memoryLog.Printf("%s: Destroy storage.", s.Name())

    return nil
}


func (ms *MemoryStorage) CreateAsset(s *storage_ifaces.Storage, path storage_ifaces.Path, r *storage_ifaces.StorageAssetReader) error {

    memoryLog.Printf("%s: Create asset: %s opts: %s", s.Name(), path, r.Opts.String())

    if _, ok := ms.assets.Load(path); ok {
        return fmt.Errorf("Asset: '%s' already exist!", path)
    }

    asset := &asset {
        path: path,
        opts: r.Opts,
    }

    var data bytes.Buffer

    w := io.Writer(&data)

    _, err := io.Copy(w, r)
    if err != nil {
        memoryLog.Printf("%s: Create asset copy error: %s", s.Name(), err)
        return err
    }

    asset.payload = data.Bytes()

    ms.assets.Store(path, asset)

    return nil
}


func (ms *MemoryStorage) ReadAsset(s *storage_ifaces.Storage, path storage_ifaces.Path) (*storage_ifaces.StorageAssetReader, error) {

    memoryLog.Printf("%s: Read asset: %s", s.Name(), path)

    iasset, ok := ms.assets.Load(path)
    if !ok {
        return nil, fmt.Errorf("Attempt to read non existing asset: %s", path)
    }

    asset, ok := iasset.(*asset)
    if !ok {
        memoryLog.Panic("Unexpected value in ms.assets!")
    }

    return &storage_ifaces.StorageAssetReader{
        Reader: bytes.NewReader(asset.payload),
        Closer: nil,
        Opts: asset.opts}, nil
}


func (ms *MemoryStorage) Range(s *storage_ifaces.Storage, callback storage_ifaces.StorageOpsCallback) {
    ms.assets.Range(func(key, value interface{}) bool {
        path, ok := key.(storage_ifaces.Path)
        if !ok {
            memoryLog.Panic("Unexpected key type!")
        }

        asset, ok := value.(*asset)
        if !ok {
            memoryLog.Panic("Unexpected value type!")
        }

        return callback(path, asset.opts)
    })
}
