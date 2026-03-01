package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

type VaultManager struct {
	root     string
	allFiles []string
}

func (v *VaultManager) Init() {
	v.Index()
}

func (v *VaultManager) Index() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir := wd
	v.root = ""
	v.allFiles = []string{}

	for {
		vault := filepath.Join(dir, ".vault")
		_, err := os.Stat(vault)
		if err == nil {
			v.root = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if v.root != "" {
		filepath.WalkDir(v.root, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				v.allFiles = append(v.allFiles, path)
			}
			return nil
		})
	}
	return nil
}

func (v *VaultManager) GetRoot() string {
	return v.root
}
