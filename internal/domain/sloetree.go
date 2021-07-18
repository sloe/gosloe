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
	Albums    map[string]*SloeNode
	Items     map[string]*SloeNode
	Nodes     map[string]*SloeNode
}

func NewSloeTree() *SloeTree {
	return &SloeTree{
		rootAlbum: NewSloeNode(),
		Albums:    map[string]*SloeNode{},
		Items:     map[string]*SloeNode{},
		Nodes:     map[string]*SloeNode{},
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
		parentSubtree := strings.Replace(filepath.Dir(subTree), "\\", "/", -1)
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
		t.Albums[newAlbum.Uuid] = newAlbum
		albumsBySubpath[subTree] = newAlbum
	}

	parentAlbumFromPath := func(objPath string) (*SloeNode, error) {
		subTree, err := filepath.Rel(rootPath, filepath.Dir(objPath))
		if err != nil {
			return nil, err
		}
		subTree = strings.Replace(subTree, "\\", "/", -1)
		if subTree == "." {
			subTree = ""
		}
		parentAlbum, ok := albumsBySubpath[subTree]
		if !ok {
			return nil, fmt.Errorf("Object has no parent album")
		}
		return parentAlbum, nil
	}

	for _, itemDef := range foundFiles["ITEM"] {
		destAlbum, err := parentAlbumFromPath(itemDef.FullPath)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"itemDef": itemDef}).Error("Object has no parent album - not loading")
		} else {
			newItem, err := loadNodeFromPath(itemDef.FullPath)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{"itemDef": itemDef}).Error("Failed to load item")
			} else {
				destAlbum.AddChildObj(newItem)
				t.Items[newItem.Uuid] = newItem
			}
		}
	}

	for objType, objList := range foundFiles {
		switch objType {
		case "ALBUM":
		case "ITEM":
		default:
			for _, objDef := range objList {
				destAlbum, err := parentAlbumFromPath(objDef.FullPath)
				if err != nil {
					log.WithError(err).WithFields(log.Fields{"objDef": objDef}).Error("Object has no parent album - not loading")
				} else {
					newItem, err := loadNodeFromPath(objDef.FullPath)
					if err != nil {
						log.WithError(err).WithFields(log.Fields{"objDef": objDef}).Error("Failed to load item")
					} else {
						destAlbum.AddChildObj(newItem)
						t.Nodes[newItem.Uuid] = newItem
					}
				}
			}
		}
	}

	return nil
}
