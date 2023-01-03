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

func ParseResources(values *TemplateValues, resources []naisiov1.Resource) ([]*unstructured.Unstructured, error) {
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

func ExtractValues(annotations map[string]string) map[string]string {
	values := make(map[string]string)
	for key, value := range annotations {
		kp := strings.Split(key, "replicator.nais.io/")
		if len(kp) != 2 {
			continue
		}
		values[kp[1]] = value
	}
	return values
}

func LoadSecrets(ctx context.Context, c client.Client, rc *naisiov1.ReplicationConfig) (map[string]string, error) {
	values := make(map[string]string)
	for _, s := range rc.Spec.ValueSecrets {

		var secret v1.Secret
		if err := c.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: os.Getenv("NAMESPACE")}, &secret); err != nil {
			return nil, err
		}

		for k, v := range secret.Data {
			values[k] = string(v)
		}
	}
	return values, nil
}
