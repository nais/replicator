package resources

import (
	"fmt"
	"github.com/mitchellh/hashstructure/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func HasChanged(existing, inputResource *unstructured.Unstructured) (bool, error) {
	existingSpec, err := spec(existing.Object)
	if err != nil {
		return false, err
	}

	resourceSpec, err := spec(inputResource.Object)
	if err != nil {
		return false, err
	}

	return existingSpec != resourceSpec, nil
}

func spec(data map[string]interface{}) (string, error) {
	return hash(data["spec"])
}

func hash(data interface{}) (string, error) {
	hash, err := hashstructure.Hash(data, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}
