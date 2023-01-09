package replicator

import (
	"fmt"
	"os"
	"testing"

	naisiov1 "nais/replicator/api/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestResources(t *testing.T) {
	values := &TemplateValues{
		Values: map[string]string{
			"foo": "bar",
		},
	}

	b, err := os.ReadFile("testdata/replicationconfig.yaml")
	assert.NoError(t, err)
	var r naisiov1.ReplicationConfig
	err = yaml.Unmarshal(b, &r)
	assert.NoError(t, err)

	resources, err := RenderResources(values, r.Spec.Resources)
	assert.NoError(t, err)
	fmt.Printf("resources: %v\n", resources[0].Object["data"])

}

func TestExtractValues(t *testing.T) {
	ns := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Annotations: map[string]string{
				"foo":             "bar",
				"some.url.io/key": "annotation_value",
			},
			Labels: map[string]string{
				"baz": "qux",
			},
		},
	}

	nsValues := naisiov1.Namespace{
		Annotations: []string{"foo", "some.url.io/key"},
		Labels:      []string{"baz"},
	}

	values := ExtractValues(ns, nsValues)

	fmt.Printf("values: %v\n", values)

	assert.Equal(t, "bar", values["foo"])
	assert.Equal(t, "qux", values["baz"])
	assert.Equal(t, "", values["some.url.io/key"])
	assert.Equal(t, "annotation_value", values["key"])
}
