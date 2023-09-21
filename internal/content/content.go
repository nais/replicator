package content

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	SpecContent       = "spec"
	DataContent       = "data"
	StringDataContent = "stringData"
	UnknownContent    = "unknown"
)

type ResourceContent interface {
	ContentHasChanged(existing *unstructured.Unstructured) (bool, error)
}

type Spec struct {
	data        *unstructured.Unstructured
	ContentType string
}

type Data struct {
	data        *unstructured.Unstructured
	ContentType string
}

type StringData struct {
	data        *unstructured.Unstructured
	ContentType string
}

type Unknown struct {
	data        *unstructured.Unstructured
	ContentType string
}

func GetHash(data *unstructured.Unstructured) (ResourceContent, error) {
	var resourceContent ResourceContent
	var err error

	switch {
	case data.UnstructuredContent()[SpecContent] != nil:
		resourceContent = &Spec{
			data:        data,
			ContentType: SpecContent,
		}
	case data.UnstructuredContent()[DataContent] != nil:
		resourceContent = &Data{
			data:        data,
			ContentType: DataContent,
		}
	case data.UnstructuredContent()[StringDataContent] != nil:
		resourceContent = &StringData{
			data:        data,
			ContentType: StringDataContent,
		}
	default:
		resourceContent = &Unknown{
			data:        data,
			ContentType: UnknownContent,
		}
		err = fmt.Errorf("resource content type not found: %v", data.UnstructuredContent())
	}
	return resourceContent, err
}

func (s *Spec) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	data := toContent(s.data, s.ContentType)
	existingData := toContent(existing, s.ContentType)

	change, err := hasChanged(data, existingData)
	if err != nil {
		return false, err
	}
	return change, nil
}

func (d *Data) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	data := toContent(d.data, d.ContentType)
	existingData := toContent(existing, d.ContentType)

	change, err := hasChanged(data, existingData)
	if err != nil {
		return false, err
	}
	return change, nil
}

func (s *StringData) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	data := toContent(s.data, s.ContentType)
	existingData := toContent(existing, DataContent)

	if len(data) != len(existingData) {
		return true, nil
	}

	change, err := hasChanged(copyToEncodedValues(data), existingData)
	if err != nil {
		return false, err
	}
	return change, nil
}

func (u *Unknown) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	data := toContent(u.data, u.ContentType)
	existingData := toContent(existing, u.ContentType)

	change, err := hasChanged(data, existingData)
	if err != nil {
		return false, err
	}
	return change, nil
}

func hasChanged(resource, existing map[string]interface{}) (bool, error) {
	existingDataHash, err := hash(existing)
	if err != nil {
		return false, err
	}

	resourceDataHash, err := hash(resource)
	if err != nil {
		return false, err
	}

	return existingDataHash != resourceDataHash, nil
}

func hash(input map[string]interface{}) (string, error) {
	hash, err := hashstructure.Hash(input, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}

func copyToEncodedValues(resources map[string]interface{}) map[string]interface{} {
	outputs := make(map[string]interface{}, len(resources))

	for k, v := range resources {
		outputs[k] = b64.StdEncoding.EncodeToString([]byte(v.(string)))
	}
	return outputs
}

func toContent(data *unstructured.Unstructured, contentType string) map[string]interface{} {
	return data.UnstructuredContent()[contentType].(map[string]interface{})
}
