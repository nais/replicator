package resources

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func TestHasChanged(t *testing.T) {
	for _, tt := range []struct {
		name     string
		existing *unstructured.Unstructured
		input    *unstructured.Unstructured
		changed  bool
	}{
		{
			name:     "resource spec has not changed, it should return false",
			changed:  false,
			existing: unstructuredData(1, nil),
			input:    unstructuredData(1, nil),
		},
		{
			name:     "resource spec has changed, it should return true",
			changed:  true,
			existing: unstructuredData(2, nil),
			input:    unstructuredData(1, nil),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			changed, err := HasChanged(tt.existing, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.changed, changed)
		})
	}
}

func unstructuredData(replicas int, inputLabels map[string]interface{}) *unstructured.Unstructured {
	var labels map[string]interface{}
	if inputLabels != nil {
		labels = inputLabels
	} else {
		labels = map[string]interface{}{
			"app": "test",
		}
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":   "test",
				"labels": labels,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
			},
		},
	}
}
