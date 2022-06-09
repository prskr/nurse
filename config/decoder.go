package config

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	_ configDecoder = (*jsonDecoder)(nil)
	_ configDecoder = (*yamlDecoder)(nil)
)

type configDecoder interface {
	DecodeConfig(into *Nurse) error
}

//nolint:ireturn // is required to fulfill a common function signature
func newJSONDecoder(r io.Reader) configDecoder {
	return &jsonDecoder{decoder: json.NewDecoder(r)}
}

type jsonDecoder struct {
	decoder *json.Decoder
}

func (j jsonDecoder) DecodeConfig(into *Nurse) error {
	return j.decoder.Decode(into)
}

//nolint:ireturn // is required to fulfill a common function signature
func newYAMLDecoder(r io.Reader) configDecoder {
	return &yamlDecoder{decoder: yaml.NewDecoder(r)}
}

type yamlDecoder struct {
	decoder *yaml.Decoder
}

func (y yamlDecoder) DecodeConfig(into *Nurse) error {
	return y.decoder.Decode(into)
}
