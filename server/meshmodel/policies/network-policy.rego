package service_pod_port_matching

service_pod_relationships[service_name] {
    print("service_target_port")

    # Select a service resource
    service := input.services[_]
#     print("service", service_name)
    
    # Check if the resource is a Service
    service.type == "Service"
    print("service")
    

    # Extract the Service spec
    service_spec := service.settings.spec
    print("service spec", service_spec)

    # Extract the Service ports
    service_port := service_spec.ports[_]
    print("service port", service_port)
    

    # Extract the Service protocol
    service_protocol := service_port.protocol
    print("service_protocol", service_protocol)

    # Extract the Service target port
    service_target_port := service_port["target Port"]

    # Select a Pod resource with the same name as the Service
    pod := input.services[_]
    print("pod", pod)

    # Check if the resource is a Pod
    pod.type == "Pod"

    # Extract the Pod spec
    pod_spec := pod.settings.spec
#     print("pod_Spec", pod_spec)

    # Extract the Pod containers
    pod_containers := pod_spec.containers[_]

    # Extract the Pod ports
    pod_ports := pod_containers.ports[_]

    # Extract the Pod port
    pod_port := pod_ports["container Port"]
#     print(pod_ports, service_port)

    # Check if the Service port matches the Pod port
    service_port.port == pod_ports["container Port"]


    # Define the Pod protocol as "TCP" since it is not provided in the data
    pod_protocol := "TCP"
    print("came here", service_protocol, pod_protocol, service_target_port, pod_port, pod_ports.name)

    # Check if the Service protocol matches the Pod protocol
    service_protocol == pod_protocol

    # Check if the Service target port matches the Pod port
#     service_target_port == pod_port | 
    service_target_port == pod_ports.name
    
    print("service_target_port", service_target_port)
    print("service_port", service_port)
    source_id := service.traits.meshmap.id
    destination_id := pod.traits.meshmap.id
    
    service_name = {"destination_name": pod.name, "source_id": source_id, "destination_id": destination_id, "source_name": service.name}
}


service_deployment_relationship[service_name] {
    print("service_target_port")

    # Select a service resource
    service := input.services[_]
#     print("service", service_name)
    
    # Check if the resource is a Service
    service.type == "Service"
    print("service")
    

    # Extract the Service spec
    service_spec := service.settings.spec
    print("service spec", service_spec)

    # Extract the Service ports
    service_port := service_spec.ports[_]
    print("service port", service_port)
    

    # Extract the Service protocol
    service_protocol := service_port.protocol
    print("service_protocol", service_protocol)

    # Extract the Service target port
    service_target_port := service_port["target Port"]

    # Select a Pod resource with the same name as the Service
    deployment := input.services[_]

    # Check if the resource is a deployment
    deployment.type == "Deployment"
    print("deployment", deployment)

    # Extract the deployment containers
    deployment_containers := deployment.settings.spec.template.spec.containers[_]
    print("deployment_con", deployment_containers)

    # Extract the deployment ports
    deployment_ports := deployment_containers.ports[_]

    # Extract the deployment port
    deployment_port := deployment_ports["container Port"]
    print(deployment_ports, service_port)

    # Check if the Service port matches the deployment port
    service_port["target Port"] == deployment_ports["container Port"] # TODO: has to be changed to OR statements 


    # Define the deployment protocol as "TCP" since it is not provided in the data
    deployment_protocol := service_port.protocol
    print("came here", service_protocol, deployment_protocol, service_target_port, deployment_port, deployment_ports.name)

    # Check if the Service protocol matches the deployment protocol
    service_protocol == deployment_protocol

    # Check if the Service target port matches the deployment port
	#     service_target_port == deployment_port | 
	#     service_target_port == deployment_ports.name
	# TODO: has to be uncommented, above lines
    
    print("service_target_port", service_target_port)
    print("service_port", service_port)
    source_id := service.traits.meshmap.id
    destination_id := deployment.traits.meshmap.id
    
    service_name = {"destination_name": deployment.name, "source_id": source_id, "destination_id": destination_id, "source_name": service.name}
}
