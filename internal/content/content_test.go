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
	}{
		{
			name:          "'unknown' should return fail",
			expectedError: true,
			existingData: unstructuredData(UnknownContent, map[string]interface{}{
				"key": "value",
			}),
			rcInput: unstructuredData("some-new-content", map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name: "'stringData' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			}),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			}),
		},

		{
			name:           "existingData 'stringData' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key":       base64.StdEncoding.EncodeToString([]byte("my-value")),
				"other-key": base64.StdEncoding.EncodeToString([]byte("my-other-value")),
			}),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			}),
		},
		{
			name:           "rcInput 'stringData' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			}),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "otherValue",
			}),
		},
		{
			name: "'data' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			}),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name:           "existingData 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			}),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name:           "rcInput 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			}),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			}),
		},
		{
			name: "'spec' has not changed, it should return false",
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			}),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			}),
		},
		{
			name:           "existingData 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			}),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			}),
		},
		{
			name:           "rcInput 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			}),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			}),
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
