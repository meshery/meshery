package meshmodel_policy

# playground with this policy: https://play.openpolicyagent.org/p/ttJQwNEGQ8

# gets the container pod protocol, default is "TCP"
get_port_protocol(container_ports) = pod_protocol {
	not has_key(container_ports, "protocol")
	pod_protocol := "TCP"
}

get_port_protocol(container_ports) = pod_protocol {
	has_key(container_ports, "protocol")
	pod_protocol := container_ports.protocol
}

# checks for key in object
has_key(x, k) { _ = x[k] }

# checkServiceToContainerPortMap(service_port, container_port) {
# 	# the service target-port should be equal to the container port name
# 	service_port["target Port"] == container_port.name
# }

checkServiceToContainerPortMap(service_port, container_port) {
	# the service target-port can be number as well and should be equal to container port in that case
	service_port["targetPort"] == container_port["containerPort"]
}

checkServiceToContainerPortMap(service_port, container_port) {
	# or the service port number equals to the the container Port Number
	service_port.port == container_port["containerPort"]
}

edge_network_relationship[service_name] {
	# Select a service resource
	service := input.services[_]

	# Check if the resource is a Service
	service.type == "Service"

	# Extract the Service spec
	service_spec := service.settings.spec

	# Extract the Service ports
	service_port := service_spec.ports[_]

	# Extract the Service protocol
	service_protocol := get_port_protocol(service_port)

	# Select a Pod resource with the same name as the Service
	pod := input.services[_]

	# Check if the resource is a Pod
	pod.type == "Pod"

	# Extract the Pod spec
	pod_spec := pod.settings.spec

	# Extract the Pod ports
	pod_ports := pod_spec.containers[_].ports[_]

	#pod protocol
	pod_protocol := get_port_protocol(pod_ports) 

	# Check if the Service protocol matches the Pod protocol
	service_protocol == pod_protocol

	# Check if the Service target port matches the Pod port
	checkServiceToContainerPortMap(service_port, pod_ports)

	source_id := service.traits.meshmap.id
	destination_id := pod.traits.meshmap.id
	service_name = {"destination_name": pod.name, "source_id": source_id, "destination_id": destination_id, "source_name": service.name, "port": service_port}
}

edge_network_relationship[service_name] {
    # Select a service resource
    service := input.services[_]
    
    # Check if the resource is a Service
    service.type == "Service"

    # Extract the Service spec
    service_spec := service.settings.spec

    # Extract the Service ports
    service_port := service_spec.ports[_]
    
    # Extract the Service protocol
    service_protocol := get_port_protocol(service_port)

    # Select a Pod resource with the same name as the Service
    deployment := input.services[_]

    # Check if the resource is a deployment
    deployment.type == "Deployment"

    # Extract the deployment containers
    deployment_containers := deployment.settings.spec.template.spec.containers[_]

    # Extract the deployment ports
    deployment_ports := deployment_containers.ports[_]

    # Check if the Service port matches the deployment port
    checkServiceToContainerPortMap(service_port, deployment_ports)

    # Define the deployment protocol as "TCP" since it is not provided in the data
    deployment_protocol := get_port_protocol(deployment_ports) 

    # Check if the Service protocol matches the deployment protocol
    service_protocol == deployment_protocol

    source_id := service.traits.meshmap.id
    destination_id := deployment.traits.meshmap.id
    
    service_name = {"destination_name": deployment.name, "source_id": source_id, "destination_id": destination_id, "source_name": service.name, "port": service_port}
}
