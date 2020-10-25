// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"path/filepath"

	"sigs.k8s.io/kustomize/api/filesys"
)

var _ filesys.FileSystem = relFsOnDisk{}

// relFsOnDisk implements FileSystem using the local filesystem
// with relative paths.
type relFsOnDisk struct {
	fs       filesys.FileSystem
	basePath string
}

// MakeFsOnDisk makes an instance of relFsOnDisk.
func MakeRelFsOnDisk(basePath string) filesys.FileSystem {
	return relFsOnDisk{
		fs:       filesys.MakeFsOnDisk(),
		basePath: basePath,
	}
}

// Create delegates to os.Create.
func (f relFsOnDisk) Create(name string) (filesys.File, error) {
	path := filepath.Join(f.basePath, name)
	return f.fs.Create(path)
}

// Mkdir delegates to os.Mkdir.
func (f relFsOnDisk) Mkdir(name string) error {
	path := filepath.Join(f.basePath, name)
	return f.fs.Mkdir(path)
}

// MkdirAll delegates to os.MkdirAll.
func (f relFsOnDisk) MkdirAll(name string) error {
	path := filepath.Join(f.basePath, name)
	return f.fs.MkdirAll(path)
}

// RemoveAll delegates to os.RemoveAll.
func (f relFsOnDisk) RemoveAll(name string) error {
	path := filepath.Join(f.basePath, name)
	return f.fs.RemoveAll(path)
}

// Open delegates to os.Open.
func (f relFsOnDisk) Open(name string) (filesys.File, error) {
	path := filepath.Join(f.basePath, name)
	return f.fs.Open(path)
}

// CleanedAbs converts the given path into a
// directory and a file name, where the directory
// is represented as a ConfirmedDir and all that implies.
// If the entire path is a directory, the file component
// is an empty string.
func (f relFsOnDisk) CleanedAbs(
	name string) (filesys.ConfirmedDir, string, error) {
	path := filepath.Join(f.basePath, name)
	return f.fs.CleanedAbs(path)
}

// Exists returns true if os.Stat succeeds.
func (f relFsOnDisk) Exists(name string) bool {
	path := filepath.Join(f.basePath, name)
	return f.fs.Exists(path)
}

// Glob returns the list of matching files
func (f relFsOnDisk) Glob(pattern string) ([]string, error) {
	path := filepath.Join(f.basePath, pattern)
	return f.fs.Glob(path)
}

// IsDir delegates to os.Stat and FileInfo.IsDir
func (f relFsOnDisk) IsDir(name string) bool {
	path := filepath.Join(f.basePath, name)
	return f.fs.IsDir(path)
}

// ReadFile delegates to ioutil.ReadFile.
func (f relFsOnDisk) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(f.basePath, name)
	return f.fs.ReadFile(path)
}

// WriteFile delegates to ioutil.WriteFile with read/write permissions.
func (f relFsOnDisk) WriteFile(name string, c []byte) error {
	path := filepath.Join(f.basePath, name)
	return f.fs.WriteFile(path, c)
}

// Walk delegates to filepath.Walk.
func (f relFsOnDisk) Walk(name string, walkFn filepath.WalkFunc) error {
	path := filepath.Join(f.basePath, name)
	return f.fs.Walk(path, walkFn)
}
