package utils_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tags-drive/core/internal/utils"
)

func TestEncrypt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		source interface{}
		result string // must have \n at the end
	}{
		{
			source: struct {
				FieldString string `json:"field_string"`
				FieldInt    int    `json:"field_int"`
			}{
				FieldString: "test",
				FieldInt:    -15,
			},
			result: `{"field_string":"test","field_int":-15}` + "\n",
		},
		{
			source: struct {
				A bool    `json:"field_string"`
				B float64 `json:"field_int"`
			}{
				A: false,
				B: 15.375,
			},
			result: `{"field_string":false,"field_int":15.375}` + "\n",
		},
	}

	for i, tt := range tests {
		buff := bytes.NewBuffer([]byte{})
		err := utils.Encode(buff, tt.source, false, [32]byte{})

		if assert.NoError(err) {
			assert.Equalf(tt.result, buff.String(), "iteration #%d", i+1)
		}
	}
}

func TestDecrypt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		source string
		result interface{}
	}{
		{
			source: `{"field_string":"12345142351623","field_int":11215}`,
			result: map[string]interface{}{
				"field_string": "12345142351623",
				"field_int":    float64(11215),
			},
		},
		{
			source: `{"a":false,"b":-15.5}`,
			result: map[string]interface{}{
				"a": false,
				"b": -15.5,
			},
		},
	}

	for i, tt := range tests {
		buff := bytes.NewBuffer([]byte(tt.source))
		target := make(map[string]interface{})

		err := utils.Decode(buff, &target, false, [32]byte{})
		if assert.NoError(err) {
			assert.Equalf(tt.result, target, "iteration #%d", i+1)
		}
	}
}
