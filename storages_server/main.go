package main

import (
    "net/http"

    "github.com/go-chi/chi"
    //"github.com/go-chi/chi/middleware"

    "./storage_server"
    "./storage"
    "./storage/ifaces"
)

func main() {
    http.ListenAndServe(":5555", ApiV1Router())
}

func ApiV1Router() http.Handler {
    router := chi.NewRouter()

    // Brokes HTTP2 streaming! Message buses will be broken!
    //router.Use(middleware.Logger)

    router.Mount("/api/v1", V1Router())
    return router
}


const (
    STANDALONE_WS = ".workspace"
)

func V1Router() http.Handler {
    r := chi.NewRouter()

    opts := storage.PrefixedStoragesOpts(STANDALONE_WS)

    opts.DefaultStorageType = storage_ifaces.StorageHashedFilesystem

    return storage_server.InitializeServer(r, opts, nil)
}
