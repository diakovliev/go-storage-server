package storage_server

import (
    "net/http"
    "net/url"
    "strings"
    "errors"

    "github.com/go-chi/chi"

    "../storage/ifaces"
    "../storage"
)

var (
    context *serverContext
)


func GetStoragesManager() (sm *storage.StoragesManager, err error) {
    if context == nil {
        sm  = nil
        err = errors.New("Not initalized storage server!")
        return
    }

    sm  = context.storages
    err = nil
    return
}


// TODO: Integrate with corvusd/auth
func InitializeServer(r *chi.Mux, opts storage_ifaces.StoragesManagerOpts, auth interface{}) *chi.Mux {

    context = newContext(opts)

    r.Route("/storage", func(r chi.Router) {

        // Pseudo FS level routines
        r.Get("/create/{type}", StorageCreate)
        r.Get("/destroy/{sid:[0-f-]+}", StorageDestroy)
        r.Get("/list/{sid:[0-f-]+}", StorageList)
        r.Route("/{sid:[0-f-]+}", func(r chi.Router) {
            r.Put("/*", StoragePutElement)
            r.Get("/*", StorageGetElement)
        })

        // Input buffers row
        r.Route("/buffer", func(r chi.Router) {
            r.Get("/create", BufferCreate)
            r.Get("/discard/{bid:[0-f-]+}", BufferDiscard)
            r.Get("/commit/{sid:[0-f-]+}/{bid:[0-f-]+}/*", BufferCommit)
            r.Put("/{bid:[0-f-]+}", BufferAppend)
        })

    })

    return r
}


func jsonResponse(w http.ResponseWriter, resp []byte) {
    w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}

type properties map[string]string

func makeProperties() properties {
    return make(properties)
}

func getProperties(args url.Values) properties {
    props := makeProperties()
    for k := range args {
        if len(args[k]) > 0 {
            props[k] = args[k][0]
        }
    }

    return props
}


func extractPath(urlPath string, sepa string) (string, bool) {
    parts := strings.Split(urlPath, "/" + sepa + "/")
    if len(parts) < 2 {
        return urlPath, false
    }

    return parts[1], true
}
