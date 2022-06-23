package hashed_filesystem_storage

import (
    "os"
    "io/ioutil"
    "bufio"
    "encoding/json"

    "../../ifaces"
)


func (hfs *HashedFilesystemStorage) loadMetadata(s *storage_ifaces.Storage) {

    hfsLog.Printf("%s: Load metadata from: %s", s.Name(), hfs.metadata)

    if _, ok := os.Stat(hfs.metadata); os.IsNotExist(ok) {
        hfsLog.Printf("%s: No file: %s. Skip loading storages metadata.", s.Name(), hfs.metadata)
        return
    }

    f, err := os.Open(hfs.metadata)
    if err != nil {
        hfsLog.Panicf("File open error: %s", err)
    }
    defer f.Close()

    err = json.NewDecoder(f).Decode(&hfs.assets)
    if err != nil {
        hfsLog.Panicf("Decoder error: %s", err)
    }
}


func (hfs *HashedFilesystemStorage) storeMetadata(s *storage_ifaces.Storage) {

    hfsLog.Printf("%s: Store metadata to: %s", s.Name(), hfs.metadata)

    f, err := ioutil.TempFile(s.Parent.Opts().TempDir, s.Parent.Opts().TempPattern)
    if err != nil {
        hfsLog.Panicf("Can't create temp file! Error: %s", err)
    }

    writer := bufio.NewWriter(f)

    writeProc := func() error {
        defer f.Close()

        err := json.NewEncoder(writer).Encode(hfs.assets)
        if err != nil {
            hfsLog.Printf("%s: Encoder error: %s", s.Name(), err)
            return err
        }

        err = writer.Flush()
        if err != nil {
            hfsLog.Printf("%s: Writer flush error: %s", s.Name(), err)
            return err
        }

        err = f.Chmod(os.FileMode(s.Parent.Opts().VaultMode))
        if err != nil {
            hfsLog.Printf("%s: File chmod error: %s", s.Name(), err)
            return err
        }

        return nil
    }

    err = writeProc()
    if err != nil {
        hfsLog.Printf("%s: Remove temp file: %s", s.Name(), f.Name())

        if rmErr := os.Remove(f.Name()); rmErr != nil {
            hfsLog.Panicf("Remove temp file error: %s", rmErr)
        }

        hfsLog.Panicf("Write proc error: %s", err)
    }

    err = os.Rename(f.Name(), hfs.metadata)
    if err != nil {
        hfsLog.Panicf("Rename error: %s", err)
    }
}
