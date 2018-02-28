package jsonschema

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/santhosh-tekuri/jsonschema"
)

// ValidationError is the error type returned by Validate.
type ValidationError = jsonschema.ValidationError

var (
	m       sync.Mutex
	schemas map[string]*jsonschema.Schema
)

// AddSchemaDir receives a path to a directory that contains only valid JSON schema
// definitions. Each definition is compiled and stored in a map for later usage.
func AddSchemaDir(dirname string) error {
	m.Lock()
	defer m.Unlock()

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return err
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("error reading files from %s: %v", dirname, err)
	}

	schemas = map[string]*jsonschema.Schema{}
	for _, file := range files {
		filename := file.Name()

		if schemas[filename] != nil {
			return fmt.Errorf("schema key conflict: JSON schema %s already exists", filename)
		}

		schemas[filename], err = jsonschema.Compile(path.Join(dirname, filename))
		if err != nil {
			return fmt.Errorf("error compiling JSON schema %s: %v", filename, err)
		}
	}

	return nil
}

// Schemas returns a slice containing all the currently available JSON schema definitions.
func Schemas() []string {
	keys := make([]string, 0, len(schemas))

	for k := range schemas {
		keys = append(keys, k)
	}

	return keys
}

// Validate receives a JSON schema name and a reader. It validates that the contents
// of the reader comply with the given schema definition. If the schema name
// does not exists, then an error is returned.
func Validate(schemaName string, r io.Reader) error {
	if schema, ok := schemas[schemaName]; ok {
		return schema.Validate(r)
	}

	return fmt.Errorf("JSON schema %s was not found", schemaName)
}
