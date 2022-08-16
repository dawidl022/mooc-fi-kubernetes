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

func readManifests(m ManifestReaders) (*manifests, error) {
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

func ApplyUntilDestroyed() {
	// TODO
}

func Apply(mr ManifestReaders) {
	m, err := readManifests(mr)
	if err != nil {
		panic(err.Error())
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// TODO create dummy-sites namespace if not exists and use it
		namespace := "default"
		_, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), m.configMap.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), m.configMap, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("ConfigMap %s created!\n", m.configMap.Name)
		} else {
			_, err := clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), m.configMap, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("ConfigMap %s updates!\n", m.configMap.Name)
		}

		_, err = clientset.AppsV1().Deployments(namespace).Get(context.TODO(), m.deployment.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), m.deployment, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Deployment %s created!\n", m.deployment.Name)
		} else {
			_, err := clientset.AppsV1().Deployments(namespace).Update(context.TODO(), m.deployment, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Deployment %s updated!\n", m.deployment.Name)
		}

		_, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), m.service.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), m.service, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Service %s created!\n", m.service.Name)
		} else {
			_, err := clientset.CoreV1().Services(namespace).Update(context.TODO(), m.service, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Service %s updated!\n", m.service.Name)
		}

		_, err = clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), m.ingress.Name, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, err := clientset.NetworkingV1().Ingresses(namespace).Create(context.TODO(), m.ingress, metav1.CreateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Ingress %s created!\n", m.service.Name)
		} else {
			_, err := clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), m.ingress, metav1.UpdateOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Ingress %s updated!\n", m.service.Name)
		}

		time.Sleep(10 * time.Second)
	}
}
