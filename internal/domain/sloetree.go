package domain

import (
	"io/fs"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type SloeTree struct {
	// root SloeNode
}

func NewSloeTree() *SloeTree {
	return &SloeTree{}
}

func (t *SloeTree) LoadFromSource(rootPath string) error {
	iniRegexp := regexp.MustCompile(`(.*)-([A-Z]+)=([0-9A-Fa-f-]{36})\.ini$`)
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
			fileExt := filepath.Ext(fileInfo.Name())
			if fileExt == ".ini" {
				match := iniRegexp.FindStringSubmatch(fileInfo.Name())
				if err != nil {
					return err
				}
				if len(match) == 0 {
					log.WithFields(log.Fields{"filename": fileInfo.Name()}).Warn("Suspicious INI filename")
				} else {
					log.Infof("Matched ini name %+v", match)
				}
			}
		}
		return nil
	}
	err := filepath.WalkDir(rootPath, walkDirFunc)
	if err != nil {
		return err
	}
	return nil
}
