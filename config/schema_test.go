package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
