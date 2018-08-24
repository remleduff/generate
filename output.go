package generate

import (
	"sort"
	"io"
	"fmt"
	"strings"
)

func getOrderedFieldNames(m map[string]Field) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func getOrderedStructNames(m map[string]Struct) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func Output(w io.Writer, structs map[string]Struct, aliases map[string]Field) {
	for _, k := range getOrderedFieldNames(aliases) {
		a := aliases[k]

		fmt.Fprintln(w, "")
		fmt.Fprintf(w, "// %s\n", a.Name)
		fmt.Fprintf(w, "type %s %s\n", a.Name, a.Type)
	}

	for _, k := range getOrderedStructNames(structs) {
		s := structs[k]

		fmt.Fprintln(w, "")
		outputNameAndDescriptionComment(s.Name, s.Description, w)
		fmt.Fprintf(w, "type %s struct {\n", s.Name)

		for _, fieldKey := range getOrderedFieldNames(s.Fields) {
			f := s.Fields[fieldKey]

			// Only apply omitempty if the field is not required.
			omitempty := ",omitempty"
			if f.Required {
				omitempty = ""
			}

			fmt.Fprintf(w, "  %s %s `json:\"%s%s\"`\n", f.Name, f.Type, f.JSONName, omitempty)
		}

		fmt.Fprintln(w, "}")
	}
}

func outputNameAndDescriptionComment(name, description string, w io.Writer) {
	if strings.Index(description, "\n") == -1 {
		fmt.Fprintf(w, "// %s %s\n", name, description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "// %s %s\n", name, strings.Join(dl, "\n// "))
}
