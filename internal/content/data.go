package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Data struct {
	contentHash string
	contentData map[string]interface{}
	annotations string
	labels      string
}

func NewData(data *unstructured.Unstructured) (*Data, error) {
	content, err := getContent(data, DataContent)
	contentHash, err := toHash(content)
	annotationsHash, err := toHash(data.GetAnnotations())
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	return &Data{
		contentHash: contentHash,
		contentData: content,
		annotations: annotationsHash,
		labels:      labelsHash,
	}, nil
}

func (d *Data) Equals(content ResourceContent) bool {
	return d.labels == content.Labels() &&
		d.annotations == content.Annotations() &&
		d.contentHash == content.Hash()
}

func (d *Data) Hash() string {
	return d.contentHash
}

func (d *Data) Annotations() string {
	return d.annotations
}

func (d *Data) Labels() string {
	return d.labels
}
