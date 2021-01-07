package helpers

/*
func getK8SDynamicClientSet(kubeconfig []byte, contextName string) (dynamic.Interface, error) {
	var clientConfig *rest.Config
	var err error
	if len(kubeconfig) == 0 {
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			err = errors.Wrap(err, "unable to load in-cluster kubeconfig")
			logrus.Error(err)
			return nil, err
		}
	} else {
		config, err := clientcmd.Load(kubeconfig)
		if err != nil {
			err = errors.Wrap(err, "unable to load kubeconfig")
			logrus.Error(err)
			return nil, err
		}
		if contextName != "" {
			config.CurrentContext = contextName
		}
		clientConfig, err = clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			err = errors.Wrap(err, "unable to create client config from config")
			logrus.Error(err)
			return nil, err
		}
	}
	clientConfig.Timeout = 2 * time.Second
	clientset, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		err = errors.Wrap(err, "unable to create client set")
		logrus.Error(err)
		return nil, err
	}
	return clientset, nil
}
*/

// func getIstioClient(kubeconfig []byte, contextName string) (*versionedclient.Clientset, error) {
// 	var clientConfig *rest.Config
// 	var err error
// 	if len(kubeconfig) == 0 {
// 		clientConfig, err = rest.InClusterConfig()
// 		if err != nil {
// 			err = errors.Wrap(err, "unable to load in-cluster kubeconfig")
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 	} else {
// 		config, err := clientcmd.Load(kubeconfig)
// 		if err != nil {
// 			err = errors.Wrap(err, "unable to load kubeconfig")
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		if contextName != "" {
// 			config.CurrentContext = contextName
// 		}
// 		clientConfig, err = clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
// 		if err != nil {
// 			err = errors.Wrap(err, "unable to create client config from config")
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 	}
// 	clientConfig.Timeout = 2 * time.Second
// 	clientset, err := versionedclient.NewForConfig(clientConfig)
// 	if err != nil {
// 		err = errors.Wrap(err, "unable to create client set")
// 		logrus.Error(err)
// 		return nil, err
// 	}
// 	return clientset, nil
// }

/*
// ScanIstio - Runs a quick scan on kubernetes to find out the version of service meshes deployed
func ScanIstio(kubeconfig []byte, contextName string) (map[string]string, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}

	dclientset, err := getK8SDynamicClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}

	res := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}
	namespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "unable to get the list of namespaces")
		logrus.Error(err)
		return nil, err
	}
	for _, ns := range namespacelist.Items {
		logrus.Debugf("Listing deployments in namespace %q", ns.GetName())
		data1, err := dclientset.Resource(res).Namespace(ns.GetName()).List(metav1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "unable to get vs in the %s namespace", ns)
			logrus.Error(err)
			return nil, err
		}

		if data1.IsList() {
			err = data1.EachListItem(func(obj runtime.Object) error {
				dataL, _ := obj.(*unstructured.Unstructured)
				dmap := dataL.UnstructuredContent()
				spec1 := dmap["spec"].(map[string]interface{})
				hosts, _ := spec1["hosts"].([]string)
				logrus.Infof("hosts: %v", hosts)
				hosts, _ := spec1["hosts"].([]string)
				return nil
			})
		}
	}
	logrus.Debugf("Derived mesh versions: %s", result)
	return result, nil
}
*/

// ScanIstio - lists VirtualServices
func ScanIstio(kubeconfig []byte, contextName string) (map[string]string, error) {
	return make(map[string]string), nil
	// clientset, err := getK8SClientSet(kubeconfig, contextName)
	// if err != nil {
	// 	return nil, err
	// }

	// dclientset, err := getIstioClient(kubeconfig, contextName)
	// if err != nil {
	// 	return nil, err
	// }
	// result := map[string]string{}

	// namespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	// if err != nil {
	// 	err = errors.Wrap(err, "unable to get the list of namespaces")
	// 	logrus.Error(err)
	// 	return nil, err
	// }
	// for _, ns := range namespacelist.Items {
	// 	logrus.Debugf("Listing deployments in namespace %q", ns.GetName())
	// 	vsList, err := dclientset.NetworkingV1alpha3().VirtualServices(ns.GetName()).List(metav1.ListOptions{})
	// 	if err != nil {
	// 		err = errors.Wrapf(err, "unable to get vs in the %s namespace", ns.GetName())
	// 		logrus.Error(err)
	// 		return nil, err
	// 	}
	// 	for i, vs := range vsList.Items {
	// 		logrus.Infof("Index: %d VirtualService Hosts: %+v\n", i, vs.Spec.GetHosts())
	// 		logrus.Infof("Index: %d VirtualService HTTP: %+v\n", i, vs.Spec.GetHttp())
	// 		// TODO: only take uri -> exact, bcoz regexes have a lot of possibilities
	// 	}
	// }
	// return result, nil
}
