package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaAcceptsExamples(t *testing.T) {
	schema := compileConfigSchema(t)

	documents := make(map[string]string)
	for _, name := range []string{"inzibat.json", "quickstart.json"} {
		data, err := os.ReadFile(filepath.Join("..", "examples", name))
		require.NoError(t, err)

		documents[name] = string(data)
	}

	documents["arrayHeaders"] = `{
		"serverPort": 8080,
		"routes": [{
			"method": "GET",
			"path": "/x",
			"fakeResponse": {
				"statusCode": 200,
				"headers": {"Content-Type": ["text/plain"]}
			}
		}]
	}`

	for name, document := range documents {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, validateConfigDocument(t, schema, document))
		})
	}
}

func TestSchemaRejectsInvalidConfigs(t *testing.T) {
	schema := compileConfigSchema(t)

	documents := []struct {
		name     string
		document string
	}{
		{"missingRoutes", `{"serverPort": 8080}`},
		{"unknownKey", `{
			"serverPort": 8080,
			"serverport": 9090,
			"routes": [{"method": "GET", "path": "/x", "fakeResponse": {"statusCode": 200}}]
		}`},
		{"routeWithoutHandler", `{
			"serverPort": 8080,
			"routes": [{"method": "GET", "path": "/x"}]
		}`},
		{"invalidMethod", `{
			"serverPort": 8080,
			"routes": [{"method": "FETCH", "path": "/x", "fakeResponse": {"statusCode": 200}}]
		}`},
	}

	for _, testCase := range documents {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Error(t, validateConfigDocument(t, schema, testCase.document))
		})
	}
}

func compileConfigSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()

	schema, err := jsonschema.NewCompiler().Compile(filepath.Join("..", "examples", "inzibat.schema.json"))
	require.NoError(t, err)

	return schema
}

func validateConfigDocument(t *testing.T, schema *jsonschema.Schema, document string) error {
	t.Helper()

	instance, err := jsonschema.UnmarshalJSON(strings.NewReader(document))
	require.NoError(t, err)

	return schema.Validate(instance)
}

func TestSchemaMatchesModel(t *testing.T) {
	schemaBytes, err := os.ReadFile(filepath.Join("..", "examples", "inzibat.schema.json"))
	require.NoError(t, err)

	var schema map[string]any
	require.NoError(t, json.Unmarshal(schemaBytes, &schema))

	definitions, _ := schema["definitions"].(map[string]any)
	resolved, _ := resolveSchemaRefs(schema, definitions).(map[string]any)
	require.NotNil(t, resolved)

	assert.Equal(t, modelPaths(), schemaPaths(resolved, ""))
}

// resolveSchemaRefs inlines the local definitions so paths can be walked once.
func resolveSchemaRefs(node any, definitions map[string]any) any {
	switch value := node.(type) {
	case map[string]any:
		if reference, ok := value["$ref"].(string); ok {
			name := strings.TrimPrefix(reference, "#/definitions/")
			return resolveSchemaRefs(definitions[name], definitions)
		}

		resolved := make(map[string]any, len(value))
		for key, child := range value {
			resolved[key] = resolveSchemaRefs(child, definitions)
		}

		return resolved
	case []any:
		resolved := make([]any, len(value))
		for index, child := range value {
			resolved[index] = resolveSchemaRefs(child, definitions)
		}

		return resolved
	default:
		return node
	}
}

func schemaPaths(node any, prefix string) map[string]struct{} {
	paths := make(map[string]struct{})

	object, ok := node.(map[string]any)
	if !ok {
		return paths
	}

	if properties, ok := object["properties"].(map[string]any); ok {
		for name, property := range properties {
			path := joinPath(prefix, name)
			paths[path] = struct{}{}

			for child := range schemaPaths(property, path) {
				paths[child] = struct{}{}
			}
		}
	}

	for child := range schemaPaths(object["items"], prefix) {
		paths[child] = struct{}{}
	}

	return paths
}

func modelPaths() map[string]struct{} {
	paths := make(map[string]struct{})

	var walk func(structType reflect.Type, prefix string)
	walk = func(structType reflect.Type, prefix string) {
		for index := range structType.NumField() {
			field := structType.Field(index)

			name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			if name == "" || name == "-" {
				continue
			}

			path := joinPath(prefix, name)
			paths[path] = struct{}{}

			fieldType := field.Type
			if fieldType.Kind() == reflect.Pointer || fieldType.Kind() == reflect.Slice {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct {
				walk(fieldType, path)
			}
		}
	}

	walk(reflect.TypeOf(Cfg{}), "")

	return paths
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}

	return prefix + "." + name
}
