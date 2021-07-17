package domain

import (
	"io"

	"gopkg.in/ini.v1"
)

type SloeNode struct {
	Location string
	Name     string
	Type     string
	Uuid     string
	SiteTag  string
	Tags     string
	Title    string
}

func NewSloeNode() SloeNode {
	return SloeNode{}
}

func NewSloeNodeFromSource(reader io.Reader) (*SloeNode, error) {
	node := NewSloeNode()
	return LoadSloeNodeFromSource(&node, reader)
}

func LoadSloeNodeFromSource(node *SloeNode, reader io.Reader) (*SloeNode, error) {
	content, err := ini.Load(reader)
	if err != nil {
		return nil, err
	}
	sectionName := content.SectionStrings()[0]
	section := content.Sections()[0]
	node.Name = section.Key("name").String()
	node.Type = sectionName
	node.Uuid = section.Key("uuid").String()

	return node, nil
}
