package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterByType_SingleEntity(t *testing.T) {
	data := []byte(`[{"_type":"amadla.org/entity/package@v1.0.0","_body":{"name":"nginx"}},{"_type":"amadla.org/entity/application@v1.0.0","_body":{"name":"myapp"}}]`)

	result, err := filterByType(data, "package")
	assert.NoError(t, err)
	assert.Contains(t, string(result), "nginx")
	assert.NotContains(t, string(result), "myapp")
}

func TestFilterByType_NoMatch(t *testing.T) {
	data := []byte(`[{"_type":"amadla.org/entity/package@v1.0.0","_body":{"name":"nginx"}}]`)

	result, err := filterByType(data, "network")
	assert.NoError(t, err)
	assert.Equal(t, "null", string(result))
}

func TestFilterByType_SingleObject(t *testing.T) {
	data := []byte(`{"_type":"amadla.org/entity/package@v1.0.0","_body":{"name":"nginx"}}`)

	result, err := filterByType(data, "package")
	assert.NoError(t, err)
	assert.Contains(t, string(result), "nginx")
}

func TestFilterByType_NotJSON(t *testing.T) {
	data := []byte(`not json at all`)

	result, err := filterByType(data, "package")
	assert.NoError(t, err)
	assert.Equal(t, data, result) // pass through
}
