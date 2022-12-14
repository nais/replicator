package parser

import (
	"fmt"
	"os"
	"testing"

	naisiov1 "nais/replicator/api/v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestTodo(t *testing.T) {
	values := &TemplateValues{
		Values: map[string]string{
			"foo": "bar",
		},
	}

	b, err := os.ReadFile("testdata/replicatorconfig.yaml")
	assert.NoError(t, err)
	var r naisiov1.ReplicatorConfiguration
	err = yaml.Unmarshal(b, &r)
	assert.NoError(t, err)

	resources, err := Resources(values, r.Spec.Resources)
	assert.NoError(t, err)
	fmt.Printf("resources: %v\n", resources[0].Object["data"])

}
