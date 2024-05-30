#!/bin/bash
#TODO: Make the cleanup smarter than basic namespace deletion
#Script is placed in /root directory in VM.

valid_namespaces=("kube-system" "default" "monitoring" "kube-flannel" "kube-node-lease" "kube-public" "meshery" "metallb-system" "projectcontour" "ingress-nginx" "layer5-cloud" "postgres")

for ns in $(kubectl get ns -o jsonpath="{.items[*].metadata.name}"); 
do
    if [[ ! " ${valid_namespaces[*]} " =~ [[:space:]]${ns}[[:space:]] ]]; then
        echo "Deleting all resources in namespace $ns"
        kubectl delete all --all -n $ns
        echo "Deleting namespace $ns"
        timeout 1m kubectl delete ns $ns

        ns_status=$(kubectl get namespace $ns -o jsonpath='{.status.phase}')
        # Check if the namespace is in the "Terminating" state
        if [ "$ns_status" == "Terminating" ]; then
            echo "Namespace $ns is stuck in 'Terminating' state. Patching finalizers field..."
            # Patch finalizers field
            kubectl get namespace $ns -o json | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/" | kubectl replace --raw /api/v1/namespaces/$ns/finalize -f -
        fi
    fi
done

for ns in $valid_namespaces;
do
    kubectl delete pods -n $ns \
    --field-selector="status.phase!=Running,status.phase!=ContainerCreating"
done

# remove mutatingwebhookconfigurations
for mwh in $(kubectl get mutatingwebhookconfigurations -o jsonpath="{.items[*].metadata.name}");
do
    if [[ "$mwh" == "consul-consul-connect-injector" ]];then
        echo "Deleting mutatingwebhookconfigurations $mwh"
        kubectl delete mutatingwebhookconfigurations "$mwh"
    fi
done