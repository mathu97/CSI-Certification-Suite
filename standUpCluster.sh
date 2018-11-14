#!/bin/bash

if ["$KUBE_PATH" = ""]; then
    printf "\nMust set KUBE_PATH environment varibale to your kubernetes path\n\n"
    exit 1
fi

(cd $KUBE_PATH && make quick-release)
printf "\nMake Complete\n\n"

#Generate gce ssh keys to stand up a cluster
ssh-keygen -t rsa -f ~/.ssh/google_compute_engine -C root
printf "\nSSH key generated for GCE\n\n"

#Stand Up a Cluster
(cd $KUBE_PATH && go run hack/e2e.go -- --provider=gce --gcp-nodes=1 --gcp-network=mselvara-e2e --gcp-project=openshift-gce-devel --gcp-zone=us-east1-d --up)

#Add bootstrapping and e2e test part here   

#Bring Cluster Down
#(cd $KUBE_PATH && go run hack/e2e.go -- --provider=gce --gcp-nodes=1 --gcp-network=mselvara-e2e --gcp-project=openshift-gce-devel --gcp-zone=us-east1-d --down)
