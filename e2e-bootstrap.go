package CSI_Certification_Suite

import (
	storage "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
			"github.com/kubernetes-csi/external-provisioner/pkg/controller"
)

var pluginName = "" //Name of the csi driver running
var deleteReclaimPolicy = v1.PersistentVolumeReclaimDelete

var targetStorageClass = &storage.StorageClass{
	TypeMeta: metav1.TypeMeta{
		Kind: "StorageClass",
	},

	ObjectMeta: metav1.ObjectMeta{
		Name:"gold",
	},

	Provisioner: pluginName,
	Parameters: nil,
	ReclaimPolicy: &deleteReclaimPolicy,
}

func getDriverName(address string) (string, error) {
	conn, err := controller.Connect(address, 10)

	if err == nil{
		name, err := controller.GetDriverName(conn, 10)
		return name, err
	}

	return "", err
}