package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

func logGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	glog.V(5).Infof("GRPC call: %s", method)
	glog.V(5).Infof("GRPC request: %+v", req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	glog.V(5).Infof("GRPC response: %+v", reply)
	glog.V(5).Infof("GRPC error: %v", err)
	return err
}

func connect(network string, endpoint string) (*grpc.ClientConn, error) {
	//Connect to the driver through the endpoint
	clientCon, conErr := grpc.Dial(
		endpoint,
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(time.Second),
		grpc.WithUnaryInterceptor(logGRPC),
		grpc.WithDialer(func(target string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout(network, target, timeout)
		}),
	)

	if conErr != nil {
		return nil, conErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for {
		if !clientCon.WaitForStateChange(ctx, clientCon.GetState()) {
			glog.V(4).Infof("Connection timed out")
			return clientCon, nil // return nil, subsequent GetPluginInfo will show the real connection error
		}
		if clientCon.GetState() == connectivity.Ready {
			glog.V(3).Infof("Connected")
			return clientCon, nil
		}
		glog.V(4).Infof("Still trying, connection is %s", clientCon.GetState())
	}

}

func getPluginInfo(connection *grpc.ClientConn) *driverInfo {
	var driver = driverInfo{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	identityClient := csi.NewIdentityClient(connection)
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

	connection.Close()

	return &driver

}

func main() {
	//Sample command line call: go run e2e-bootstrap.go --endpoint 127.0.0.1:10000 --network tcp
	var endPoint string
	var network string
	var kubeConfig *string

	flag.StringVar(&endPoint, "endpoint", "", "Provide the driver's endpoint")
	flag.StringVar(&network, "network", "", "Provide the network for the driver endpoint (ex: tcp or unix)")

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

	//Get the driver's name
	clientConnection, connectionErr := connect(network, endPoint)
	if connectionErr != nil {
		fmt.Printf("Unable to connect to Driver via gRPC\n")
		connectionErr.Error()
	}

	driverInfo := getPluginInfo(clientConnection)
	newStorageClass := createStorageClass(driverInfo.name)

	config, confErr := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if confErr != nil {
		panic(confErr.Error())
	}
	clientSet, csErr := kubernetes.NewForConfig(config)
	if csErr != nil {
		panic(csErr.Error())
	}

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
