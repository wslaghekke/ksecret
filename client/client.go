package client

import (
	"fmt"
	"os"
	"path/filepath"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"k8s.io/api/core/v1"
)

type secretData struct {
	data []secretDataItem
}

type secretDataItem struct {
	Key string
	Value string
}

func CreateClient(kubeconfig string, namespace string) (*kubernetes.Clientset, *string, error) {
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homeDir(), ".kube", "config")
	}

	var clientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{ClusterInfo: api.Cluster{Server: ""}})

	// use the current context in kubeconfig
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	if namespace == "" {
		// fallback to namespace from kubeconfig
		configNamespace, _, err := clientConfig.Namespace()
		if err != nil {
			return nil, nil, err
		}
		namespace = configNamespace
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return clientset, &namespace, nil
}
func EditSecret(clientset *kubernetes.Clientset, namespace *string, secretName string) {
	// Examples for error handling:
	// - Use helper functions like e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	secret, err := clientset.CoreV1().Secrets(*namespace).Get(secretName, metav1.GetOptions{})
	checkErr(err)

	bytes, err := secretDataToYaml(secret)
	checkErr(err)

	tempFile, err := ioutil.TempFile("", "")
	checkErr(err)

	tempFile.Write(bytes)
	tempFile.Close()
	tempFileName := tempFile.Name()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	editorCmd := exec.Command(editor, tempFileName)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	checkErr(editorCmd.Run())

	readBytes, err := ioutil.ReadFile(tempFileName)
	checkErr(err)

	checkErr(yamlToSecretData(secret, readBytes))

	checkErr(os.Remove(tempFileName))

	clientset.CoreV1().Secrets(*namespace).Update(secret)
}

func ListSecrets(clientset *kubernetes.Clientset, namespace *string) {
	secrets, err := clientset.CoreV1().Secrets(*namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, secret := range secrets.Items {
		fmt.Println(secret.Name)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func secretDataToYaml(secret *v1.Secret) (out []byte, err error) {
	data := make(map[string]string)
	for key, value := range secret.Data {
		data[key] = string(value)
	}
	return yaml.Marshal(data)
}

func yamlToSecretData(secret *v1.Secret, bytes []byte) error {
	data := make(map[string]string)

	err := yaml.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	for key, value := range data {
		secret.Data[key] = []byte(value)
	}

	return nil
}