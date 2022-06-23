package hashed_filesystem_storage

import (
    "sync"
    "encoding/json"

    "../../ifaces"
)

type asset struct {
    Path    storage_ifaces.Path             `json:"path"`
    Object  string                          `json:"object"`
    Opts    storage_ifaces.StorageAssetOpts `json:"opts"`
}


func (a *asset) VaultAsset() storage_ifaces.VaultAsset {
    return storage_ifaces.VaultAsset{
        Path:   a.Path,
        Object: a.Object,
    }
}


type assetsMap struct {
    sync.Mutex

    values map[storage_ifaces.Path]*asset
}


func makeAssetsMap() assetsMap {
    return assetsMap{
        values: make(map[storage_ifaces.Path]*asset),
    }
}


func (m *assetsMap) Store(path storage_ifaces.Path, asset *asset) {
    m.Lock()
    defer m.Unlock()

    m.values[path] = asset
}


func (m *assetsMap) Load(path storage_ifaces.Path) (*asset, bool) {
    m.Lock()
    defer m.Unlock()

    a, ok := m.values[path]
    if !ok {
        return nil, false
    }

    return a, true
}


func (m *assetsMap) Range(callback func(path storage_ifaces.Path, asset *asset) bool) {
    m.Lock()
    defer m.Unlock()

    for p, a := range m.values {
        if !callback(p, a) {
            break
        }
    }
}


func (m *assetsMap) UnmarshalJSON(b []byte) (err error) {
    m.Lock()
    defer m.Unlock()

    return json.Unmarshal(b, &m.values)
}


func (m assetsMap) MarshalJSON() ([]byte, error) {
    m.Lock()
    defer m.Unlock()

    b, err := json.Marshal(m.values)
    return b, err
}
