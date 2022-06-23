package buffers

import (
    "os"
    "io"
    "fmt"
    "path/filepath"

    "../ifaces"
    "../filesystem"
)


type BuffersManager struct {
    Opts storage_ifaces.BuffersManagerOpts
}


func NewBuffersManager(opts storage_ifaces.BuffersManagerOpts) *BuffersManager {

    buffersLog.Printf("Create buffers manager. Opts: %s", opts.String())

    if err := filesystem_utils.EnsureDir(opts.StorageRoot, os.FileMode(opts.StorageRootMode)); err != nil {
        buffersLog.Panicf("Can't initialize buffers storage root directory. Err: %s", err)
    }

    return &BuffersManager{Opts: opts}
}


func (bm *BuffersManager) Abspath(bid string) string {
    return filepath.Join(bm.Opts.StorageRoot, bid)
}


func (bm *BuffersManager) EnsureBuffer(bid string) error {

    if _, err := os.Stat(bm.Abspath(bid)); os.IsNotExist(err) {
        return err
    }

    return nil
}


func (bm *BuffersManager) Create() (string, error) {

    bufferId := storage_ifaces.MakeNewStorageId()
    bid      := bufferId.String()

    if err := bm.EnsureBuffer(bid); err == nil {
        buffersLog.Panicf("Buffer { id: %s } storage file already exists!", bid)
    }

    if _, err := os.OpenFile(bm.Abspath(bid), os.O_CREATE, os.FileMode(bm.Opts.FilesMode)); err != nil {
        rerr := fmt.Errorf("Unable to create storage file for new buffer {id: %s}. Err: %w", bid, err)
        buffersLog.Print(rerr.Error())
        return "", rerr
    }

    buffersLog.Printf("Create new buffer {id: %s}", bid)

    return bid, nil
}


func (bm *BuffersManager) Discard(bid string) error {

    if err := bm.EnsureBuffer(bid); err != nil {
        rerr := fmt.Errorf("Unable to remove storage file for buffer {id: %s}. Err: %w", bid, err)
        buffersLog.Print(rerr.Error())
        return rerr
    }

    buffersLog.Printf("Discard buffer {id: %s}", bid)

    return os.Remove(bm.Abspath(bid))
}


func (bm *BuffersManager) Append(bid string, source io.Reader) (int, error) {

    if err := bm.EnsureBuffer(bid); err != nil {
        rerr := fmt.Errorf("Unable to append data to storage file for buffer {id: %s}. Err: %w", bid, err)
        buffersLog.Print(rerr.Error())
        return 0, rerr
    }

    f, err := os.OpenFile(bm.Abspath(bid), os.O_WRONLY|os.O_APPEND, os.FileMode(bm.Opts.FilesMode))
    if err != nil {
        rerr := fmt.Errorf("Unable to open storage file for buffer {id: %s}. Err: %w", bid, err)
        buffersLog.Print(rerr.Error())
        return 0, rerr
    }
    defer f.Close()

    copied, err := io.Copy(f, source)

    buffersLog.Printf("Appended %d bytes of data to buffer {id: %s}", copied, bid)

    return int(copied), err
}
