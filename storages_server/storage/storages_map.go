package storage

import (
    "sync"
    "encoding/json"

    "./ifaces"
)


type storagesMap struct {
    sync.Mutex

    values map[string]*storage_ifaces.Storage
}


func makeStoragesMap() storagesMap {
    return storagesMap{
        values: make(map[string]*storage_ifaces.Storage),
    }
}


func (m *storagesMap) Store(id storage_ifaces.StorageId, s *storage_ifaces.Storage) {
    m.Lock()
    defer m.Unlock()

    m.values[id.String()] = s
}


func (m *storagesMap) Load(id storage_ifaces.StorageId) (*storage_ifaces.Storage, bool) {
    m.Lock()
    defer m.Unlock()

    s, ok := m.values[id.String()]
    if !ok {
        return nil, false
    }

    return s, true
}


func (m *storagesMap) Range(callback func(id storage_ifaces.StorageId, storage *storage_ifaces.Storage) bool) {
    m.Lock()
    defer m.Unlock()

    for i, s := range m.values {
        if !callback(storage_ifaces.MakeStorageId(i), s) {
            break
        }
    }
}


func (m *storagesMap) Delete(id storage_ifaces.StorageId) {
    m.Lock()
    defer m.Unlock()

    delete(m.values, id.String())
}


func (m *storagesMap) UnmarshalJSON(b []byte) (err error) {
    m.Lock()
    defer m.Unlock()

    return json.Unmarshal(b, &m.values)
}


func (m storagesMap) MarshalJSON() ([]byte, error) {
    m.Lock()
    defer m.Unlock()

    b, err := json.Marshal(m.values)
    return b, err
}
