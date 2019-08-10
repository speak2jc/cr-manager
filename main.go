package main

import (
	"crypto/rand"
	"fmt"
	"github.com/speak2jc/k-op/pkg/apis/example/v1alpha1"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	keevakindGVR = schema.GroupVersionResource{
		Group:    "example.keeva.com",
		Version:  "v1alpha1",
		Resource: "keevakinds",
	}
)

func main() {
	log.Print("Loading client config")
	config, err := clientcmd.BuildConfigFromFlags("", userConfig())
	errExit("Failed to load client conifg", err)

	log.Print("Loading dynamic client")
	client, err := dynamic.NewForConfig(config)
	errExit("Failed to create client", err)

	name := "keeva-" + RandomString(5)
	namespace := "james"

	RegisterRuntimeClassCRD(config)
	CreateSampleKeevaKinds(client, name, namespace)
	kk := GetSampleKeevaKinds(client, name, namespace)
	log.Printf("Found existing Keevakind %s", kk)
	//UpdateSampleKeevaKinds(client, name, namespace)
}

func RegisterRuntimeClassCRD(config *rest.Config) {
	apixClient, err := apixv1beta1client.NewForConfig(config)
	errExit("Failed to load apiextensions client", err)

	crds := apixClient.CustomResourceDefinitions()

	runtimeClassCRD := &apixv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "keevakinds.example.keeva.com",
		},
		Spec: apixv1beta1.CustomResourceDefinitionSpec{
			Group:   "example.keeva.com",
			Version: "v1alpha1",
			Versions: []apixv1beta1.CustomResourceDefinitionVersion{{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
			}},
			Names: apixv1beta1.CustomResourceDefinitionNames{
				Plural:   "keevakinds",
				Singular: "keevakind",
				Kind:     "Keevakind",
			},
			Scope: apixv1beta1.ClusterScoped,
			Validation: &apixv1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &apixv1beta1.JSONSchemaProps{
					Properties: map[string]apixv1beta1.JSONSchemaProps{
						"spec": {
							Properties: map[string]apixv1beta1.JSONSchemaProps{
								"runtimeHandler": {
									Type: "string",
									//Pattern: "abc",
								},
								"kind": {
									Type: "string",
								},
							},
						},
					},
				},
			},
		},
	}
	log.Print("Registering Keevakind CRD")
	_, err = crds.Create(runtimeClassCRD)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Print("Keevakind CRD already registered")
		} else {
			errExit("Failed to create Keevakind CRD", err)
		}
	}
}

func CreateSampleKeevaKinds(client dynamic.Interface, name string, namespace string) {

	res := client.Resource(keevakindGVR)

	var count int32 = 1
	var port int32 = 8080
	group := "Group-" + RandomString(5)
	image := "Image-" + RandomString(5)

	log.Printf("Creating Keevakind %s", name)
	rc := NewKeevaKind(name, namespace, count, group, image, port)
	log.Printf("rc %s", rc)
	_, err := res.Create(rc, metav1.CreateOptions{})
	errExit(fmt.Sprintf("Failed to create Keevakind %#v", rc), err)

}

func UpdateSampleKeevaKinds(client dynamic.Interface, name string, namespace string) {

	res := client.Resource(keevakindGVR)

	//var count int32 = 2
	//var port int32 = 8080
	//group := "Group-" + RandomString(5)
	//image := "Image-" + RandomString(5)

	log.Printf("Getting Keevakind %s", name)
	//rc := NewKeevaKind(name, namespace, count, group, image, port)
	//log.Printf("rc %s", rc)

	rc, err := res.Get(name, metav1.GetOptions{})

	errExit(fmt.Sprintf("Failed to create Keevakind %#v", rc), err)

}

func GetSampleKeevaKinds(client dynamic.Interface, name string, namespace string) *unstructured.Unstructured {

	res := client.Resource(keevakindGVR)

	log.Printf("Getting Keevakind %s", name)

	var existingKeeva *unstructured.Unstructured
	existingKeeva, err := res.Get(name, metav1.GetOptions{})

	errExit(fmt.Sprintf("Failed to create Keevakind %s in namespace %s", name, namespace), err)

	return existingKeeva
}

func errExit(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %#v", msg, err)
	}
}

func userConfig() string {
	usr, err := user.Current()
	errExit("Failed to get current user", err)

	return filepath.Join(usr.HomeDir, ".kube", "config")
}

func NewKeevaKind(name string, namespace string, count int32, group string, image string, port int32) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Keevakind",
			"apiVersion": keevakindGVR.Group + "/v1alpha1",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": v1alpha1.KeevakindSpec{
				Count: count,
				Group: group,
				Image: image,
				Port:  port,
			},
		},
	}
}

func RandomString(len int) string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	fmt.Println(s)

	return strings.ToLower(s)
}
