package main

import (
	"encoding/json"
	"strings"

	"github.com/itchyny/gojq"
)

func prettifyJson(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

func processJson(input any, filter string) string {
	query, err := gojq.Parse(filter)
	if err != nil {
		return err.Error()
	}

	iter := query.Run(input)

	var result strings.Builder

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err != nil {
				result.Reset()
				result.WriteString(err.Error())
				break
			}
		}
		b, err := prettifyJson(v)
		if err != nil {
			result.Reset()
			result.WriteString(err.Error())
			break
		}

		result.Write(b)
		result.WriteString("\n")
	}

	return result.String()
}
