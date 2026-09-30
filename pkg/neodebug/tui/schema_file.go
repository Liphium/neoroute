package tui

import (
	"fmt"
	"os"
)

var _ SchemaNode = &FileNode{}

// FileNode edits a byte array as a file path, without path completion.
type FileNode struct {
	ValueNode[[]byte]
}

func newFileNode() *FileNode {
	return &FileNode{
		prefix: "file(\"",
		suffix: "\")",
		value:  []byte{},
		convert: func(path string) ([]byte, error) {
			if path == "" {
				return []byte{}, nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("Could not read file: %w", err)
			}
			return data, nil
		}}
}

func (f *FileNode) Init() {
	f.ValueNode.Init()
	// Display the path, not the byte array's formatted value.
	f.input.SetValue("")
}
