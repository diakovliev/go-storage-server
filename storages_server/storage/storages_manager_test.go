package storage

import (
    "testing"
    "./ifaces"

    "io"
    "io/ioutil"
    "bytes"
    "strings"
    _ "io/ioutil"
    "log"
)

const (
    TESTING_WS = ".testing_workspace"
)

func checkStorageOps_NewAsset(s *storage_ifaces.Storage, t *testing.T, path string, payload string, mode int) {

    asset0          := path

    r0, err := s.ReadAsset(asset0)
    if err == nil {
        t.Fatal("Unexpected asset existance!")
    }

    asset0_opts     := storage_ifaces.StorageAssetOpts{Mode: mode}
    asset0_payload  := payload
    asset0_reader   := strings.NewReader(asset0_payload)

    err = s.CreateAsset(asset0, &storage_ifaces.StorageAssetReader{Reader: asset0_reader, Opts: asset0_opts})

    r0, err = s.ReadAsset(asset0)
    if err != nil {
        t.Fatal(err)
    }
    defer r0.Close()

    //b0, err := ioutil.ReadAll(r0)

    var buf bytes.Buffer
    writer := io.Writer(&buf)
    _, err = io.Copy(writer, r0)
    if err != nil {
        t.Fatal(err)
    }

    b0 := buf.Bytes()

    if asset0_opts != r0.Opts {
        t.Fatal("Unexpected asset opts!")
    }

    if asset0_payload != string(b0) {
        t.Fatal("Unexpected asset payload!")
    }

    // Try to create asset one more time
    err = s.CreateAsset(asset0, &storage_ifaces.StorageAssetReader{Reader: asset0_reader, Opts: asset0_opts})
    if err == nil {
        t.Fatal("Unexpected non error!")
    }
}


func checkStorageOps(s *storage_ifaces.Storage, t *testing.T) {

    checkStorageOps_NewAsset(s, t, "asset0", "payload0", 0o666)
    checkStorageOps_NewAsset(s, t, "asset1", "payload1", 0o444)
    checkStorageOps_NewAsset(s, t, "asset2", "payload2", 0o777)
    checkStorageOps_NewAsset(s, t, "asset3", "", 0o777)
    checkStorageOps_NewAsset(s, t, "asset4", "", 0o777)
    checkStorageOps_NewAsset(s, t, "asset5", "", 0o777)
    checkStorageOps_NewAsset(s, t, "asset6", "", 0o777)

    s.Range(func(path storage_ifaces.Path, opts storage_ifaces.StorageAssetOpts) bool {
        log.Printf("Asset: %s mode: %d", path, opts.Mode)
        return true
    })

    resp, err := s.List()
    if err != nil {
        t.Fatal(err)
    }

    log.Printf("List: %s", string(resp))
    //expected_list := `[{"path":"asset0","properties":{"mode":"438"}},{"path":"asset1","properties":{"mode":"292"}},{"path":"asset2","properties":{"mode":"511"}}]`
    //if expected_list != string(resp) {
    //    t.Fatalf("Unexpected List() result!\nGot\t\t: %s\nExpected\t: %s", string(resp), expected_list)
    //}
}


func TestCreateStorage(t *testing.T) {

    opts := PrefixedStoragesOpts(TESTING_WS)

    storagesManager := NewStoragesManager(opts)

    sm := storagesManager.Create(storage_ifaces.StorageDefault)
    if sm == nil {
        t.Fatal("Can't create storage in memory!")
    }

    checkStorageOps(sm, t)

    storagesManager.Destroy(sm.Id)

    spf := storagesManager.Create(storage_ifaces.StoragePlainFilesystem)
    if spf == nil {
        t.Fatal("Can't create plain storage on disk!")
    }

    checkStorageOps(spf, t)

    storagesManager.Destroy(spf.Id)

    shf := storagesManager.Create(storage_ifaces.StorageHashedFilesystem)
    if shf == nil {
        t.Fatal("Can't create hashed storage on disk!")
    }

    checkStorageOps(shf, t)

    storagesManager.Destroy(shf.Id)
}


func TestCreateStorages(t *testing.T) {

    opts := PrefixedStoragesOpts(TESTING_WS)

    storagesManager := NewStoragesManager(opts)
    storagesManager.Create(storage_ifaces.StorageDefault)
    storagesManager.Create(storage_ifaces.StorageDefault)
    storagesManager.Create(storage_ifaces.StorageDefault)

    storagesManager1 := NewStoragesManager(opts)
    storagesManager1.Create(storage_ifaces.StorageDefault)
}


func TestCreateBuffer(t *testing.T) {

    opts := PrefixedStoragesOpts(TESTING_WS)

    storagesManager := NewStoragesManager(opts)

    storage := storagesManager.Create(storage_ifaces.StorageDefault)
    if storage == nil {
        t.Fatal("Can't create storage!")
    }


    bid, err := storagesManager.Buffers().Create()
    if err != nil {
        t.Fatal(err)
    }
    defer storagesManager.Buffers().Discard(bid)

    log.Printf("bid: %s", bid)

    _, err = storagesManager.Buffers().Append(bid, strings.NewReader("123"))
    if err != nil {
        t.Fatal(err)
    }
    _, err = storagesManager.Buffers().Append(bid, strings.NewReader("456"))
    if err != nil {
        t.Fatal(err)
    }
    _, err = storagesManager.Buffers().Append(bid, strings.NewReader("789"))
    if err != nil {
        t.Fatal(err)
    }
    _, err = storagesManager.Buffers().Append(bid, strings.NewReader("0"))
    if err != nil {
        t.Fatal(err)
    }

    err = storagesManager.CreateStorageAssetFromBuffer(storage.Id, "test_1234567890", bid, storage_ifaces.StorageAssetOpts{Mode: 0o666})

    if err != nil {
        t.Fatal("Can't create storage asset!")
    }

    reader, err := storage.ReadAsset("test_1234567890")
    if err != nil {
        t.Fatal(err)
    }

    b, err := ioutil.ReadAll(reader)
    if err != nil {
        t.Fatal(err)
    }

    if string(b) != "1234567890" {
        t.Fatalf("Unexpected buffer content. Got: '%s' expected: '%s'", string(b), "1234567890")
    }
}
