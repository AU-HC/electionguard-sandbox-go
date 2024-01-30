package generation

import (
	"electionguard-sandbox-go/models"
	"encoding/json"
	"io"
	"os"
)

func LoadManifest(path string) models.Manifest {
	manifest := models.Manifest{}

	// Open json file and print error if any
	file, fileErr := os.Open(path)
	if fileErr != nil {
		panic(fileErr)
	}

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	if byteErr != nil {
		panic(fileErr)
	}

	// Unmarshal the bytearray into empty instance of variable of type E
	jsonErr := json.Unmarshal(jsonByte, &manifest)
	if jsonErr != nil {
		panic(fileErr)
	}

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(file)

	return manifest
}
