package parser

import (
	"bytes"
	"fmt"
	"text/template"

	naisiov1 "nais/replicator/api/v1"

	"gopkg.in/yaml.v2"
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
		resource, err := RenderTemplate(values, r.Template)
		if err != nil {
			return nil, err
		}

		fmt.Printf("resource Kind: %s\n", resource.GetKind())
		fmt.Printf("resource Namespace: %s\n", resource.GetNamespace())
		objects = append(objects, resource)
	}
	return objects, nil
}

func RenderTemplate(values any, tpl string) (*unstructured.Unstructured, error) {
	rdr, err := renderString(values, tpl)
	if err != nil {
		return nil, err
	}

	var v any
	if err := yaml.Unmarshal([]byte(rdr), &v); err != nil {
		return nil, err
	}
	v = repairMapAny(v)

	u := &unstructured.Unstructured{
		Object: v.(map[string]interface{}),
	}

	return u, nil
}

func repairMapAny(v any) any {
	switch t := v.(type) {
	case []any:
		for i, v := range t {
			t[i] = repairMapAny(v)
		}
	case map[any]any:
		nm := make(map[string]any)
		for k, v := range t {
			nm[k.(string)] = repairMapAny(v)
		}
		return nm
	}
	return v
}

func renderString(values any, tpl string) (string, error) {
	t := template.New("tpl")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, values); err != nil {
		return "", err
	}
	return buf.String(), nil
}
