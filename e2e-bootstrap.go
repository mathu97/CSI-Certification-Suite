package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"os"
	"path/filepath"
	"time"
)

var deleteReclaimPolicy = v1.PersistentVolumeReclaimDelete

type driverInfo struct {
	name         string
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

func getPluginInfo(network string, endpoint string) *driverInfo {
	var driver = driverInfo{}
	ctx := context.Background()

	//Connect to the driver through the endpoint
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	clientCon, conErr := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
		grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
			return net.Dial(network, target)
		}),
	)

	if conErr != nil {
		fmt.Printf("Unable to connect to Driver via gRPC\n")
		conErr.Error()
	}

	identityClient := csi.NewIdentityClient(clientCon)

	pInfoReq := &csi.GetPluginInfoRequest{}
	res, pluginInfoErr := identityClient.GetPluginInfo(ctx, pInfoReq)

	if pluginInfoErr != nil {
		fmt.Printf("Unable to get plugin info from identity server\n")
		fmt.Printf("Error: %s\n", pluginInfoErr)
		pluginInfoErr.Error()
	}

	driver.name = res.GetName()
	fmt.Printf("Driver Name: %s\n", driver.name)

	plugCapReq := &csi.GetPluginCapabilitiesRequest{}
	capRes, plugCapErr := identityClient.GetPluginCapabilities(ctx, plugCapReq)

	if plugCapErr != nil {
		fmt.Printf("Unable to get plugin capabilities from identity server\n")
		plugCapErr.Error()
	}

	driver.capabilities = capRes.Capabilities

	clientCon.Close()

	return &driver

}

func main() {
	//Sample command line call: go run e2e-bootstrap.go --endpoint 127.0.0.1:10000 --network tcp
	var endPoint string
	var network string
	var kubeConfig *string

	flag.StringVar(&endPoint, "endpoint", "", "Provide the driver's endpoint")
	flag.StringVar(&network, "network", "", "Provide the netowrk for the driver endpoint (ex: tcp or unix)")

	if home := homeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	if endPoint == "" || network == "" {
		fmt.Printf("Need to provide Driver Network and Endpoint\n")
		return
	}

	config, confErr := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if confErr != nil {
		panic(confErr.Error())
	}

	//Get the driver's name
	driverInfo := getPluginInfo(network, endPoint)

	clientSet, csErr := kubernetes.NewForConfig(config)
	if csErr != nil {
		panic(csErr.Error())
	}

	newStorageClass := createStorageClass(driverInfo.name)

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
