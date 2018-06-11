package main

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pkg/xattr"
)

const (
	tempPrefix = "tmp-"

	xattrMD5 = "user.md5"
)

type TempFile struct {
	path string
	h    hash.Hash
	f    *os.File
}

func NewTempFile(path string) (t *TempFile, err error) {
	base := tempPrefix + filepath.Base(path)
	f, err := ioutil.TempFile(filepath.Dir(path), base)
	if err != nil {
		err = errors.Wrap(err, "create temp file")
		return
	}

	t = &TempFile{
		path: path,
		h:    md5.New(),
		f:    f,
	}
	return
}

func (t *TempFile) Write(b []byte) (n int, err error) {
	n, err = t.f.Write(b)
	t.h.Write(b[:n])
	return
}

func (t *TempFile) Close() (path string, err error) {
	checksum := t.h.Sum(nil)
	md5 := []byte(hex.EncodeToString(checksum))
	if err = xattr.Set(t.f.Name(), xattrMD5, md5); err != nil {
		err = errors.Wrapf(err, "set xattr for %s", xattrMD5)
		return
	}

	if err = t.f.Close(); err != nil {
		err = errors.Wrap(err, "close temp file")
		return
	}

	path = t.path
	if err = os.Rename(t.f.Name(), path); err != nil {
		err = errors.Wrap(err, "rename temp file")
		return
	}

	return
}
