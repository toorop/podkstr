package core

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// GetSHA256File retuen the SHA256 hash of file @ path
func GetSHA256File(path string) (h string, err error) {
	fd, err := os.Open(path)
	if err != nil {
		return
	}
	defer fd.Close()
	hasher := sha256.New()
	if _, err = io.Copy(hasher, fd); err != nil {
		return
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
