package parquet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFromGenericConfig_Success(t *testing.T) {
	c := map[string]any{
		"schema": []map[string]any{
			{
				"name": "uuid",
				"type": "BYTE_ARRAY",
			},
			{
				"name": "name",
				"type": "BYTE_ARRAY",
			},
		},
	}

	_, err := NewFromGenericConfig(c)
	assert.NoError(t, err)
}
