package replicator

import (
	"context"
	"os"
	"strings"

	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/internal/template"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TemplateValues struct {
	Values map[string]string
}

func RenderResources(values *TemplateValues, resources []naisiov1.Resource) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured
	for _, r := range resources {
		resource, err := template.RenderTemplate(values, r.Template)
		if err != nil {
			return nil, err
		}
		objects = append(objects, resource)
	}
	return objects, nil
}

func ExtractValues(namespace v1.Namespace, namespaceValues naisiov1.Namespace) map[string]string {
	labels := filter(namespace.Labels, namespaceValues.Labels)
	return Merge(labels, filter(namespace.Annotations, namespaceValues.Annotations))
}

func Merge(a, b map[string]string) map[string]string {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}

func filter(m map[string]string, keys []string) map[string]string {
	values := make(map[string]string)
	for _, k := range keys {
		if v, ok := m[k]; ok {
			values[normalizeKey(k)] = v
		}
	}
	return values
}

func normalizeKey(key string) string {
	if strings.Contains(key, "/") {
		return strings.Split(key, "/")[1]
	}
	return key
}

func LoadSecrets(ctx context.Context, c client.Client, rc *naisiov1.ReplicationConfig) (map[string]string, error) {
	values := make(map[string]string)
	for _, s := range rc.Spec.TemplateValues.Secrets {

		var secret v1.Secret
		if err := c.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: os.Getenv("POD_NAMESPACE")}, &secret); err != nil {
			return nil, err
		}

		for k, v := range secret.Data {
			values[k] = string(v)
		}
	}
	return values, nil
}
