// The schema-generate binary reads the JSON schema files passed as arguments
// and outputs the corresponding Go structs.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"github.com/a-h/generate"
	"github.com/a-h/generate/jsonschema"
	generate2 "github.com/a-h/generate/output.json"
)

var (
	o = flag.String("o", "", "The output file for the schema.")
	p = flag.String("p", "main", "The package that the structs are created in.")
	i = flag.String("i", "", "A single file path (used for backwards compatibility).")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "  paths")
		fmt.Fprintln(os.Stderr, "\tThe input JSON Schema files.")
	}

	flag.Parse()

	inputFiles := flag.Args()
	if *i != "" {
		inputFiles = append(inputFiles, *i)
	}
	if len(inputFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No input JSON Schema files.")
		flag.Usage()
		os.Exit(1)
	}

	schemas := make([]*jsonschema.Schema, len(inputFiles))
	for i, file := range inputFiles {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read the input file with error ", err)
			return
		}

		schemas[i], err = jsonschema.Parse(string(b))
		if err != nil {
			if jsonError, ok := err.(*json.SyntaxError); ok {
				line, character, lcErr := lineAndCharacter(b, int(jsonError.Offset))
				fmt.Fprintf(os.Stderr, "Cannot parse JSON schema due to a syntax error at %s line %d, character %d: %v\n", file, line, character, jsonError.Error())
				if lcErr != nil {
					fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
				}
				return
			}
			if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
				line, character, lcErr := lineAndCharacter(b, int(jsonError.Offset))
				fmt.Fprintf(os.Stderr, "The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file %s line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, file, line, character)
				if lcErr != nil {
					fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
				}
				return
			}
			fmt.Fprintf(os.Stderr, "Failed to parse the input JSON schema file %s with error %v\n", file, err)
			return
		}
	}

	g := generate.New(schemas...)

	structs, aliases, err := g.CreateTypes()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failure generating structs: ", err)
	}

	var w io.Writer = os.Stdout

	if *o != "" {
		w, err = os.Create(*o)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening output file: ", err)
			return
		}
	}

	fmt.Fprintln(w, "// Code generated by schema- DO NOT EDIT.")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "package %v\n", *p)

	generate.Output(w, structs, aliases)
}

func lineAndCharacter(bytes []byte, offset int) (line int, character int, err error) {
	lf := byte(0x0A)

	if offset > len(bytes) {
		return 0, 0, fmt.Errorf("couldn't find offset %d in %d bytes", offset, len(bytes))
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range bytes {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			return line, character, nil
		}
	}

	return 0, 0, fmt.Errorf("couldn't find offset %d in %d bytes", offset, len(bytes))
}

