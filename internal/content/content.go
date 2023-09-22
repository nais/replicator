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
	CompareTo(existing *unstructured.Unstructured) (bool, error)
}

type Spec struct {
	contentHash string
	contentType string
	data        *unstructured.Unstructured
}

type Data struct {
	contentHash string
	contentType string
	data        *unstructured.Unstructured
}

type StringData struct {
	contentHash string
	contentType string
	data        *unstructured.Unstructured
}

func GetContentHash(data *unstructured.Unstructured) (ResourceContent, error) {
	var resourceContent ResourceContent
	var content map[string]interface{}
	var err error
	var hash string

	switch {
	case data.UnstructuredContent()[SpecContent] != nil:
		content, err = getContent(data, SpecContent)
		hash, err = toHash(content)
		if err != nil {
			return nil, err
		}
		resourceContent = &Spec{
			contentHash: hash,
			contentType: SpecContent,
			data:        data,
		}
	case data.UnstructuredContent()[DataContent] != nil:
		content, err = getContent(data, DataContent)
		hash, err = toHash(content)
		if err != nil {
			return nil, err
		}
		resourceContent = &Data{
			contentHash: hash,
			contentType: DataContent,
			data:        data,
		}
	case data.UnstructuredContent()[StringDataContent] != nil:
		content, err = getContent(data, StringDataContent)
		encodedValues := copyToEncodedValues(content)
		hash, err = toHash(encodedValues)
		resourceContent = &StringData{
			contentHash: hash,
			contentType: StringDataContent,
			data:        data,
		}
	default:
		_, err = getContent(data, UnknownContent)
	}
	return resourceContent, err
}

func (s *Spec) CompareTo(existing *unstructured.Unstructured) (bool, error) {
	existingData, err := getContent(existing, s.contentType)
	if err != nil {
		return false, err
	}
	existingHash, err := toHash(existingData)
	if err != nil {
		return false, err
	}
	return hasChanged(s.contentHash, existingHash), nil
}

func (d *Data) CompareTo(existing *unstructured.Unstructured) (bool, error) {
	existingData, err := getContent(existing, d.contentType)
	if err != nil {
		return false, err
	}
	existingHash, err := toHash(existingData)
	if err != nil {
		return false, err
	}
	return hasChanged(d.contentHash, existingHash), nil
}

func (s *StringData) CompareTo(existing *unstructured.Unstructured) (bool, error) {
	existingData, err := getContent(existing, DataContent)
	if err != nil {
		return false, err
	}
	existingHash, err := toHash(existingData)
	if err != nil {
		return false, err
	}
	return hasChanged(s.contentHash, existingHash), nil
}

func hasChanged(resourceHash, existingHash string) bool {
	return resourceHash != existingHash
}

func toHash(input map[string]interface{}) (string, error) {
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

func getContent(data *unstructured.Unstructured, contentType string) (map[string]interface{}, error) {
	if data.UnstructuredContent()[contentType] == nil {
		return nil, fmt.Errorf("content type %q not found with data %v", contentType, data.UnstructuredContent())
	}

	return data.UnstructuredContent()[contentType].(map[string]interface{}), nil
}
