package storage_ifaces

type StorageType int

const (
    StorageDefault StorageType = iota
    StorageMemory
    StoragePlainFilesystem
    StorageHashedFilesystem
)

func StorageType_fromString(input string) StorageType {
    switch input {
    case "default":
        return StorageDefault
    case "memory":
        return StorageMemory
    case "filesystem":
        return StoragePlainFilesystem
    case "hashed":
        return StorageHashedFilesystem
    }
    return StorageDefault
}

func StorageType_toString(st StorageType) string {
    switch st {
    case StorageDefault:
        return "default"
    case StorageMemory:
        return "memory"
    case StoragePlainFilesystem:
        return "filesystem"
    case StorageHashedFilesystem:
        return "hashed"
    }
    return "default"
}
