package storage_server

import (
    "../storage/ifaces"
    "../storage"
)

type serverContext struct {
    storages *storage.StoragesManager
}

func newContext(opts storage_ifaces.StoragesManagerOpts) *serverContext {
    return &serverContext{
        storages: storage.NewStoragesManager(opts),
    }
}
