package go_gen

import (
	"fmt"
	"os"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

const httpStart = `{{ .generationLine }}
package {{ .packageName }}

type {{ .transporterName }} struct{}

func New{{ .transporterName }}() *{{ .transporterName }} {
	return &{{ .transporterName }}{}
}

func (c *{{ .transporterName }}) SetURL() {
	fmt.Println("Hello, neogen!")
}

`

func GenerateHTTPTransporter(name string, genLine string, transporter neoschema.TransporterSchema) (string, error) {
	transporterName := util.ToCamelCase(name+".Connector", true)

	file := fmt.Sprintf(httpStart, genLine, os.Getenv("GOPACKAGE"), transporterName, transporterName, transporterName, transporterName, transporterName)

	// Generate the stuff for all route schemas
	for name, schema := range transporter.Routes {
		generated, err := GenerateRoutes(transporterName, "c.receiver", name, schema)
		if err != nil {
			return file, fmt.Errorf("Couldn't generate route %s: %v", name, err)
		}

		file += generated + "\n\n"
	}

	return file, nil
}
