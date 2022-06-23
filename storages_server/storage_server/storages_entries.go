package storage_server

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strconv"
    "fmt"

    "github.com/go-chi/chi"

    "../storage/ifaces"
)


func StorageCreate(w http.ResponseWriter, r *http.Request) {
    storageTypeString := chi.URLParam(r, "type")
    if len(storageTypeString) < 1 {
        http.Error(w, "Empty storage type!", http.StatusNotFound)
        return
    }

    st := storage_ifaces.StorageType_fromString(storageTypeString)

    log.Printf("Create ne '%s' storage.", storage_ifaces.StorageType_toString(st))

    s := context.storages.Create(st)
    if s == nil {
        http.Error(w, "Failed to create storage!", http.StatusInternalServerError)
        return
    }

    log.Printf("Created new storage with id: %s", s.Id.String())

    jsonResponse(w, s.Id.Json())
}


func StorageDestroy(w http.ResponseWriter, r *http.Request) {
    sid := chi.URLParam(r, "sid")
    if len(sid) < 1 {
        http.Error(w, "Empty storage id!", http.StatusNotFound)
        return
    }

    id := storage_ifaces.MakeStorageId(sid)

    if s := context.storages.Get(id); s == nil {
        http.Error(w, "Unknown storage id!", http.StatusNotFound)
        return
    }

    if err := context.storages.Destroy(id); err != nil {
        http.Error(w, "Failed to delete storage!", http.StatusInternalServerError)
        return
    }

    log.Printf("Deleted storage with id: %s", id.String())

    resp, err := json.Marshal(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonResponse(w, resp)
}


func StorageList(w http.ResponseWriter, r *http.Request) {
    sid := chi.URLParam(r, "sid")
    if len(sid) < 1 {
        http.Error(w, "Empty storage id!", http.StatusNotFound)
        return
    }

    id := storage_ifaces.MakeStorageId(sid)

    s := context.storages.Get(id)
    if s == nil {
        http.Error(w, "Unknown storage id!", http.StatusNotFound)
        return
    }

    resp, err := s.List()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    jsonResponse(w, resp)
}


func StoragePutElement(w http.ResponseWriter, r *http.Request) {
    sid := chi.URLParam(r, "sid")
    if len(sid) < 1 {
        http.Error(w, "Empty storage id!", http.StatusNotFound)
        return
    }

    log.Printf("PUT url path: %s", r.URL.Path)

    path, ok := extractPath(r.URL.Path, sid)
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


    err := s.CreateAsset(path, &storage_ifaces.StorageAssetReader{
        Reader: r.Body,
        Opts: storage_ifaces.StorageAssetOpts{Mode: mode},
    })
    if err != nil {
        log.Printf("Error on creating storage element: %w. File: %s", err, path)
        http.Error(w, "Error on creating storage element!", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}


func StorageGetElement(w http.ResponseWriter, r *http.Request) {
    sid := chi.URLParam(r, "sid")
    if len(sid) < 1 {
        http.Error(w, "Empty storage id!", http.StatusNotFound)
        return
    }

    log.Printf("GET url path: %s", r.URL.Path)

    path, ok := extractPath(r.URL.Path, sid)
    if !ok {
        http.Error(w, "Empty path!", http.StatusNotFound)
        return
    }

    id := storage_ifaces.MakeStorageId(sid)

    s := context.storages.Get(id)
    if s == nil {
        http.Error(w, "Unknown storage id!", http.StatusNotFound)
        return
    }

    reader, err := s.ReadAsset(path)
    if err != nil {
        http.Error(w, fmt.Sprintf("%w", err), http.StatusNotFound)
        return
    }
    defer reader.Close()

    _, err = io.Copy(w, reader)
    if err != nil {
        http.Error(w, "Error on reading storage element!", http.StatusInternalServerError)
        return
    }

    log.Printf("get sid: %s path: %s", sid, path)
}
