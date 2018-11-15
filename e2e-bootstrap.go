package main

import (
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
		"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"os"
	"google.golang.org/grpc"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"context"
)

var deleteReclaimPolicy = v1.PersistentVolumeReclaimDelete

type driverInfo struct {
	name string
	capabilities []*csi.PluginCapability
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

func getPluginInfo(endpoint string) (*driverInfo){
	var driver = driverInfo{}
	ctx := context.Background()

	//Connect to the driver through the endpoint
	clientCon, conErr := grpc.DialContext(ctx, endpoint)

	if conErr != nil {
		fmt.Printf("Unable to connect to Driver via gRPC")
		conErr.Error()
	}

	identityClient := csi.NewIdentityClient(clientCon)
	pInfoReq := &csi.GetPluginInfoRequest{}
	res, pluginInfoErr := identityClient.GetPluginInfo(ctx, pInfoReq)

	if pluginInfoErr != nil {
		fmt.Printf("Unable to get plugin info from identity server")
		pluginInfoErr.Error()
	}

	driver.name = res.GetName()

	plugCapReq := &csi.GetPluginCapabilitiesRequest{}
	capRes, plugCapErr := identityClient.GetPluginCapabilities(ctx, plugCapReq)

	if plugCapErr != nil {
		fmt.Printf("Unable to get plugin capabilities from identity server")
		plugCapErr.Error()
	}

	driver.capabilities = capRes.Capabilities

	clientCon.Close()

	return &driver

}

func main() {
	endPointPtr := flag.String("endpoint", "foo", "a string")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	if *endPointPtr == "foo" {
		fmt.Printf("Need to provide Driver Endpoint\n")
		return
	}

	config, confErr := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if confErr != nil {
		panic(confErr.Error())
	}

	//Get the driver's name
	driverName, err := getDriverName(*endPointPtr)
	if err != nil {
		fmt.Printf("Unable to get Driver Name\n")
	}

	clientSet, csErr := kubernetes.NewForConfig(config)
	if csErr != nil {
		panic(csErr.Error())
	}

	newStorageClass := createStorageClass(driverName)

	//Create the e2e-csi-test storage class object in the cluster
	if _, scErr := clientSet.StorageV1().StorageClasses().Create(newStorageClass); scErr != nil {
		panic(scErr.Error())
	}



}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

