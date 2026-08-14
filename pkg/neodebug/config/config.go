package config

// This is just here so we don't get into import cycles.

type DebugConfig struct {
	// The name of the transpoter you want to launch (in the schema).
	TransporterName string

	// The URL of the transporter endpoint (HTTP).
	TransporterURL string

	// The command you have to use to get the schema of the server.
	GenerateCommand string
}

var Config DebugConfig
