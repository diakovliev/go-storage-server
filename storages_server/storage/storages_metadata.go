package storage

import (
    "os"
    "io/ioutil"
    "bufio"
    "log"
    "encoding/json"
)


func (sm *StoragesManager) loadMetadata() {

    log.Printf("Load storages metadata from: %s", sm.opts.Metadata)

    if _, ok := os.Stat(sm.opts.Metadata); os.IsNotExist(ok) {
        log.Printf("No file: %s. Skip loading storages metadata.", sm.opts.Metadata)
        return
    }

    f, err := os.Open(sm.opts.Metadata)
    if err != nil {
        log.Panicf("Open file error: %s", err)
    }
    defer f.Close()

    decoder := json.NewDecoder(f)

    err = decoder.Decode(&sm.storages)
    if err != nil {
        log.Panicf("Storages decode error: %s", err)
    }
}


func (sm *StoragesManager) storeMetadata() {

    log.Printf("Store storages metadata to: %s", sm.opts.Metadata)

    f, err := ioutil.TempFile(sm.opts.TempDir, sm.opts.TempPattern)
    if err != nil {
        log.Fatalf("Can't create temp file! Error: %s", err)
    }

    writer := bufio.NewWriter(f)

    writeProc := func() error {
        defer f.Close()

        encoder := json.NewEncoder(writer)

        err := encoder.Encode(sm.storages)
        if err != nil {
            log.Printf("Storages encoding error: %s", err)
            return err
        }

        err = writer.Flush()
        if err != nil {
            log.Printf("Writer flush error: %s", err)
            return err
        }

        err = f.Chmod(os.FileMode(sm.opts.VaultMode))
        if err != nil {
            log.Printf("Chmod error: %s", err)
            return err
        }

        return nil
    }

    err = writeProc()
    if err != nil {
        log.Printf("Remove temp file: %s", f.Name())
        if rmErr := os.Remove(f.Name()); rmErr != nil {
            log.Printf("Remove temp file error: %s", rmErr)
        }

        log.Panicf("Write proc error: %s", err)
    }

    err = os.Rename(f.Name(), sm.opts.Metadata)
    if err != nil {
        log.Panicf("Rename error: %s", err)
    }
}
