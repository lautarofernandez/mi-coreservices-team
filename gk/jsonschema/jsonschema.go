package jsonschema

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

// ValidationError is the error type returned by Validate.
type ValidationError struct {
	Schema string
	Errors []gojsonschema.ResultError
}

// ErrorsDescription returns a map that contains JSON attribute paths as keys,
// and an error description detailing what the error cause was.
func (v ValidationError) ErrorsDescription() map[string]string {
	errors := map[string]string{}

	for _, err := range v.Errors {
		field := err.Field()
		if err.Type() == "required" {
			field = fmt.Sprintf("%s.%s", err.Context().String(), err.Field())
			field = strings.Replace(field, "(root).", "", 1)
		}

		errors[field] = err.Description()
	}

	return errors
}

func (v ValidationError) Error() string {
	buf := bytes.NewBuffer(nil)

	for field, desc := range v.ErrorsDescription() {
		fmt.Fprintf(buf, "%s: %s\n", field, desc)
	}

	return buf.String()
}

var (
	m       sync.Mutex
	schemas map[string]*gojsonschema.Schema
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

	schemas = map[string]*gojsonschema.Schema{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		if schemas[filename] != nil {
			return fmt.Errorf("schema key conflict: JSON schema %s already exists", filename)
		}

		bytes, err := ioutil.ReadFile(path.Join(dirname, filename))
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", file.Name(), err)
		}

		schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(bytes))
		if err != nil {
			return fmt.Errorf("error compiling JSON schema %s: %v", filename, err)
		}

		schemas[filename] = schema
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
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	loader := gojsonschema.NewBytesLoader(bytes)

	schema, exists := schemas[schemaName]
	if !exists {
		return fmt.Errorf("JSON schema %s was not found", schemaName)
	}

	res, err := schema.Validate(loader)
	if err != nil {
		return fmt.Errorf("error validating JSON through the schema: %v", err)
	}

	if res.Valid() {
		return nil
	}

	return &ValidationError{
		Schema: schemaName,
		Errors: res.Errors(),
	}
}
