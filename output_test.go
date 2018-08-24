package generate

import (
	"testing"
	"github.com/a-h/generate/jsonschema"
	"strings"
)

func TestGeneration(t *testing.T) {
	s := `{
	    "$schema": "http://json-schema.org/draft-04/schema#",
	    "name": "Example",
	    "type": "object",
	    "properties": {
	        "name": {
	            "type": ["object", "array", "integer"],
	            "description": "name"
			},
			"time": {
				"type": "string",
				"format": "date"
			}
	    }
	}`

	schema, err := jsonschema.Parse(s)

	if err != nil {
		t.Fatal(err)
	}

	generator := New(schema)

	structs, aliases, err := generator.CreateTypes()

	if err != nil {
		t.Fatal(err)
	}

	builder := strings.Builder{}

	Output(&builder, structs, aliases)

	t.Log(builder.String())
}



