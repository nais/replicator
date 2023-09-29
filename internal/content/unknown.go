package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Unknown struct {
	contentHash string
	contentData map[string]interface{}
	annotations string
	labels      string
}

func NewUnknown(data *unstructured.Unstructured) (*Unknown, error) {
	content, err := getContent(data, UnknownContent)
	if err != nil {
		return nil, err
	}
	contentHash, err := toHash(content)
	annotationsHash, err := toHash(data.GetAnnotations())
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	return &Unknown{
		contentHash: contentHash,
		contentData: content,
		annotations: annotationsHash,
		labels:      labelsHash,
	}, nil
}

func (u *Unknown) Equals(content ResourceContent) bool {
	return u.labels == content.Labels() &&
		u.annotations == content.Annotations() &&
		u.contentHash == content.Hash()
}

func (u *Unknown) Annotations() string {
	return u.annotations
}

func (u *Unknown) Labels() string {
	return u.labels
}

func (u *Unknown) Hash() string {
	return u.contentHash
}
