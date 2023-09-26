package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Spec struct {
	contentHash string
	contentData map[string]interface{}
}

func NewSpec(data *unstructured.Unstructured) (*Spec, error) {
	content, err := getContent(data, SpecContent)
	if err != nil {
		return nil, err
	}
	hash, err := toHash(content)
	if err != nil {
		return nil, err
	}
	return &Spec{
		contentHash: hash,
		contentData: content,
	}, nil
}

func (s *Spec) Equals(existingHash string) bool {
	return s.contentHash == existingHash
}

func (s *Spec) Hash() string {
	return s.contentHash
}
