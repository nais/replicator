package content

import (
	b64 "encoding/base64"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type StringData struct {
	contentHash string
	contentData map[string]interface{}
	annotations string
	labels      string
}

func NewStringData(data *unstructured.Unstructured) (*StringData, error) {
	content, err := getContent(data, StringDataContent)
	contentHash, err := toHash(copyToEncodedValues(content))
	annotationsHash, err := toHash(data.GetAnnotations())
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	return &StringData{
		contentHash: contentHash,
		contentData: content,
		annotations: annotationsHash,
		labels:      labelsHash,
	}, nil
}

func (s *StringData) Equals(content ResourceContent) bool {
	return s.labels == content.Labels() &&
		s.annotations == content.Annotations() &&
		s.contentHash == content.Hash()
}

func (s *StringData) Hash() string {
	return s.contentHash
}

func (s *StringData) Annotations() string {
	return s.annotations
}

func (s *StringData) Labels() string {
	return s.labels
}

func copyToEncodedValues(data map[string]interface{}) map[string]interface{} {
	outputs := make(map[string]interface{}, len(data))

	for k, v := range data {
		outputs[k] = b64.StdEncoding.EncodeToString([]byte(v.(string)))
	}
	return outputs
}
