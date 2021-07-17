package domain

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SloeTree struct {
	rootAlbum SloeNode
}

func NewSloeTree() *SloeTree {
	return &SloeTree{
		rootAlbum: NewSloeNode(),
	}
}

func loadNodeFromPath(filePath string) (*SloeNode, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	node, err := NewSloeNodeFromSource(fp)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (t *SloeTree) LoadFromSource(rootPath string) error {
	iniRegexp := regexp.MustCompile(`(.*)-([A-Z]+)=([0-9A-Fa-f-]{36})\.ini$`)

	foundFiles := map[string][]FileListEntry{}
	foundUuids := map[string]string{}
	orderFiles := []string{}

	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
			fileExt := filepath.Ext(fileInfo.Name())
			if fileExt == ".ini" {
				match := iniRegexp.FindStringSubmatch(fileInfo.Name())
				if len(match) == 0 {
					log.WithFields(log.Fields{"filename": fileInfo.Name()}).Warn("Suspicious INI filename")
				} else {
					nodeType, nodeUuid := match[2], match[3]

					if err != nil {
						return err
					}
					foundFilesForType, ok := foundFiles[nodeType]
					if !ok {
						foundFilesForType = nil
					}
					foundFiles[nodeType] = append(foundFilesForType, *NewFileListEntry(path, filepath.Base(path), nodeUuid))
					foundUuids[nodeUuid] = path
				}
			} else if fileInfo.Name() == "order.txt" {
				orderFiles = append(orderFiles, path)
			}
		}
		return nil
	}
	err := filepath.WalkDir(rootPath, walkDirFunc)
	if err != nil {
		return err
	}

	messages := []string{}
	for k, v := range foundFiles {
		messages = append(messages, fmt.Sprintf("%s:%d", k, len(v)))
	}
	log.Infof("Found objects: %s", strings.Join(messages, ", "))

	albumsBySubpath := map[string]*SloeNode{"": &t.rootAlbum}

	for _, album := range foundFiles["ALBUM"] {
		subTree, err := filepath.Rel(rootPath, filepath.Dir(album.FullPath))
		if err != nil {
			return err
		}
		subTree = strings.Replace(subTree, "\\", "/", -1)
		if subTree == "." {
			subTree = ""
		}
		parentSubtree := filepath.Dir(subTree)
		if parentSubtree == "." {
			parentSubtree = ""
		}
		newAlbum, err := loadNodeFromPath(album.FullPath)
		parentAlbum, ok := albumsBySubpath[parentSubtree]
		if !ok {
			log.WithFields(log.Fields{"subTree": subTree, "uuid": album.FilenameUuid}).Error("Album has no parent")
		} else {
			parentAlbum.AddChildObj(newAlbum)
		}
		albumsBySubpath[subTree] = newAlbum
	}

	return nil
}
