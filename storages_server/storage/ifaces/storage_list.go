package storage_ifaces

import (
    "encoding/json"
    "fmt"
)

type element struct {
    Path   string               `json:"path"`
    Props  map[string]string    `json:"properties"`
}

func (s *Storage) List() ([]byte, error) {
    elements := make([]element, 0, 100)

    s.Range(func(path Path, opts StorageAssetOpts) bool {
        props := make(map[string]string)
        props["mode"] = fmt.Sprintf("%d", opts.Mode)
        elements = append(elements, element{Path: path, Props: props})
        return true
    })

    resp, err := json.Marshal(elements)
    return resp, err
}
