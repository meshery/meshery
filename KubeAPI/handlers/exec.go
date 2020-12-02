package handlers

import (
	"github.com/layer5io/meshery/KubeAPI/structs"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/polymorphichelpers"
	"log"
	"os"
	"path/filepath"
)

var clientset *kubernetes.Clientset

func init() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err = kubernetes.NewForConfig(config)
}

func ExecFunction(resourceName string, command []string, namespace string) error {

	var err error

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	options := &structs.ExecOptions{
		StreamOptions: structs.StreamOptions{
			IOStreams: ioStreams,
		},

		Executor: &structs.DefaultRemoteExecutor{},
	}

	options.ResourceName = resourceName
	options.Command = command
	options.Namespace = namespace
	options.EnforceNamespace = false
	options.ExecutablePodFn = polymorphichelpers.AttachablePodForObjectFn
	options.Builder = f.NewBuilder
	options.Config, err = f.ToRESTConfig()
	if err != nil {
		return err
	}
	options.RestClientGetter = f
	options.PodClient = clientset.CoreV1()

	cmdutil.CheckErr(options.Validate())
	cmdutil.CheckErr(options.Run())

	return nil
}
