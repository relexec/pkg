package basic

import (
	_ "embed"
)

var (
	//go:embed testcreds.txt
	testCredentials string
)

// Config contains configuration information for the HTTP Basic Authorizer.
// This Authorizer should only be used for local development and testing.
type Config struct {
	// LoadFrom indicates where credentials should be loaded from.
	LoadFrom string `json:"load_from"`
}

// Development returns a Config populated with local development configuration.
// Should not be used for anything other than local development and testing.
func Development() Config {
	return Config{
		LoadFrom: "testcreds",
	}
}
