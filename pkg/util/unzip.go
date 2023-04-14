package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// based on https://golangcode.com/unzip-files-in-go/
// Code shown Licenced under MIT Licence
func Unzip(src []byte, path string) (filenames []string, err error) {
	br := bytes.NewReader(src)
	ra := io.ReaderAt(br)
	r, err := zip.NewReader(ra, br.Size())
	if err != nil {
		return filenames, err
	}
	rc := zip.ReadCloser{Reader: *r}
	defer rc.Close()

	dest, err := filepath.Abs(path)
	if err != nil {
		return filenames, err
	}

	// refuse to overwrite existing directories
	first := r.File[0].FileInfo()
	if first.IsDir() {
		p := filepath.Join(dest, first.Name())
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			return filenames, fmt.Errorf("path %q already exists", p)
		}
	}

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
