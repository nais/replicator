package content

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Generic struct {
	annotations string
	labels      string
}

func NewGeneric(data *unstructured.Unstructured) (*Generic, error) {
	annotationsHash, err := toHash(data.GetAnnotations())
	if err != nil {
		return nil, err
	}
	labelsHash, err := toHash(data.GetLabels())
	if err != nil {
		return nil, err
	}
	return &Generic{
		annotations: annotationsHash,
		labels:      labelsHash,
	}, nil
}

func (g *Generic) Annotations() string {
	return g.annotations
}

func (g *Generic) Labels() string {
	return g.labels
}

func (g *Generic) Equals(content ResourceContent) bool {
	return g.labels == content.Labels() &&
		g.annotations == content.Annotations() &&
		g.Hash() == content.Hash()
}

func (g *Generic) Hash() string {
	return ""
}
