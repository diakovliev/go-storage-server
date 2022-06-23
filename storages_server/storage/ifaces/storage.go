package storage_ifaces

import "fmt"

type Storage struct {
    Id      StorageId
    Type    StorageType

    Parent  StoragesManager `json:"-"`
    Ops     StorageOps      `json:"-"`
}


func (s Storage) Name() string {
    return fmt.Sprintf("%s(%d)", s.Id.String(), int(s.Type))
}

func (s *Storage) CreateAsset(path Path, r *StorageAssetReader) error {
    return s.Ops.CreateAsset(s, path, r)
}

func (s *Storage) ReadAsset(path Path) (*StorageAssetReader, error) {
    r, err := s.Ops.ReadAsset(s, path)
    return r, err
}

func (s *Storage) Range(callback StorageOpsCallback) {
    s.Ops.Range(s, callback)
}
