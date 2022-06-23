package vault

import (
    "os"
    "fmt"
    "sync"
    "io/ioutil"
    "bufio"
    "encoding/json"
    "errors"
    "../ifaces"
)

// Reference to storage object descriptor.
type Ref struct {

    // Storage id.
    StorageId   storage_ifaces.StorageId

    // Storage path.
    Path        storage_ifaces.Path
}


// Type alias.
type RefsSlice = []*Ref


// References collection (refs database).
type Refs struct {
    sync.Mutex

    opts    storage_ifaces.StoragesManagerOpts

    //  Location of the file in what refs database will be
    // serialized.
    dbpath  storage_ifaces.Path

    //  References database data. This is a map where key
    // is the referenced object id, value is a slice with
    // references descriptors.
    values  map[string]RefsSlice
}


// Declare errors
var (
    REF_EXIST       = errors.New("Reference already exist!")
    REF_NOT_EXIST   = errors.New("Reference not exist!")
)


// Make new references database instance.
func NewRefs(dbpath storage_ifaces.Path, opts storage_ifaces.StoragesManagerOpts) *Refs {
    r := &Refs{
        opts    : opts,
        dbpath  : dbpath,
        values  : make(map[string]RefsSlice),
    }

    // Load refs database from file if the file exist.
    return r.loadDb();
}


//  Load references database from file located by r.dbpath location,
// if a file exists. If the file not exists - do nothing.
// Returns pointer to struct itself. Not thread safe.
func (r *Refs) loadDb() *Refs {
    vaultLog.Printf("Load references database from: %s", r.dbpath)

    if _, ok := os.Stat(r.dbpath); os.IsNotExist(ok) {
        vaultLog.Printf("No file: %s. Skip loading references database.", r.dbpath)
        return r
    }

    f, err := os.Open(r.dbpath)
    if err != nil {
        vaultLog.Panicf("Open file error: %s", err)
    }
    defer f.Close()

    decoder := json.NewDecoder(f)

    err = decoder.Decode(r)
    if err != nil {
        vaultLog.Panicf("Storages decode error: %s", err)
    }

    return r
}


// Store references database to a file located by r.dbpath location.
// The method guarantee sucessfull data write or panic. Not thhread safe.
func (r *Refs) storeDb() {
    vaultLog.Printf("Store references database to: %s", r.dbpath)

    f, err := ioutil.TempFile(r.opts.TempDir, r.opts.TempPattern)
    if err != nil {
        vaultLog.Fatalf("Can't create temp file! Error: %s", err)
    }

    writer := bufio.NewWriter(f)

    writeProc := func() error {
        defer f.Close()

        encoder := json.NewEncoder(writer)

        err := encoder.Encode(r)
        if err != nil {
            vaultLog.Printf("Refs encoding error: %s", err)
            return err
        }

        err = writer.Flush()
        if err != nil {
            vaultLog.Printf("Writer flush error: %s", err)
            return err
        }

        err = f.Chmod(os.FileMode(r.opts.VaultMode))
        if err != nil {
            vaultLog.Printf("Chmod error: %s", err)
            return err
        }

        return nil
    }

    err = writeProc()
    if err != nil {
        vaultLog.Printf("Remove temp file: %s", f.Name())
        if rmErr := os.Remove(f.Name()); rmErr != nil {
            vaultLog.Printf("Remove temp file error: %s", rmErr)
        }

        vaultLog.Panicf("Write proc error: %s", err)
    }

    err = os.Rename(f.Name(), r.dbpath)
    if err != nil {
        vaultLog.Panicf("Rename error: %s", err)
    }
}


//  Add new reference to object in storage. If the reference already exist
// (total object refs clount, REF_EXIST) will be return. On success will be
// return (total object refs count, nil).
func (r *Refs) Add(object string, id storage_ifaces.StorageId, path storage_ifaces.Path) (int, error) {
    r.Lock()
    defer r.Unlock()

    vaultLog.Printf("vault refs: new object: %s reference for storage: %s path: %s", object, id.Id, path)

    refs, ok := r.values[object]
    if !ok {
        r.values[object] = make(RefsSlice, 0, 100)
    } else {
        // Check id such reference already exist.
        for _, ref := range refs {
            if ref.StorageId == id && ref.Path == path {
                return len(refs), REF_EXIST
            }
        }
    }

    refsCount := len(refs)

    // Add new reference to object references collection.
    refs = append(refs, &Ref{ StorageId: id, Path: path })

    vaultLog.Printf("vault refs add: object: %s refs count: %d", object, refsCount)

    r.storeDb()

    return refsCount, nil
}


// Get references count by object id.
func (r *Refs) RefsCount(object string) int {
    r.Lock()
    defer r.Unlock()

    refs, ok := r.values[object]
    if !ok {
        return 0
    } else {
        return len(refs)
    }
}


//  Remove reference to object in storage.
func (r *Refs) Remove(object string, id storage_ifaces.StorageId, path storage_ifaces.Path) (int, error) {
    r.Lock()
    defer r.Unlock()

    refs, ok := r.values[object]
    if !ok {
        panic(fmt.Errorf("No refs record for object: %s! Called by storage: %s for asset: %s", object, id.Id, path))
    }

    toRemove := -1
    for idx, ref := range refs {
        if ref.StorageId == id && ref.Path == path {
            toRemove = idx
            break
        }
    }

    if toRemove > -1 {
        r.values[object] = append(refs[:toRemove], refs[toRemove + 1:]...)
    }

    refsCount := len(r.values[object])
    if refsCount == 0 {
        delete(r.values, object)
    }

    vaultLog.Printf("vault refs remove: object: %s refs count: %d", object, refsCount)

    r.storeDb()

    return refsCount, nil
}


func (r *Refs) UnmarshalJSON(b []byte) (err error) {
    return json.Unmarshal(b, &r.values)
}


func (r Refs) MarshalJSON() ([]byte, error) {
    b, err := json.Marshal(r.values)
    return b, err
}
