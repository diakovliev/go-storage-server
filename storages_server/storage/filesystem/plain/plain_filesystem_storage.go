package plain_filesystem_storage

import (
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "fmt"
    "strings"
    "errors"
    "path/filepath"

    filesystem_utils ".."
    "../../ifaces"
)


type PlainFilesystemStorage struct {
    root storage_ifaces.Path
}


func (pfs *PlainFilesystemStorage) Initialize(s *storage_ifaces.Storage) error {

    pfsLog.Printf("%s: Initialize storage.", s.Name())

    pfs.root = filepath.Join(s.Parent.Opts().StoragesRoot, s.Id.Id)

    pfsLog.Printf("%s: Ensure root: %s", s.Name(), pfs.root)

    err := filesystem_utils.EnsureDir(pfs.root, os.FileMode(s.Parent.Opts().DirsMode))
    if err != nil {
        pfsLog.Printf("%s: Ensure root error: %s", s.Name(), err)
        return err
    }

    return nil
}


func (pfs *PlainFilesystemStorage) Destroy(s *storage_ifaces.Storage) error {

    pfsLog.Printf("%s: Destroy storage.", s.Name())

    pfsLog.Printf("%s: Remove root: %s", s.Name(), pfs.root)

    return os.RemoveAll(pfs.root)
}


func (pfs *PlainFilesystemStorage) CreateAsset(s *storage_ifaces.Storage, path storage_ifaces.Path, r *storage_ifaces.StorageAssetReader) error {

    pfsLog.Printf("%s: Create asset: %s opts: %s", s.Name(), path, r.Opts.String())

    assetPath := filepath.Join(pfs.root, path)

    if _, ok := os.Stat(assetPath); !os.IsNotExist(ok) {
        err := fmt.Errorf("Asset: '%s' already exist!", path)

        pfsLog.Printf("%s: Create asset error: %s", s.Name(), err)

        return err
    }

    err := filesystem_utils.EnsureDir(filepath.Dir(assetPath), os.FileMode(s.Parent.Opts().DirsMode))
    if err != nil {
        pfsLog.Printf("%s: Ensure asset parent dir error: %s", s.Name(), err)
        return err
    }

    f, err := ioutil.TempFile(s.Parent.Opts().TempDir, s.Parent.Opts().TempPattern)
    if err != nil {
        pfsLog.Printf("%s: Can't create temp file! Error: %s", s.Name(), err)
        return nil
    }

    writeProc := func() error {
        defer f.Close()

        writer := bufio.NewWriter(f)

        _, err  = io.Copy(writer, r)
        if err != nil {
            pfsLog.Printf("%s: Create asset copy error: %s", s.Name(), err)
            return err
        }

        err = writer.Flush()
        if err != nil {
            pfsLog.Printf("%s: Create asset flush error: %s", s.Name(), err)
            return err
        }

        err = f.Chmod(os.FileMode(r.Opts.Mode))
        if err != nil {
            pfsLog.Printf("%s: Create asset chmod error: %s", s.Name(), err)
            return err
        }

        return nil
    }

    err = writeProc()
    if err != nil {
        pfsLog.Printf("%s: Remove temp file: %s", s.Name(), f.Name())
        if rmErr := os.Remove(f.Name()); rmErr != nil {
            pfsLog.Panicf("%s: Remove temp file error: %s", s.Name(), rmErr)
        }

        pfsLog.Panicf("%s: Write proc error: %s", s.Name(), err)
    }

    if err := os.Rename(f.Name(), assetPath); err != nil {
        pfsLog.Panicf("%s: File rename error: %s", s.Name(), err)
    }

    return err
}


func (pfs *PlainFilesystemStorage) ReadAsset(s *storage_ifaces.Storage, path storage_ifaces.Path) (*storage_ifaces.StorageAssetReader, error) {

    pfsLog.Printf("%s: Read asset: %s", s.Name(), path)

    assetPath := filepath.Join(pfs.root, path)

    if _, ok := os.Stat(assetPath); os.IsNotExist(ok) {
        return nil, fmt.Errorf("Attempt to read non existing asset: %s", path)
    }

    fi, err := os.Lstat(assetPath)
    if err != nil {
        pfsLog.Printf("%s: Lstat error: %s", s.Name(), err)
        return nil, err
    }

    f, err := os.Open(assetPath)
    if err != nil {
        pfsLog.Printf("%s: Open file error: %s", s.Name(), err)
        return nil, err
    }

    return &storage_ifaces.StorageAssetReader{
        Reader: bufio.NewReader(f),
        Closer: f,
        Opts: storage_ifaces.StorageAssetOpts{Mode: int(fi.Mode().Perm())}}, nil
}


func (pfs *PlainFilesystemStorage) Range(s *storage_ifaces.Storage, callback storage_ifaces.StorageOpsCallback) {
    filepath.Walk(pfs.root, func(path string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }

        assetPath := strings.TrimPrefix(path, pfs.root + string(filepath.Separator))
        if !callback(assetPath, storage_ifaces.StorageAssetOpts{Mode: int(info.Mode().Perm())}) {
            return errors.New("Stop enumeration!")
        }
        return nil
    })
}
