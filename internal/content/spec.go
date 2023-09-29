package content

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Spec struct {
	contentHash string
	annotations string
	labels      string
}

func NewSpec(data *unstructured.Unstructured) (*Spec, error) {
	content, err := getContent(data, SpecContent)
	if err != nil {
		return nil, err
	}
	contentHash, err := toHash(content)
	if err != nil {
		return nil, err
	}
	annotationsHash, err := toHash(data.GetAnnotations())
	if err != nil {
		return nil, err
	}
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &Spec{
		contentHash: contentHash,
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
