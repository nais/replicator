package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Unknown struct {
	contentHash string
	contentData map[string]interface{}
}

func NewUnknown(data *unstructured.Unstructured) (*Unknown, error) {
	content, err := getContent(data, UnknownContent)
	if err != nil {
		return nil, err
	}
	hash, err := toHash(content)
	if err != nil {
		return nil, err
	}
	return &Unknown{
		contentHash: hash,
		contentData: content,
	}, nil
}

func (u *Unknown) Equals(existingHash string) bool {
	return u.contentHash == existingHash
}

func (u *Unknown) Hash() string {
	return u.contentHash
}
