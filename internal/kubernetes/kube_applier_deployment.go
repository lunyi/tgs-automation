package kubernetes

import (
	"context"
	"fmt"

	"tgs-automation/internal/log"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

type DeploymentApplier struct{}

func (d *DeploymentApplier) Apply(clientset *kubernetes.Clientset, doc string) error {
	var deployment appsv1.Deployment
	if err := yaml.Unmarshal([]byte(doc), &deployment); err != nil {
		return fmt.Errorf("failed to unmarshal Deployment: %v", err)
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err := clientset.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
		if err != nil {
			if apierrors.IsAlreadyExists(err) {
				log.LogInfo(fmt.Sprintf("Deployment %s already exists. Updating it.", deployment.Name))
				existing, getErr := clientset.AppsV1().Deployments(deployment.Namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
				if getErr != nil {
					return fmt.Errorf("failed to get existing Deployment: %v", getErr)
				}
				deployment.ResourceVersion = existing.ResourceVersion
				_, updateErr := clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
				if updateErr != nil {
					return fmt.Errorf("failed to update Deployment: %v", updateErr)
				}
				log.LogInfo(fmt.Sprintf("Deployment %s updated successfully\n", deployment.Name))
				return nil
			}
			return fmt.Errorf("failed to create Deployment: %v", err)
		}
		log.LogInfo(fmt.Sprintf("Deployment %s created successfully\n", deployment.Name))
		return nil
	})
}
