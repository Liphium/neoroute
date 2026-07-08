package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Liphium/magic/v3"
	"github.com/Liphium/magic/v3/mconfig"
)

var TestDirectory = "test"

func TestMain(m *testing.M) {
	magic.PrepareTesting(m, magic.Config{
		AppName: "tests-neogen",
		PlanDeployment: func(ctx *mconfig.Context) {
			TestDirectory = filepath.Join(ctx.ProjectDirectory(), "neogen-env")
			if err := os.RemoveAll(TestDirectory); err != nil {
				panic(err)
			}
			if err := os.Mkdir(TestDirectory, os.ModePerm); err != nil {
				panic(err)
			}
		},
		StartFunction: magic.AppStarted,
	})
}

func TestGenerate(t *testing.T) {
}
