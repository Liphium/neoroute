package languages

import (
	"github.com/Liphium/neoroute/v2/cmd/neogen/engine"
	"github.com/Liphium/neoroute/v2/cmd/neogen/internal/testruntime"
	"github.com/Liphium/neoroute/v2/cmd/neogen/internal/testschema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypeScriptRuntime(t *testing.T) {
	schema, err := testschema.Schema()
	require.NoError(t, err)
	files, err := engine.NewGenerationEngine(NewTSConfig()).Generate(schema)
	require.NoError(t, err)
	testruntime.Run(t, "../tests/typescript", files)
}
