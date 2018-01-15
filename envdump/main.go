// json-env dumps azure settings to stdout
// as environment variable bindings
package main

import (
	_ "azure-request-limitometer/common"
	"os"
)

func main() {
	// `common` sets the various environment variables
	for _, v := range os.Environ() {
		if v[:len("AZURE_")] == "AZURE_" {
			fmt.Println(v)
		}
	}
}
