package cli

import (
	"encoding/json"
	"os"
)

func writeOutput(o interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	// nolint:errcheck,wrapcheck // ignore error, we are writing to stdout
	return encoder.Encode(o)
}
