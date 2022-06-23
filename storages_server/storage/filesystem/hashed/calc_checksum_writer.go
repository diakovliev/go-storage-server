package hashed_filesystem_storage

import (
    "io"
    "hash"
    "fmt"
    "crypto/sha256"
)

type CalcChecksumsWriter struct {
    io.Writer

    // Selected checksum builder
    checksum    hash.Hash
}


func NewCalcChecksumsWriter(dst io.Writer) *CalcChecksumsWriter {

    checksum    := sha256.New()

    w := &CalcChecksumsWriter{
        Writer      : io.MultiWriter(dst, checksum),
        checksum    : checksum,
    }

    return w
}


func (w *CalcChecksumsWriter) String() string {
    return fmt.Sprintf("%x", w.checksum.Sum(nil))
}
