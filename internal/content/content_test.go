package content

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
		error    bool
	}{
		{
			name:    "resource 'unknown' is a unknown resource, it should return fail",
			changed: false,
			error:   true,
			existing: unstructuredData("unknown", map[string]interface{}{
				"key": "value",
			}),
			input: unstructuredData("unknown", map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name:    "resource 'stringData' has not changed, it should return false",
			changed: false,
			existing: unstructuredData("stringData", map[string]interface{}{
				"key": "value",
			}),
			input: unstructuredData("stringData", map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name:    "resource 'data' has not changed, it should return false",
			changed: false,
			existing: unstructuredData("data", map[string]interface{}{
				"key": "value",
			}),
			input: unstructuredData("data", map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name:    "resource 'spec' has not changed, it should return false",
			changed: false,
			existing: unstructuredData("spec", map[string]interface{}{
				"replicas": "1",
			}),
			input: unstructuredData("spec", map[string]interface{}{
				"replicas": "1",
			}),
		},
		{
			name:    "existing resource 'spec' has changed, it should return true",
			changed: true,
			existing: unstructuredData("spec", map[string]interface{}{
				"replicas": "2",
			}),
			input: unstructuredData("spec", map[string]interface{}{
				"replicas": "1",
			}),
		},
		{
			name:    "input resource 'spec' has changed, it should return true",
			changed: true,
			existing: unstructuredData("spec", map[string]interface{}{
				"replicas": "1",
			}),
			input: unstructuredData("spec", map[string]interface{}{
				"replicas": "2",
			}),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			getHash, err := GetHash(tt.input)
			if tt.error {
				assert.Error(t, err)
				return
			}

			changed, err := getHash.ContentHasChanged(tt.existing)
			assert.NoError(t, err)
			assert.Equal(t, tt.changed, changed)
		})
	}
}

func unstructuredData(inputKey string, inputValue map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "test",
			},
			inputKey: inputValue,
		},
	}
}
