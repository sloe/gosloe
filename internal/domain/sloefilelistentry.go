package domain

type FileListEntry struct {
	FullPath     string
	Name         string
	FilenameUuid string
}

func NewFileListEntry(fullPath string, name string, filenameUuid string) *FileListEntry {
	return &FileListEntry{
		FullPath:     fullPath,
		Name:         name,
		FilenameUuid: filenameUuid,
	}
}
