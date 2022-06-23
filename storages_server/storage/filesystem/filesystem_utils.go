package filesystem_utils

import (
    "os"
    "fmt"
)


func EnsureDir(path string, mode os.FileMode) error {

    if _, ok := os.Stat(path); !os.IsNotExist(ok) {
        fsUtilsLog.Printf("The directory already exist: %s", path)
        return nil
    }

    err := os.MkdirAll(path, mode)
    if err != nil {
        return fmt.Errorf("Directory creation error: %s. Path: %s", err, path)
    }

    return nil
}
