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
			},
				nil,
				nil,
			),
			rcInput: unstructuredData("some-new-content", map[string]interface{}{
				"key": "value",
			},
				nil,
				nil,
			),
		},
		{
			name: "'stringData' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			},
				nil,
				nil,
			),
		},
		{
			name:           "existingData 'stringData' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key":       base64.StdEncoding.EncodeToString([]byte("my-value")),
				"other-key": base64.StdEncoding.EncodeToString([]byte("my-other-value")),
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "my-value",
			},
				nil,
				nil,
			),
		},
		{
			name:           "rcInput 'stringData' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": base64.StdEncoding.EncodeToString([]byte("my-value")),
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(StringDataContent, map[string]interface{}{
				"key": "otherValue",
			},
				nil,
				nil,
			),
		},
		{
			name: "'data' has not changed, it should return false",
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				nil,
				nil,
			),
		},
		{
			name:           "existingData 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			},
				nil,
				nil),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				nil,
				nil,
			),
		},
		{
			name:           "rcInput 'data' changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(DataContent, map[string]interface{}{
				"other-key": "value",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(DataContent, map[string]interface{}{
				"key": "value",
			},
				nil,
				nil,
			),
		},
		{
			name: "'spec' has not changed, it should return false",
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
		},
		{
			name:           "existingData 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
		},
		{
			name:           "rcInput 'spec' has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				nil,
				nil,
			),
		},
		{
			name:           "rcInput annotations has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				map[string]string{"my-annotation": "my-value"},
				nil,
			),
		},

		{
			name:           "existingData annotations has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				map[string]string{"my-annotation": "my-value"},
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			}, nil, nil,
			),
		},

		{
			name:           "rcInput labels has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				map[string]string{"my-label": "my-value"},
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			},
				nil,
				nil,
			),
		},

		{
			name:           "existingData labels has changed, it should return true",
			expectedChange: true,
			existingData: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "1",
			},
				nil,
				nil,
			),
			rcInput: unstructuredData(SpecContent, map[string]interface{}{
				"replicas": "2",
			}, nil,
				map[string]string{"my-label": "my-value"},
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

func unstructuredData(
	contentKey string,
	contentValues map[string]interface{},
	annotations map[string]string,
	labels map[string]string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":        "test",
				"annotations": annotations,
				"labels":      labels,
			},
			contentKey: contentValues,
		},
	}
}
