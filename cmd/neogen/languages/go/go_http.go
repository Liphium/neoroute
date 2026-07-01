package go_gen

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

var httpStart = template.Must(template.New("").Parse(`{{ .generationLine }}
package {{ .packageName }}

import (
	"net/url"

	"github.com/Liphium/neoroute/client"
	"github.com/Liphium/neoroute/client/transporter/http"
)

type {{ .transporterName }} struct{
	*http.HTTPTransporter
	receiver *client.Receiver
}

func New{{ .transporterName }}(config client.Config, method string, u *url.URL) *{{ .transporterName }} {
	r := client.NewReceiver(config)

	return &{{ .transporterName }}{
		HTTPTransporter: http.NewHTTPTransporter(r, method, u),
		receiver: r,
	}
}

`))

func GenerateHTTPTransporter(name string, genLine string, transporter neoschema.TransporterSchema) (string, error) {
	transporterName := util.ToCamelCase(name+".Connector", true)

	var builder strings.Builder
	httpStart.Execute(&builder, map[string]string{
		"generationLine":  genLine,
		"packageName":     os.Getenv("GOPACKAGE"),
		"transporterName": transporterName,
	})
	file := builder.String()

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
