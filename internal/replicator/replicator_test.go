package replicator

import (
	"fmt"
	"os"
	"testing"

	naisiov1 "nais/replicator/api/v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestResources(t *testing.T) {
	values := &TemplateValues{
		Values: map[string]string{
			"foo": "bar",
		},
	}

	b, err := os.ReadFile("testdata/replicatorconfig.yaml")
	assert.NoError(t, err)
	var r naisiov1.ReplicationConfig
	err = yaml.Unmarshal(b, &r)
	assert.NoError(t, err)

	resources, err := RenderResources(values, r.Spec.Resources)
	assert.NoError(t, err)
	fmt.Printf("resources: %v\n", resources[0].Object["data"])

}
