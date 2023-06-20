#!/bin/bash
#TODO: Make the cleanup smarter than basic namespace deletion
#Script is placed in /root directory in VM.
for ns in $(kubectl get ns -o jsonpath="{.items[*].metadata.name}");
do
    if [[ "$ns" != "kube-system" && "$ns" != "default" && "$ns" != "monitoring" && "$ns" != "kube-flannel" && "$ns" != "kube-node-lease" && "$ns" != "kube-public" && "$ns" != "meshery" && "$ns" != "metallb-system" && "$ns" != "projectcontour" ]];then
          echo "Deleting namespace $ns"
          kubectl delete ns "$ns"
    fi
done