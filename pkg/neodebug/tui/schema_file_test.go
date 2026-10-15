package tui

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Liphium/neoroute/neoschema"
	"github.com/stretchr/testify/assert"
)

func TestByteArrayFileInput(t *testing.T) {
	node := createNode(&neoschema.ArrayType{
		Element: &neoschema.BasicType{ActualType: neoschema.TypeByte},
	}, nil)
	file, ok := node.(*FileNode)
	if !assert.True(t, ok, "expected FileNode, got %T", node) {
		return
	}
	file.Init()
	assert.Empty(t, file.input.Value())
	assert.NoError(t, file.input.Err)

	path := filepath.Join(t.TempDir(), "binary file")
	data := []byte{0, 1, 128, 255}
	assert.NoError(t, os.WriteFile(path, data, 0600))
	file.input.SetValue(path)
	assert.NoError(t, file.input.Err)
	assert.True(t, bytes.Equal(data, file.Request().([]byte)))

	file.input.SetValue(path + ".missing")
	assert.Error(t, file.input.Err)
	assert.Equal(t, 2, file.Height())
	file.input.SetValue(filepath.Dir(path))
	assert.Error(t, file.input.Err)
	file.input.SetValue("")
	assert.NoError(t, file.input.Err)
	assert.Empty(t, file.Request().([]byte))
}

func TestNonByteArrayUsesSliceNode(t *testing.T) {
	node := createNode(&neoschema.ArrayType{
		Element: &neoschema.BasicType{ActualType: neoschema.TypeString},
	}, nil)
	assert.IsType(t, &SliceNode{}, node)
}
