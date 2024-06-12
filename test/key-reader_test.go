package test

import (
	"demo-ui/internal/key_reader"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestValidDirPath(t *testing.T) {
	reader := key_reader.CreateKeyReader()
	pwd, _ := os.Getwd()
	path := pwd + "/testdata/container.024"
	assert.Equal(t, reader.IsZip(path), false)
	key, err := reader.OpenDir(path)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, key.Path, path)
}

func TestCreateReader(t *testing.T) {
	archive, err := os.Create("archive.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer archive.Close()

}

func TestInValidDirPath(t *testing.T) {
	reader := key_reader.CreateKeyReader()
	pwd, _ := os.Getwd()
	path := pwd + "/testdata/invalid.container.024"
	assert.Equal(t, reader.IsZip(path), false)
	_, err := reader.OpenDir(path)
	if err == nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, err.Error(), "header.key, masks.key, masks2.key, name.key, primary.key, primary2.key файлы не найдены")
}

func TestInValidZip(t *testing.T) {
	reader := key_reader.CreateKeyReader()
	pwd, _ := os.Getwd()
	path := pwd + "/testdata/invalid.container.024.zip"
	assert.Equal(t, reader.IsZip(path), true)
	_, err := reader.OpenZip(path, "")
	if err == nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, err.Error(), "header.key, masks.key, masks2.key, name.key, primary.key, primary2.key файлы не найдены")
}

func TestValidZip(t *testing.T) {
	reader := key_reader.CreateKeyReader()
	pwd, _ := os.Getwd()
	path := pwd + "/testdata/container.024.zip"
	assert.Equal(t, reader.IsZip(path), true)
	key, err := reader.OpenZip(path, "")
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, key.Path, path)
}
