package template

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type RenderOption func(*template.Template) *template.Template

func WithOption(option string) RenderOption {
	return func(t *template.Template) *template.Template {
		return t.Option(option)
	}
}

func RenderTemplate(values any, tpl string, options ...RenderOption) (*unstructured.Unstructured, error) {
	if options == nil {
		options = []RenderOption{WithOption("missingkey=error")}
	}

	rdr, err := renderString(values, tpl, options...)
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

func renderString(values any, tpl string, tplOptions ...RenderOption) (string, error) {
	t := template.New("tpl").Delims("[[", "]]").Funcs(template.FuncMap{
		"b64enc": b64enc,
	})
	for _, option := range tplOptions {
		t = option(t)
	}
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

func b64enc(in any) string {
	s, _ := in.(string)
	return base64.StdEncoding.EncodeToString([]byte(s))
}
