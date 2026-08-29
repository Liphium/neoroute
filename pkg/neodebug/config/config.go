package config

// This is just here so we don't get into import cycles.

type DebugConfig struct {
	// The name of the transpoter you want to launch (in the schema).
	TransporterName string

	// The URL of the transporter endpoint (HTTP).
	TransporterURL string

	// Which HTTP method to use to connect to the transporter (e.g. POST, GET). Some transporters ignore this, but the HTTP transporter for example needs it.
	//
	// Default for HTTP: POST
	TransporterMethod string

	// The command you have to use to get the schema of the server.
	GenerateCommand string
}

var Config DebugConfig
