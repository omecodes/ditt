package ditt

import (
	"bytes"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Files is a convenience for UserData file persistence
type Files interface {
	Save(userId string, data string) error
	Delete(userId string) error
	Get(userId string) (string, error)
}

type memoryFiles struct {
	fs afero.Fs
}

func (m *memoryFiles) Save(userId string, data string) error {
	return afero.WriteFile(m.fs, userId, []byte(data), os.ModePerm)
}

func (m *memoryFiles) Delete(id string) error {
	return m.fs.Remove(id)
}

func (m *memoryFiles) Get(userId string) (string, error) {
	file, err := m.fs.Open(userId)
	if err != nil {
		if os.IsNotExist(err) {
			return "", NotFound
		}
		return "", Internal
	}

	defer func() {
		_ = file.Close()
	}()

	data, err := ioutil.ReadAll(file)
	return string(data), err
}

func NewMemoryFiles() Files {
	return &memoryFiles{
		fs: afero.NewMemMapFs(),
	}
}

type dirFiles struct {
	rootDir string
}

func (d *dirFiles) Save(userId string, data string) error {
	filename := filepath.Join(d.rootDir, userId)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, bytes.NewBufferString(data))
	return err
}

func (d *dirFiles) Delete(userId string) error {
	filename := filepath.Join(d.rootDir, userId)
	return os.Remove(filename)
}

func (d *dirFiles) Get(userId string) (string, error) {
	filename := filepath.Join(d.rootDir, userId)
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return "", NotFound
		}
		return "", Internal
	}

	defer func() {
		_ = file.Close()
	}()

	stats, _ := file.Stat()
	if stats.IsDir() {
		return "", nil
	}

	data, err := ioutil.ReadAll(file)
	return string(data), err
}

func NewDirFiles(rootDir string) Files {
	return &dirFiles{rootDir: rootDir}
}
