package app

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

type ContainerKey struct {
	Files map[string][]byte
}

type Key struct {
	Path      string
	Container io.Reader
}

type KeyReader struct {
}

func CreateKeyReader() KeyReader {
	return KeyReader{}
}

func (obj KeyReader) IsZip(path string) bool {
	return strings.HasSuffix(path, ".zip")
}

func (obj KeyReader) OpenZip(path string, password string) (*Key, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	container := ContainerKey{Files: make(map[string][]byte, 5)}
	validator := NewKeyContainerValidator()
	for _, f := range r.File {
		validator.Add(f.Name)
		localReader, err := f.OpenRaw()
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(localReader)
		if err != nil {
			return nil, err
		}
		container.Files[f.Name] = content
	}

	defer r.Close()

	if err := validator.GetError(); err != nil {
		return nil, err
	}

	reader, err := obj.createArchive(container)
	if err != nil {
		return nil, err
	}

	return &Key{Path: path, Container: reader}, nil
}

func (obj KeyReader) OpenDir(path string) (*Key, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	container := ContainerKey{Files: make(map[string][]byte, 5)}
	if fileInfo.IsDir() {
		dir, err := file.Readdir(-1)
		if err != nil {
			return nil, err
		}

		validator := NewKeyContainerValidator()

		for _, file := range dir {
			validator.Add(file.Name())

			localFile, err := os.Open(path + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			defer localFile.Close()

			content, err := io.ReadAll(localFile)
			if err != nil {
				return nil, err
			}

			container.Files[file.Name()] = content
		}

		if err := validator.GetError(); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("не директория")
	}

	reader, err := obj.createArchive(container)
	if err != nil {
		return nil, err
	}

	return &Key{Path: path, Container: reader}, nil
}

func (obj KeyReader) createArchive(keyContainer ContainerKey) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, content := range keyContainer.Files {
		file, err := os.CreateTemp("", "part-key-*")
		if err != nil {
			return nil, err
		}
		_, err = file.Write(content)
		if err != nil {
			file.Close()
			return nil, err
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}

		header, err := tar.FileInfoHeader(fi, name)
		header.Name = name
		err = tw.WriteHeader(header)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(tw, bytes.NewBuffer(content))
		if err != nil {
			return nil, err
		}
	}

	ok := tw.Close()
	if ok != nil {
		return nil, ok
	}

	return bufio.NewReader(&buf), nil
}
