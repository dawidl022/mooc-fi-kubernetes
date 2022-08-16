package gateway

import (
	"context"
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/dummy-site/controller/processor"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Controller struct {
	jobs      map[string]*applierJob
	applier   *processor.KubernetesApplier
	crdClient *rest.RESTClient
}

func NewController() (*Controller, error) {
	applier, err := processor.NewApplier()
	if err != nil {
		return nil, err
	}

	crdClient, err := initCRDClient()
	if err != nil {
		return nil, err
	}

	return &Controller{
		jobs:      make(map[string]*applierJob),
		applier:   applier,
		crdClient: crdClient,
	}, nil
}

func (c *Controller) Start() error {

	for {
		dummySites := DummySiteList{}
		c.crdClient.Get().Namespace("default").Resource("dummysites").Do(context.TODO()).Into(&dummySites)

		c.addOrUpdateDummySites(&dummySites)
		c.garbageCollectDummySites(&dummySites)
		time.Sleep(5 * time.Second)
	}
}

func (c *Controller) addOrUpdateDummySites(dummySites *DummySiteList) {
	for _, site := range dummySites.Items {
		job := c.jobs[site.Name]
		if job == nil || job.url != site.Spec.WebsiteUrl {
			c.jobs[site.Name] = newApplierJob(c.applier, site.Name, site.Spec.WebsiteUrl)
		}
	}
}

func (c *Controller) garbageCollectDummySites(dummySites *DummySiteList) {
	clusterSiteNames := make(map[string]struct{})
	for _, site := range dummySites.Items {
		clusterSiteNames[site.Name] = struct{}{}
	}

	for _, job := range c.jobs {
		if _, inCluster := clusterSiteNames[job.website]; !inCluster {
			if job.status == processor.StatusDone {
				delete(c.jobs, job.website)
			} else {
				go func() { job.term <- true }()
			}
		}
	}
}

func initCRDClient() (*rest.RESTClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}
	crdConfig.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}
	return exampleRestClient, nil
}
