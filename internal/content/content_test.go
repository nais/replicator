package content

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func TestContentHasChanged(t *testing.T) {
	for _, tt := range []struct {
		name           string
		existingData   *unstructured.Unstructured
		rcInput        *unstructured.Unstructured
		expectedChange bool
		expectedError  bool
		labels         bool
		annotations    bool
	}{
		{
			name: "'stringData' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			},
				false,
				false,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			},
				false,
				false,
			),
		},
		{
			name:           "existingData 'stringData' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key":       base64.StdEncoding.EncodeToString([]byte("my-value")),
				"other-key": base64.StdEncoding.EncodeToString([]byte("my-other-value")),
			},
				false,
				false,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			},
				false,
				false,
			),
		},
		{
			name:           "rcInput 'stringData' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			},
				false,
				false,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "otherValue",
			},
				false,
				false,
			),
		},
		{
			name: "'data' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				false,
				false,
			),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				false,
				false,
			),
		},
		{
			name:           "existingData 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			},
				false,
				false,
			),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				false,
				false,
			),
		},
		{
			name:           "rcInput 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			},
				false,
				false,
			),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				false,
				false,
			),
		},
		{
			name: "'spec' has not changed, it should return false",
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				false,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				false,
				false,
			),
		},
		{
			name:           "existingData 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				false,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				false,
				false,
			),
		},
		{
			name:           "rcInput 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				false,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				false,
				false,
			),
		},
		{
			name:           "rcInput annotations has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{},
				true,
				false,
			),
		},
		{
			name:           "existingData annotations has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{},
				true,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				false,
			),
		},

		{
			name:           "rcInput labels has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				false,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				true,
			),
		},

		{
			name:           "existingData labels has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				true,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{},
				false,
				false,
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rcContent, err := Get(tt.rcInput)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			existingContent, err := Get(tt.existingData)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.expectedChange, !rcContent.Equals(existingContent))
		})
	}
}

func unstructuredData(contentKey string, contentValues map[string]interface{}, annotations bool, labels bool) *unstructured.Unstructured {
	if annotations {
		return &unstructured.Unstructured{
			Object: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        "test",
					"annotations": map[string]interface{}{"my-annotation": "my-value"},
				},
				contentKey: contentValues,
			},
		}
	}
	if labels {
		return &unstructured.Unstructured{
			Object: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":   "test",
					"labels": map[string]interface{}{"my-label": "my-value"},
				},
				contentKey: contentValues,
			},
		}
	}
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "test",
			},
			contentKey: contentValues,
		},
	}
}
