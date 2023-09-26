package content

import (
	b64 "encoding/base64"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type StringData struct {
	contentHash string
	contentData map[string]interface{}
}

func NewStringData(data *unstructured.Unstructured) (*StringData, error) {
	content, err := getContent(data, StringDataContent)
	if err != nil {
		return nil, err
	}
	hash, err := toHash(copyToEncodedValues(content))
	if err != nil {
		return nil, err
	}
	return &StringData{
		contentHash: hash,
		contentData: content,
	}, nil
}

func (s *StringData) Equals(existingHash string) bool {
	return s.contentHash == existingHash
}

func (s *StringData) Hash() string {
	return s.contentHash
}

func copyToEncodedValues(data map[string]interface{}) map[string]interface{} {
	outputs := make(map[string]interface{}, len(data))

	for k, v := range data {
		outputs[k] = b64.StdEncoding.EncodeToString([]byte(v.(string)))
	}
	return outputs
}
