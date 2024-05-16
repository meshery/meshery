#!/bin/bash
#TODO: Make the cleanup smarter than basic namespace deletion
#Script is placed in /root directory in VM.

valid_namespaces=("kube-system" "default" "monitoring" "kube-flannel" "kube-node-lease" "kube-public" "meshery" "metallb-system" "projectcontour" "ingress-nginx" "layer5-cloud" "postgres")

for ns in $(kubectl get ns -o jsonpath="{.items[*].metadata.name}");
do
    if [[ ! " ${valid_namespaces[*]} " =~ [[:space:]]${ns}[[:space:]] ]]; then
        echo "Deleting namespace $ns"
        timeout 2m kubectl delete ns "$ns"
        ns_status=$(kubectl get namespace $ns -o jsonpath='{.status.phase}')
        # Check if the namespace is in the "Terminating" state
        if [ "$ns_status" == "Terminating" ]; then
            echo "Namespace $ns is stuck in 'Terminating' state. Patching finalizers field..."
            # Patch finalizers field
            kubectl patch namespace $ns -p '{"metadata":{"finalizers":null}}'
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