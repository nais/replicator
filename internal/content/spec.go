package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Spec struct {
	contentHash string
	contentData map[string]interface{}
	annotations string
	labels      string
}

func NewSpec(data *unstructured.Unstructured) (*Spec, error) {
	content, err := getContent(data, SpecContent)
	contentHash, err := toHash(content)
	annotationsHash, err := toHash(data.GetAnnotations())
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	return &Spec{
		contentHash: contentHash,
		contentData: content,
		annotations: annotationsHash,
		labels:      labelsHash,
	}, nil
}

func (s *Spec) Equals(content ResourceContent) bool {
	return s.labels == content.Labels() &&
		s.annotations == content.Annotations() &&
		s.contentHash == content.Hash()
}

func (s *Spec) Annotations() string {
	return s.annotations
}

func (s *Spec) Hash() string {
	return s.contentHash
}

func (s *Spec) Labels() string {
	return s.labels
}
