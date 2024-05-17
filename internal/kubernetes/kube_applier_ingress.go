package kubernetes

import (
	"context"
	"fmt"

	"tgs-automation/internal/log"

	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type IngressApplier struct{}

func (i *IngressApplier) Apply(clientset *kubernetes.Clientset, doc string) error {
	var ingress networkingv1.Ingress
	if err := yaml.Unmarshal([]byte(doc), &ingress); err != nil {
		return fmt.Errorf("failed to unmarshal Ingress: %v", err)
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := clientset.NetworkingV1().Ingresses(ingress.Namespace).Create(context.TODO(), &ingress, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				log.LogInfo(fmt.Sprintf("Ingress %s already exists. Updating it.", ingress.Name))
				existing, getErr := clientset.NetworkingV1().Ingresses(ingress.Namespace).Get(context.TODO(), ingress.Name, metav1.GetOptions{})
				if getErr != nil {
					return fmt.Errorf("failed to get existing Ingress: %v", getErr)
				}
				ingress.ResourceVersion = existing.ResourceVersion
				_, updateErr := clientset.NetworkingV1().Ingresses(ingress.Namespace).Update(context.TODO(), &ingress, metav1.UpdateOptions{})
				if updateErr != nil {
					return fmt.Errorf("failed to update Ingress: %v", updateErr)
				}
				log.LogInfo(fmt.Sprintf("Ingress %s updated successfully\n", ingress.Name))
				return nil
			}
			return fmt.Errorf("failed to create Ingress: %v", err)
		}
		log.LogInfo(fmt.Sprintf("Ingress %s created successfully\n", ingress.Name))
		return nil
	})
}
