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
	if change, err := hasChanged(s.data, existing, s.ContentType); err != nil {
		return false, err
	} else if change {
		return true, nil
	}
	return false, nil
}

func (d *Data) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	if change, err := hasChanged(d.data, existing, d.ContentType); err != nil {
		return false, err
	} else if change {
		return true, nil
	}
	return false, nil
}

func (s *StringData) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	if change, err := hasChanged(s.data, existing, s.ContentType); err != nil {
		return false, err
	} else if change {
		return true, nil
	}
	return false, nil
}

func (u *Unknown) ContentHasChanged(existing *unstructured.Unstructured) (bool, error) {
	if change, err := hasChanged(u.data, existing, u.ContentType); err != nil {
		return false, err
	} else if change {
		return true, nil
	}
	return false, nil
}

func hasChanged(resource, existing *unstructured.Unstructured, contentType string) (bool, error) {
	existingDataHash, err := hash(existing, contentType)
	if err != nil {
		return false, err
	}

	resourceDataHash, err := hash(resource, contentType)
	if err != nil {
		return false, err
	}

	return existingDataHash != resourceDataHash, nil
}

func hash(input *unstructured.Unstructured, contentType string) (string, error) {
	hash, err := hashstructure.Hash(input.UnstructuredContent()[contentType], hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}
