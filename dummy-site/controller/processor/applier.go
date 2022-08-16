package processor

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type kubernetesApplier struct {
	clientset *kubernetes.Clientset
}

func NewApplier() (*kubernetesApplier, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubernetesApplier{clientset: clientset}, nil
}

type status string

const (
	StatusWorking status = "WORKING"
	StatusDone    status = "DONE"
	StatusError   status = "ERROR"
	namespace            = "default"
)

type applier interface {
	apply(*manifests) error
	cleanupResources(*manifests) error
	sleepDuration() time.Duration
}

func (a *kubernetesApplier) ApplyDummySite(website string, url string, term chan bool, status chan status) error {
	body, err := NewScraper().Scrape(url)
	if err != nil {
		return err
	}
	mr, err := GenerateManifests(website, string(body))
	if err != nil {
		return err
	}
	m, err := a.readManifests(mr)
	if err != nil {
		return err
	}

	return applyUntilDestroyed(a, m, term, status)
}

const sleepDuration = 10 * time.Second

func applyUntilDestroyed(a applier, m *manifests, term chan bool, status chan status) error {
	defer close(status)

	terminating := false
	for {
		select {
		case <-term:
			terminating = true
		default:
		}

		var err error

		if terminating {
			err = a.cleanupResources(m)
			if err == nil {
				status <- StatusDone
				return nil
			}
		} else {
			err = a.apply(m)
		}

		if err != nil {
			status <- StatusError
		} else {
			status <- StatusWorking
		}

		time.Sleep(a.sleepDuration())
	}
}

type ManifestReaders struct {
	deploymentReader io.Reader
	serviceReader    io.Reader
	ingressReader    io.Reader
}

type manifests struct {
	deployment *appsv1.Deployment
	service    *v1.Service
	ingress    *networkingv1.Ingress
	configMap  *v1.ConfigMap
}

//go:embed manifests/configmap.yml
var nginxConfigMap []byte

func (a *kubernetesApplier) readManifests(m *ManifestReaders) (*manifests, error) {
	deployment := appsv1.Deployment{}
	err := yaml.NewYAMLOrJSONDecoder(m.deploymentReader, 1).Decode(&deployment)
	if err != nil {
		return nil, err
	}

	service := v1.Service{}
	err = yaml.NewYAMLOrJSONDecoder(m.serviceReader, 1).Decode(&service)
	if err != nil {
		return nil, err
	}

	ingress := networkingv1.Ingress{}
	err = yaml.NewYAMLOrJSONDecoder(m.ingressReader, 1).Decode(&ingress)
	if err != nil {
		return nil, err
	}

	configMap := v1.ConfigMap{}
	err = yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(nginxConfigMap), 1).Decode(&configMap)
	if err != nil {
		return nil, err
	}

	return &manifests{
		deployment: &deployment,
		service:    &service,
		ingress:    &ingress,
		configMap:  &configMap,
	}, nil
}

func (a *kubernetesApplier) apply(m *manifests) error {
	// TODO create dummy-sites namespace if not exists and use it
	_, err := a.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), m.configMap.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := a.clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), m.configMap, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("ConfigMap %s created!\n", m.configMap.Name)
	} else {
		_, err := a.clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), m.configMap, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("ConfigMap %s updates!\n", m.configMap.Name)
	}

	_, err = a.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), m.deployment.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := a.clientset.AppsV1().Deployments(namespace).Create(context.TODO(), m.deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Deployment %s created!\n", m.deployment.Name)
	} else {
		_, err := a.clientset.AppsV1().Deployments(namespace).Update(context.TODO(), m.deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Deployment %s updated!\n", m.deployment.Name)
	}

	_, err = a.clientset.CoreV1().Services(namespace).Get(context.TODO(), m.service.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := a.clientset.CoreV1().Services(namespace).Create(context.TODO(), m.service, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Service %s created!\n", m.service.Name)
	} else {
		_, err := a.clientset.CoreV1().Services(namespace).Update(context.TODO(), m.service, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Service %s updated!\n", m.service.Name)
	}

	_, err = a.clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), m.ingress.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := a.clientset.NetworkingV1().Ingresses(namespace).Create(context.TODO(), m.ingress, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Ingress %s created!\n", m.service.Name)
	} else {
		_, err := a.clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), m.ingress, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Ingress %s updated!\n", m.service.Name)
	}
	return nil
}

// TODO test manually by calling directly
func (a *kubernetesApplier) cleanupResources(m *manifests) error {
	_, err := a.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), m.deployment.Name, metav1.GetOptions{})
	if err == nil {
		err := a.clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), m.deployment.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		fmt.Printf("Deployment %s deleted!\n", m.deployment.Name)
	} else if !errors.IsNotFound(err) {
		return err
	}

	_, err = a.clientset.CoreV1().Services(namespace).Get(context.TODO(), m.service.Name, metav1.GetOptions{})
	if err == nil {
		err := a.clientset.CoreV1().Services(namespace).Delete(context.TODO(), m.service.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	_, err = a.clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), m.ingress.Name, metav1.GetOptions{})
	if err == nil {
		err := a.clientset.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), m.ingress.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	} else if !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (a *kubernetesApplier) sleepDuration() time.Duration {
	return sleepDuration
}
