package processor

import "k8s.io/client-go/rest"

func Apply() {
	_, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
}
