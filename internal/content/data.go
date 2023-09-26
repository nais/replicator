package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Data struct {
	contentHash string
	contentData map[string]interface{}
}

func NewData(data *unstructured.Unstructured) (*Data, error) {
	content, err := getContent(data, DataContent)
	if err != nil {
		return nil, err
	}
	hash, err := toHash(content)
	if err != nil {
		return nil, err
	}
	return &Data{
		contentHash: hash,
		contentData: content,
	}, nil
}

func (d *Data) Equals(content ResourceContent) bool {
	return d.contentHash == content.Hash()
}

func (d *Data) Hash() string {
	return d.contentHash
}
