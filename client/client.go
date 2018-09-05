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
	"path"
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
func ListSecrets(clientset *kubernetes.Clientset, namespace *string) {
	secrets, err := clientset.CoreV1().Secrets(*namespace).List(metav1.ListOptions{})
	checkErr(err)
	for _, secret := range secrets.Items {
		fmt.Println(secret.Name)
	}
}
func ListSecretKeys(clientset *kubernetes.Clientset, namespace *string, secretName string) {
	secret, err := clientset.CoreV1().Secrets(*namespace).Get(secretName, metav1.GetOptions{})
	checkErr(err)
	for key := range secret.Data {
		fmt.Println(key)
	}
}
func EditSecret(clientset *kubernetes.Clientset, namespace *string, secretName string) {
	secret, err := clientset.CoreV1().Secrets(*namespace).Get(secretName, metav1.GetOptions{})
	checkErr(err)

	inBytes, err := secretDataToYaml(secret)
	checkErr(err)

	outBytes, err := editInEditor(inBytes, secretName)
	checkErr(err)

	checkErr(yamlToSecretData(secret, outBytes))

	clientset.CoreV1().Secrets(*namespace).Update(secret)
}

func EditSecretKey(clientset *kubernetes.Clientset, namespace *string, secretName string, secretKey string) {
	secret, err := clientset.CoreV1().Secrets(*namespace).Get(secretName, metav1.GetOptions{})
	checkErr(err)

	inBytes := secret.Data[secretKey]

	outBytes, err := editInEditor(inBytes, secretKey)
	checkErr(err)

	secret.Data[secretKey] = outBytes

	clientset.CoreV1().Secrets(*namespace).Update(secret)
}

func editInEditor(in []byte, name string) (out []byte, err error)  {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	tempFileName := path.Join(tempDir, name)
	if err := ioutil.WriteFile(tempFileName, in, 0664); err != nil {
		return nil, err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	editorCmd := exec.Command(editor, tempFileName)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return nil, err
	}

	readBytes, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		return nil, err
	}

	if os.Remove(tempFileName) != nil {
		return nil, err
	}

	return readBytes, nil
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

	keysDeleted := make(map[string]bool)
	for secretKey, _ := range secret.Data {
		keysDeleted[secretKey] = true
	}

	if secret.Data == nil {
		secret.Data =  map[string][]byte{}
        }
	
	for key, value := range data {
		secret.Data[key] = []byte(value)
		keysDeleted[key] = false
	}

	for key, deleted := range keysDeleted {
		if deleted {
			delete(secret.Data, key)
		}
	}

	return nil
}
