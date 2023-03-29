package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

type SampleResource struct {
	runtime.TypeMeta   `json:",inline"`
	runtime.ObjectMeta `json:"metadata,omitempty"`

	Spec   SampleResourceSpec   `json:"spec,omitempty"`
	Status SampleResourceStatus `json:"status,omitempty"`
}

type SampleResourceSpec struct {
	Message string `json:"message,omitempty"`
}

type SampleResourceStatus struct {
	Processed bool `json:"processed,omitempty"`
}

func main() {
	// Set up logging.
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logging: %w", err))
	}
	defer zapLog.Sync()
	log := zapr.NewLogger(zapLog)

	// Set up command-line flags.
	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file")
	flag.Parse()

	// Load the Kubernetes configuration.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Error(err, "failed to load kubeconfig")
		panic(err)
	}

	// Create the Kubernetes client.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "failed to create Kubernetes clientset")
		panic(err)
	}

	// Create the controller object.
	controller := &Controller{
		clientset: clientset,
		workqueue: workqueue.NewNamed("sample-resource"),
		informer: cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
					return clientset.SampleResource().List(context.Background(), metav1.ListOptions{})
				},
				WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
					return clientset.SampleResource().Watch(context.Background(), metav1.ListOptions{})
				},
			},
			&SampleResource{},
			0,
			cache.Indexers{},
		),
		log: log,
	}

	// Start the controller.
	controller.Run()
}

type Controller struct {
	clientset kubernetes.Interface
	workqueue workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
	log       logr.Logger
}

func (c *Controller) Run()
