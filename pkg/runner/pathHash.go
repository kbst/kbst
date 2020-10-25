package runner

import (
	"crypto/sha512"
	"encoding/hex"
	"path/filepath"
)

// Hash of absolute path used to ensure unique
// docker image tags and volumes per repository
func pathHash(path string) (hash string, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return hash, err
	}

	h := sha512.New()
	h.Write([]byte(absPath))
	hash = hex.EncodeToString(h.Sum(nil))[0:7]

	return hash, nil
}
