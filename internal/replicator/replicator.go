package replicator

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	naisiov1 "nais/replicator/api/v1"
	"nais/replicator/internal/template"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// for Templating
type TemplateValues struct {
	Values map[string]string
}

func LoadValues(ctx context.Context, c client.Client, rc *naisiov1.ReplicatorConfiguration) (*TemplateValues, error) {
	values := &TemplateValues{
		Values: map[string]string{},
	}

	if err := loadSecrets(ctx, c, rc, values); err != nil {
		return nil, err
	}

	if err := loadConfigMaps(ctx, c, rc, values); err != nil {
		return nil, err
	}

	return values, nil
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

func AddAnnotations(annotations map[string]string, values *TemplateValues) error {
	for key, value := range annotations {
		kp := strings.Split(key, "replicator.nais.io/")
		if len(kp) != 2 {
			continue
		}
		values.Values[kp[1]] = value
	}
	return nil
}

func loadSecrets(ctx context.Context, c client.Client, rc *naisiov1.ReplicatorConfiguration, values *TemplateValues) error {
	for _, s := range rc.Spec.Values.Secrets {

		var secret v1.Secret
		if err := c.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, &secret); err != nil {
			return err
		}

		for k, v := range secret.Data {
			values.Values[k] = string(v)
		}
	}
	return nil
}

func loadConfigMaps(ctx context.Context, c client.Client, rc *naisiov1.ReplicatorConfiguration, values *TemplateValues) error {
	for _, s := range rc.Spec.Values.ConfigMaps {

		var configMap v1.ConfigMap
		if err := c.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, &configMap); err != nil {
			return err
		}

		for k, v := range configMap.Data {
			values.Values[k] = v
		}
	}
	return nil
}
