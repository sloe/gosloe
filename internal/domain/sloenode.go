package domain

import (
	"io"

	"gopkg.in/ini.v1"
)

type SloeNode struct {
	Location     string
	Name         string
	Type         string
	Uuid         string
	SiteTag      string
	Tags         string
	Title        string
	ChildObjects []*SloeNode
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
	for _, section := range content.Sections() {
		if section.Name() != "DEFAULT" {
			node.Name = section.Key("name").String()
			node.Type = section.Name()
			node.Uuid = section.Key("uuid").String()
		}
	}
	return node, nil
}

func (n *SloeNode) AddChildObj(node *SloeNode) {
	n.ChildObjects = append(n.ChildObjects, node)
}
