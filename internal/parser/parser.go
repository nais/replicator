package parser

import (
	"fmt"
	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/internal/util"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// from CRD
type Values struct {
	Secrets []Secret `json:"secrets,omitempty"`
}

type Secret struct {
	Name string
}

// for Templating
type TemplateValues struct {
	Values map[string]string
}

func Resources(values *TemplateValues, resources []naisiov1.Resource) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured
	for _, r := range resources {
		resource, err := util.RenderTemplate(values, r.Template)
		if err != nil {
			return nil, err
		}
		objects = append(objects, resource)
	}
	return objects, nil
}

func ParseAnnotations(annotations map[string]string, values *TemplateValues) error {
	for key, value := range annotations {
		kp := strings.Split(key, "replicator.nais.io/")
		if len(kp) < 2 {
			fmt.Printf("invalid annotation: %s", key)
			continue
		}
		values.Values[kp[1]] = value
	}
	return nil
}
