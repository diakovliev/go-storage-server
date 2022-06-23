package storage_server

import (
    "log"
    "fmt"
    "strconv"
    "net/http"

    "github.com/go-chi/chi"

    "../storage/ifaces"
)


func BufferCreate(w http.ResponseWriter, r *http.Request) {

    bid, err := context.storages.Buffers().Create()
    if err != nil {
        log.Printf("Create buffer error: %w", err)
        http.Error(w, fmt.Sprintf("Create buffer error: %w", err), http.StatusInternalServerError)
        return
    }

    jbid := storage_ifaces.MakeStorageId(bid)

    jsonResponse(w, jbid.Json())
}


func BufferDiscard(w http.ResponseWriter, r *http.Request) {
    bid := chi.URLParam(r, "bid")
    if len(bid) < 1 {
        http.Error(w, "Empty buffer id!", http.StatusNotFound)
        return
    }

    if err := context.storages.Buffers().Discard(bid); err != nil {
        log.Printf("Discard buffer error: %w", err)
        http.Error(w, fmt.Sprintf("Discard buffer error: %w", err), http.StatusInternalServerError)
        return
    }
}


func BufferCommit(w http.ResponseWriter, r *http.Request) {
    sid := chi.URLParam(r, "sid")
    if len(sid) < 1 {
        http.Error(w, "Empty storage id!", http.StatusNotFound)
        return
    }

    bid := chi.URLParam(r, "bid")
    if len(bid) < 1 {
        http.Error(w, "Empty buffer id!", http.StatusNotFound)
        return
    }

    log.Printf("PUT url path: %s", r.URL.Path)

    path, ok := extractPath(r.URL.Path, bid)
    if !ok {
        log.Printf("Empty path!")
        http.Error(w, "Empty path!", http.StatusNotFound)
        return
    }

    id := storage_ifaces.MakeStorageId(sid)

    s := context.storages.Get(id)
    if s == nil {
        log.Printf("Unknown storage id!")
        http.Error(w, "Unknown storage id!", http.StatusNotFound)
        return
    }

    mode := 0o644
    modeStr := getProperties(r.URL.Query())["mode"]
    log.Printf("mode str: %s", modeStr)
    if len(modeStr) != 0 {
        modeVal, err := strconv.ParseInt(modeStr, 0, 32)
        if err != nil {
            log.Printf("Permission conversion error: %w. File: %s", err, path)
            http.Error(w, fmt.Sprintf("Permission conversion error: %w. File: %s", err, path), http.StatusInternalServerError)
            return
        }
        mode = int(modeVal)
    }

    if err := context.storages.CreateStorageAssetFromBuffer(id, path, bid, storage_ifaces.StorageAssetOpts{Mode: mode}); err != nil {
        log.Printf("Create storage asset error: %w. File: %s", err, path)
        http.Error(w, fmt.Sprintf("Create storage asset error: %w. File: %s", err, path), http.StatusInternalServerError)
        return
    }
}


func BufferAppend(w http.ResponseWriter, r *http.Request) {
    bid := chi.URLParam(r, "bid")
    if len(bid) < 1 {
        http.Error(w, "Empty buffer id!", http.StatusNotFound)
        return
    }

    if _, err := context.storages.Buffers().Append(bid, r.Body); err != nil {
        log.Printf("Append buffer error: %w", err)
        http.Error(w, fmt.Sprintf("Append buffer error: %w", err), http.StatusInternalServerError)
        return
    }
}
