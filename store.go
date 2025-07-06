package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "daitya"

func CASPathTranformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) // 20bytes => []bytes => [:]
	hashStr := hex.EncodeToString(hash[:])

	depth := 3
	sliceLen := len(hashStr) / depth
	paths := make([]string, depth)
	for i := 0; i < depth; i++ {
		from, to := i*sliceLen, i*sliceLen+sliceLen
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}

}

type PathKey struct {
	PathName string
	FileName string
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type Store struct {
	StoreOpts StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func (s *Store) Delete(key string) error {
	pathKey := s.StoreOpts.PathTransformFunc(key)

	defer func() {
		fmt.Printf("")
	}()
	return os.RemoveAll(pathKey.FullPath())
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) Has(id string, key string) bool {
	pathKey := s.StoreOpts.PathTransformFunc(key)

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.StoreOpts.Root, id, pathKey.FullPath())

	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.StoreOpts.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) WriteStream(id string, key string, r io.Reader) error {
	pathKey := s.StoreOpts.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.StoreOpts.Root, id, pathKey.FullPath())

	if err := os.MkdirAll(fullPathWithRoot, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk %s ", n, fullPathWithRoot)

	return nil
}
