package kubernetes

import (
	"context"
	"fmt"
	"tgs-automation/internal/log"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type ServiceApplier struct{}

func (s *ServiceApplier) Apply(clientset *kubernetes.Clientset, doc string) error {
	var service corev1.Service
	if err := yaml.Unmarshal([]byte(doc), &service); err != nil {
		return fmt.Errorf("failed to unmarshal Service: %v", err)
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := clientset.CoreV1().Services(service.Namespace).Create(context.TODO(), &service, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				log.LogInfo(fmt.Sprintf("Service %s already exists. Updating it.", service.Name))
				existing, getErr := clientset.CoreV1().Services(service.Namespace).Get(context.TODO(), service.Name, metav1.GetOptions{})
				if getErr != nil {
					return fmt.Errorf("failed to get existing Service: %v", getErr)
				}
				service.ResourceVersion = existing.ResourceVersion
				_, updateErr := clientset.CoreV1().Services(service.Namespace).Update(context.TODO(), &service, metav1.UpdateOptions{})
				if updateErr != nil {
					return fmt.Errorf("failed to update Service: %v", updateErr)
				}
				log.LogInfo(fmt.Sprintf("Service %s updated successfully\n", service.Name))
				return nil
			}
			return fmt.Errorf("failed to create Service: %v", err)
		}
		log.LogInfo(fmt.Sprintf("Service %s created successfully\n", service.Name))
		return nil
	})
}
