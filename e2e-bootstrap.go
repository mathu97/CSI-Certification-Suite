package main

import (
	"flag"
	"fmt"
	"github.com/kubernetes-csi/external-provisioner/pkg/controller"
	"k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var deleteReclaimPolicy = v1.PersistentVolumeReclaimDelete

func getDriverName(address string) (string, error) {
	conn, err := controller.Connect(address, 10)

	if err == nil {
		name, err := controller.GetDriverName(conn, 10)
		return name, err
	}

	return "", err
}

func createStorageClass(driverName string) *storage.StorageClass {

	var targetStorageClass = &storage.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind: "StorageClass",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name: "e2e-csi-test",
		},

		Provisioner:   driverName,
		Parameters:    nil,
		ReclaimPolicy: &deleteReclaimPolicy,
	}

	return targetStorageClass
}

func main() {
	endPointPtr := flag.String("endpoint", "foo", "a string")
	flag.Parse()

	if *endPointPtr == "foo" {
		fmt.Errorf("Need to provide Driver Endpoint\n")
		return
	}

	//Get the driver's name
	driverName, err := getDriverName(*endPointPtr)
	if err != nil {
		fmt.Errorf("Unable to get Driver Name\n")
	}

	newStorageClass := createStorageClass(driverName)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientSet, cErr := kubernetes.NewForConfig(config)
	if cErr != nil {
		panic(cErr.Error())
	}

	//Create the e2e-csi-test storage class object in the cluster
	if _, scErr := clientSet.StorageV1().StorageClasses().Create(newStorageClass); scErr != nil {
		panic(scErr.Error())
	}

}
