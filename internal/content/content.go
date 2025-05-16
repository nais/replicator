package content

import (
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	SpecContent       = "spec"
	DataContent       = "data"
	StringDataContent = "stringData"
)

type ResourceContent interface {
	Annotations() string
	Labels() string
	Equals(content ResourceContent) bool
	Hash() string
}

func Get(data *unstructured.Unstructured) (ResourceContent, error) {
	switch {
	case data.UnstructuredContent()[SpecContent] != nil:
		return NewSpec(data)
	case data.UnstructuredContent()[DataContent] != nil:
		return NewData(data)
	case data.UnstructuredContent()[StringDataContent] != nil:
		return NewStringData(data)
	default:
		return NewGeneric(data)
	}
}

func toHash(input any) (string, error) {
	hash, err := hashstructure.Hash(input, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}

func getContent(data *unstructured.Unstructured, contentType string) (map[string]interface{}, error) {
	content := data.UnstructuredContent()[contentType]
	if content == nil {
		return nil, fmt.Errorf("content type %q not found with data %v", contentType, data.UnstructuredContent())
	}
	return content.(map[string]interface{}), nil
}
