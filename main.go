package main

import (
	"flag"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	klog "k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)

	var kubeconfig string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := buildConfig(kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	client := clientset.NewForConfigOrDie(config)

	preferredResources, err := client.Discovery().ServerPreferredResources()
	if err != nil {
		klog.Errorln("Failed to do client.Discovery().ServerPreferredResources()")
	}

	klog.Infoln("client.Discovery().ServerPreferredResources():")
	for _, r := range preferredResources {
		klog.Infoln("========", r.GroupVersion)
		for _, res := range r.APIResources {
			klog.Infoln(res.Kind)
		}
	}

}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
