package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTemplateWithOptions(t *testing.T) {
	file, err := os.ReadFile("testdata/test.yaml")
	assert.NoError(t, err)
	u, err := RenderTemplate(map[string]string{}, string(file), WithOption("missingkey=invalid"))
	assert.NoError(t, err)
	fmt.Printf("map created: %v", u.Object)
}
