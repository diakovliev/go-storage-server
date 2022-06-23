package storage_ifaces

import (
    "encoding/json"
    "github.com/google/uuid"
)


type StorageId struct {
    Id string `json:"sid"`
}

func MakeNewStorageId() StorageId {
    return StorageId{Id: uuid.New().String()}
}

func MakeStorageId(id string) StorageId {
    return StorageId{Id: id}
}

func (sid *StorageId) String() string {
    return sid.Id
}

func (sid *StorageId) Json() []byte {
    data, err := json.Marshal(sid)
    if err != nil {
        return make([]byte, 0)
    }
    return data
}
