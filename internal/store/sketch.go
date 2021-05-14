package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go-bootstraper/internal/models"

	"gopkg.in/yaml.v3"
)

type SketchStore interface {
	Load() (*models.Drawing, error)
	Store(drawing *models.Drawing) error
}

var _ SketchStore = &sketchStore{}

type sketchStore struct {
	path string
}

const daFile = ".bt.yml"
const filePerm = os.ModePerm

func NewSketch(path string) *sketchStore {
	if strings.Contains(path, "file://") {
		return &sketchStore{path: strings.Split(path, "file://")[1]}
	}
	return &sketchStore{path: filepath.Join(path, daFile)}
}

func (s *sketchStore) Load() (*models.Drawing, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("cannot load config %w", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot load config %w", err)
	}
	var drawing models.Drawing
	if err := yaml.Unmarshal(b, &drawing); err != nil {
		return nil, fmt.Errorf("cannot load config %w", err)
	}
	return &drawing, nil
}

func (s *sketchStore) Store(drawing *models.Drawing) error {
	os.Create(s.path)
	f, err := os.OpenFile(s.path, os.O_WRONLY, filePerm)
	if err != nil {
		return err
	}
	var conf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&conf)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(&drawing); err != nil {
		return err
	}
	if _, err := f.Write(conf.Bytes()); err != nil {
		return err
	}
	return nil
}
