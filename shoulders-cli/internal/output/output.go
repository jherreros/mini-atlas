package output

import (
	"encoding/json"
	"errors"

	"github.com/pterm/pterm"
	"sigs.k8s.io/yaml"
)

type Format string

const (
	Table Format = "table"
	JSON  Format = "json"
	YAML  Format = "yaml"
)

func ParseFormat(value string) (Format, error) {
	switch Format(value) {
	case Table, JSON, YAML:
		return Format(value), nil
	default:
		return "", errors.New("invalid output format, must be table|json|yaml")
	}
}

func Render(data any, format Format) ([]byte, error) {
	switch format {
	case JSON:
		return json.MarshalIndent(data, "", "  ")
	case YAML:
		return yaml.Marshal(data)
	default:
		return nil, nil
	}
}

func PrintTable(headers []string, rows [][]string) error {
	data := pterm.TableData{headers}
	data = append(data, rows...)
	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
