package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	file, err := os.ReadFile("testdata/test.yaml")
	assert.NoError(t, err)
	u, err := RenderTemplate(map[string]string{}, string(file), WithOption("missingkey=invalid"))
	assert.NoError(t, err)
	fmt.Printf("map created: %v", u.Object)
	b, err := yaml.Marshal(u.Object)
	assert.NoError(t, err)

	os.WriteFile("testdata/test2.yaml", b, 0644)
}
