package template

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TemplateValues struct {
	Values map[string]string
}

func TestRenderTemplateWithOptions(t *testing.T) {
	file, err := os.ReadFile("testdata/test.yaml")
	assert.NoError(t, err)
	u, err := RenderTemplate(TemplateValues{Values: map[string]string{"tommy.johnny": "foo"}}, string(file), WithOption("missingkey=invalid"))
	assert.NoError(t, err)
	fmt.Printf("map created: %v", u.Object)
}
