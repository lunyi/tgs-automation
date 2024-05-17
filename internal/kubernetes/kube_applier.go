package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

func ApplyYamlDocument(clientset *kubernetes.Clientset, doc string) error {
	decoder := yaml.NewYAMLToJSONDecoder(strings.NewReader(doc))
	var rawObj map[string]interface{}
	if err := decoder.Decode(&rawObj); err != nil {
		return fmt.Errorf("failed to decode YAML document: %v", err)
	}

	kind, ok := rawObj["kind"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'kind' in YAML document")
	}

	applier, err := getResourceApplier(kind)
	if err != nil {
		return err
	}

	err = applier.Apply(clientset, doc)
	if err != nil {
		return fmt.Errorf("failed to apply %s resource: %v", kind, err)
	}

	return nil
}

type ResourceApplier interface {
	Apply(clientset *kubernetes.Clientset, doc string) error
}

var appliers = map[string]func() (ResourceApplier, error){
	"Deployment": func() (ResourceApplier, error) { return &DeploymentApplier{}, nil },
	"Service":    func() (ResourceApplier, error) { return &ServiceApplier{}, nil },
	"Ingress":    func() (ResourceApplier, error) { return &IngressApplier{}, nil },
}

func getResourceApplier(kind string) (ResourceApplier, error) {
	if applier, found := appliers[kind]; found {
		return applier()
	}
	return nil, fmt.Errorf("unsupported resource kind: %s", kind)
}
