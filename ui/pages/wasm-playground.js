import dynamic from 'next/dynamic'
import { loadPolicy } from "@open-policy-agent/opa-wasm"
import { useEffect } from 'react';

const WasmComponent = dynamic({
  loader : async () => {
    const go = new globalThis.Go();
    WebAssembly.instantiateStreaming(fetch("static/wasm/main.wasm"), go.importObject).then((result) => {
      go.run(result.instance)
    });
  },
})

async function loadWasmPolicy() {
  try {
    const policyWasmResp = await fetch("static/wasm/schemaInfoExtractor.wasm");
    const policyWasm = await policyWasmResp.arrayBuffer();
    const policy = await loadPolicy(policyWasm);

    policy.setData({
      resource_types : {
        "ports" : "ports",
        "cluster IP" : "Cluster IP"
      }
    })

    const res = policy.evaluate(serviceSchema);
    console.log(res[0].result)
  } catch (e) {
    console.log("[..er", e)
  }
}

function WasmSchemaValidator() {
  useEffect(() => {
    loadWasmPolicy()
  })

  return <>
  Wasm schema Validator
  </>
}

export default function Index() {
  return (
    <div>
      <WasmComponent />
      <WasmSchemaValidator />
    </div>
  )
}

const serviceSchema = {
  "description" : "Service is a named abstraction of software service (for example, mysql) consisting of local port (for example 3306) that the proxy listens on, and the selector that determines which pods will answer requests sent through the proxy.",
  "properties" : {
    "spec" : {
      "allOf" : [
        {
          "description" : "Service Spec describes the attributes that a user creates on a service.",
          "properties" : {
            "allocate Load Balancer Node Ports" : {
              "description" : "allocate Load Balancer Node Ports defines if Node Ports will be automatically allocated for services with type Load Balancer. Default is \"true\". It may be set to \"false\" if the cluster load-balancer does not rely on Node Ports. If the caller requests specific Node Ports (by specifying a value), those requests will be respected, regardless of this field. This field may only be set for services with type Load Balancer and will be cleared if the type is changed to any other type.",
              "type" : "boolean"
            },
            "cluster IP" : {
              "description" : "cluster IP is the IP address of the service and is usually assigned randomly. If an address is specified manually, is in-range (as per system configuration), and is not in use, it will be allocated to the service; otherwise creation of the service will fail. This field may not be changed through updates unless the type field is also being changed to External Name (which requires this field to be blank) or the type field is being changed from External Name (in which case this field may optionally be specified, as describe above). Valid values are \"None\", empty string (\"\"), or a valid IP address. Setting this to \"None\" makes a \"headless service\" (no virtual IP), which is useful when direct endpoint connections are preferred and proxying is not required. Only applies to types Cluster IP, Node Port, and Load Balancer. If this field is specified when creating a Service of type External Name, creation will fail. This field will be wiped when updating a Service to type External Name. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies",
              "type" : "string"
            },
            "cluster IPs" : {
              "description" : "Cluster I Ps is a list of IP addresses assigned to this service, and are usually assigned randomly. If an address is specified manually, is in-range (as per system configuration), and is not in use, it will be allocated to the service; otherwise creation of the service will fail. This field may not be changed through updates unless the type field is also being changed to External Name (which requires this field to be empty) or the type field is being changed from External Name (in which case this field may optionally be specified, as describe above). Valid values are \"None\", empty string (\"\"), or a valid IP address. Setting this to \"None\" makes a \"headless service\" (no virtual IP), which is useful when direct endpoint connections are preferred and proxying is not required. Only applies to types Cluster IP, Node Port, and Load Balancer. If this field is specified when creating a Service of type External Name, creation will fail. This field will be wiped when updating a Service to type External Name. If this field is not specified, it will be initialized from the cluster IP field. If this field is specified, clients must ensure that cluster I Ps[0] and cluster IP have the same value. This field may hold a maximum of two entries (dual-stack I Ps, in either order). These I Ps must correspond to the values of the ip Families field. Both cluster I Ps and ip Families are governed by the ip Family Policy field. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies",
              "items" : {
                "default" : "",
                "type" : "string"
              },
              "type" : "array",
              "x-kubernetes-list-type" : "atomic"
            },
            "external IPs" : {
              "description" : "external I Ps is a list of IP addresses for which nodes in the cluster will also accept traffic for this service. These I Ps are not managed by Kubernetes. The user is responsible for ensuring that traffic arrives at a node with this IP. A common example is external load-balancers that are not part of the Kubernetes system.",
              "items" : {
                "default" : "",
                "type" : "string"
              },
              "type" : "array"
            },
            "external Name" : {
              "description" : "external Name is the external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record). No proxying will be involved. Must be a lowercase RFC-1123 hostname (https://tools.ietf.org/html/rfc1123) and requires `type` to be \"External Name\".",
              "type" : "string"
            },
            "external Traffic Policy" : {
              "description" : "external Traffic Policy describes how nodes distribute service traffic they receive on one of the Service's \"externally-facing\" addresses (Node Ports, External I Ps, and Load Balancer I Ps). If set to \"Local\", the proxy will configure the service in a way that assumes that external load balancers will take care of balancing the service traffic between nodes, and so each node will deliver traffic only to the node-local endpoints of the service, without masquerading the client source IP. (Traffic mistakenly sent to a node with no endpoints will be dropped.) The default value, \"Cluster\", uses the standard behavior of routing to all endpoints evenly (possibly modified by topology and other features). Note that traffic sent to an External IP or Load Balancer IP from within the cluster will always get \"Cluster\" semantics, but clients sending to a Node Port from within the cluster may need to take traffic policy into account when picking a node. Possible enum values: - `\"Cluster\"` routes traffic to all endpoints. - `\"Local\"` preserves the source IP of the traffic by routing only to endpoints on the same node as the traffic was received on (dropping the traffic if there are no local endpoints).",
              "enum" : [
                "Cluster",
                "Local"
              ],
              "type" : "string"
            },
            "health Check Node Port" : {
              "description" : "health Check Node Port specifies the healthcheck node Port for the service. This only applies when type is set to Load Balancer and external Traffic Policy is set to Local. If a value is specified, is in-range, and is not in use, it will be used. If not specified, a value will be automatically allocated. External systems (e.g. load-balancers) can use this port to determine if a given node holds endpoints for this service or not. If this field is specified when creating a Service which does not need it, creation will fail. This field will be wiped when updating a Service to no longer need it (e.g. changing type). This field cannot be updated once set.",
              "format" : "int32",
              "type" : "integer"
            },
            "internal Traffic Policy" : {
              "description" : "Internal Traffic Policy describes how nodes distribute service traffic they receive on the Cluster IP. If set to \"Local\", the proxy will assume that pods only want to talk to endpoints of the service on the same node as the pod, dropping the traffic if there are no local endpoints. The default value, \"Cluster\", uses the standard behavior of routing to all endpoints evenly (possibly modified by topology and other features).",
              "type" : "string"
            },
            "ip Families" : {
              "description" : "IP Families is a list of IP families (e.g. I Pv4, I Pv6) assigned to this service. This field is usually assigned automatically based on cluster configuration and the ip Family Policy field. If this field is specified manually, the requested family is available in the cluster, and ip Family Policy allows it, it will be used; otherwise creation of the service will fail. This field is conditionally mutable: it allows for adding or removing a secondary IP family, but it does not allow changing the primary IP family of the Service. Valid values are \"I Pv4\" and \"I Pv6\". This field only applies to Services of types Cluster IP, Node Port, and Load Balancer, and does apply to \"headless\" services. This field will be wiped when updating a Service to type External Name. This field may hold a maximum of two entries (dual-stack families, in either order). These families must correspond to the values of the cluster I Ps field, if specified. Both cluster I Ps and ip Families are governed by the ip Family Policy field.",
              "items" : {
                "default" : "",
                "type" : "string"
              },
              "type" : "array",
              "x-kubernetes-list-type" : "atomic"
            },
            "ip Family Policy" : {
              "description" : "IP Family Policy represents the dual-stack-ness requested or required by this Service. If there is no value provided, then this field will be set to Single Stack. Services can be \"Single Stack\" (a single IP family), \"Prefer Dual Stack\" (two IP families on dual-stack configured clusters or a single IP family on single-stack clusters), or \"Require Dual Stack\" (two IP families on dual-stack configured clusters, otherwise fail). The ip Families and cluster I Ps fields depend on the value of this field. This field will be wiped when updating a service to type External Name.",
              "type" : "string"
            },
            "load Balancer Class" : {
              "description" : "load Balancer Class is the class of the load balancer implementation this Service belongs to. If specified, the value of this field must be a label-style identifier, with an optional prefix, e.g. \"internal-vip\" or \"example.com/internal-vip\". Unprefixed names are reserved for end-users. This field can only be set when the Service type is 'Load Balancer'. If not set, the default load balancer implementation is used, today this is typically done through the cloud provider integration, but should apply for any default implementation. If set, it is assumed that a load balancer implementation is watching for Services with a matching class. Any default load balancer implementation (e.g. cloud providers) should ignore Services that set this field. This field can only be set when creating or updating a Service to type 'Load Balancer'. Once set, it can not be changed. This field will be wiped when a service is updated to a non 'Load Balancer' type.",
              "type" : "string"
            },
            "load Balancer IP" : {
              "description" : "Only applies to Service Type: Load Balancer. This feature depends on whether the underlying cloud-provider supports specifying the load Balancer IP when a load balancer is created. This field will be ignored if the cloud-provider does not support the feature. Deprecated: This field was under-specified and its meaning varies across implementations, and it cannot support dual-stack. As of Kubernetes v1.24, users are encouraged to use implementation-specific annotations when available. This field may be removed in a future API version.",
              "type" : "string"
            },
            "load Balancer Source Ranges" : {
              "description" : "If specified and supported by the platform, this will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client I Ps. This field will be ignored if the cloud-provider does not support the feature.\" More info: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/",
              "items" : {
                "default" : "",
                "type" : "string"
              },
              "type" : "array"
            },
            "ports" : {
              "description" : "The list of ports that are exposed by this service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies",
              "items" : {
                "allOf" : [
                  {
                    "description" : "Service Port contains information on service's port.",
                    "properties" : {
                      "app Protocol" : {
                        "description" : "The application protocol for this port. This field follows standard Kubernetes label syntax. Un-prefixed names are reserved for IANA standard service names (as per RFC-6335 and https://www.iana.org/assignments/service-names). Non-standard protocols should use prefixed names such as mycompany.com/my-custom-protocol.",
                        "type" : "string"
                      },
                      "name" : {
                        "description" : "The name of this port within the service. This must be a DNS_LABEL. All ports within a Service Spec must have unique names. When considering the endpoints for a Service, this must match the 'name' field in the Endpoint Port. Optional if only one Service Port is defined on this service.",
                        "type" : "string"
                      },
                      "node Port" : {
                        "description" : "The port on each node on which this service is exposed when type is Node Port or Load Balancer. Usually assigned by the system. If a value is specified, in-range, and not in use it will be used, otherwise the operation will fail. If not specified, a port will be allocated if this Service requires one. If this field is specified when creating a Service which does not need it, creation will fail. This field will be wiped when updating a Service to no longer need it (e.g. changing type from Node Port to Cluster IP). More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport",
                        "format" : "int32",
                        "type" : "integer"
                      },
                      "port" : {
                        "default" : 0,
                        "description" : "The port that will be exposed by this service.",
                        "format" : "int32",
                        "type" : "integer"
                      },
                      "protocol" : {
                        "default" : "TCP",
                        "description" : "The IP protocol for this port. Supports \"TCP\", \"UDP\", and \"SCTP\". Default is TCP. Possible enum values: - `\"SCTP\"` is the SCTP protocol. - `\"TCP\"` is the TCP protocol. - `\"UDP\"` is the UDP protocol.",
                        "enum" : [
                          "SCTP",
                          "TCP",
                          "UDP"
                        ],
                        "type" : "string"
                      },
                      "target Port" : {
                        "allOf" : [
                          {
                            "description" : "Int Or String is a type that can hold an int32 or a string. When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type. This allows you to have, for example, a JSON field that can accept a name or number.",
                            "format" : "int-or-string",
                            "oneOf" : [
                              {
                                "type" : "integer"
                              },
                              {
                                "type" : "string"
                              }
                            ]
                          }
                        ],
                        "default" : {},
                        "description" : "Number or name of the port to access on the pods targeted by the service. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME. If this is a string, it will be looked up as a named port in the target Pod's container ports. If this is not specified, the value of the 'port' field is used (an identity map). This field is ignored for services with cluster IP=None, and should be omitted or set equal to the 'port' field. More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service"
                      }
                    },
                    "required" : [
                      "port"
                    ],
                    "type" : "object"
                  }
                ],
                "default" : {}
              },
              "type" : "array",
              "x-kubernetes-list-map-keys" : [
                "port",
                "protocol"
              ],
              "x-kubernetes-list-type" : "map",
              "x-kubernetes-patch-merge-key" : "port",
              "x-kubernetes-patch-strategy" : "merge"
            },
            "publish Not Ready Addresses" : {
              "description" : "publish Not Ready Addresses indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready. The primary use case for setting this field is for a Stateful Set's Headless Service to propagate SRV DNS records for its Pods for the purpose of peer discovery. The Kubernetes controllers that generate Endpoints and Endpoint Slice resources for Services interpret this to mean that all endpoints are considered \"ready\" even if the Pods themselves are not. Agents which consume only Kubernetes generated endpoints through the Endpoints or Endpoint Slice resources can safely assume this behavior.",
              "type" : "boolean"
            },
            "selector" : {
              "additionalProperties" : {
                "default" : "",
                "type" : "string"
              },
              "description" : "Route service traffic to pods with label keys and values matching this selector. If empty or not present, the service is assumed to have an external process managing its endpoints, which Kubernetes will not modify. Only applies to types Cluster IP, Node Port, and Load Balancer. Ignored if type is External Name. More info: https://kubernetes.io/docs/concepts/services-networking/service/",
              "type" : "object",
              "x-kubernetes-map-type" : "atomic"
            },
            "session Affinity" : {
              "description" : "Supports \"Client IP\" and \"None\". Used to maintain session affinity. Enable client IP based session affinity. Must be Client IP or None. Defaults to None. More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies Possible enum values: - `\"Client IP\"` is the Client IP based. - `\"None\"` - no session affinity.",
              "enum" : [
                "ClientIP",
                "None"
              ],
              "type" : "string"
            },
            "session Affinity Config" : {
              "allOf" : [
                {
                  "description" : "Session Affinity Config represents the configurations of session affinity.",
                  "properties" : {
                    "client IP" : {
                      "allOf" : [
                        {
                          "description" : "Client IP Config represents the configurations of Client IP based session affinity.",
                          "properties" : {
                            "timeout Seconds" : {
                              "description" : "timeout Seconds specifies the seconds of Client IP type session sticky time. The value must be >0 && <=86400(for 1 day) if Service Affinity == \"Client IP\". Default value is 10800(for 3 hours).",
                              "format" : "int32",
                              "type" : "integer"
                            }
                          },
                          "type" : "object"
                        }
                      ],
                      "description" : "client IP contains the configurations of Client IP based session affinity."
                    }
                  },
                  "type" : "object"
                }
              ],
              "description" : "session Affinity Config contains the configurations of session affinity."
            },
            "type" : {
              "description" : "type determines how the Service is exposed. Defaults to Cluster IP. Valid options are External Name, Cluster IP, Node Port, and Load Balancer. \"Cluster IP\" allocates a cluster-internal IP address for load-balancing to endpoints. Endpoints are determined by the selector or if that is not specified, by manual construction of an Endpoints object or Endpoint Slice objects. If cluster IP is \"None\", no virtual IP is allocated and the endpoints are published as a set of endpoints rather than a virtual IP. \"Node Port\" builds on Cluster IP and allocates a port on every node which routes to the same endpoints as the cluster IP. \"Load Balancer\" builds on Node Port and creates an external load-balancer (if supported in the current cloud) which routes to the same endpoints as the cluster IP. \"External Name\" aliases this service to the specified external Name. Several other fields do not apply to External Name services. More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types Possible enum values: - `\"Cluster IP\"` means a service will only be accessible inside the cluster, via the cluster IP. - `\"External Name\"` means a service consists of only a reference to an external name that kubedns or equivalent will return as a CNAME record, with no exposing or proxying of any pods involved. - `\"Load Balancer\"` means a service will be exposed via an external load balancer (if the cloud provider supports it), in addition to 'Node Port' type. - `\"Node Port\"` means a service will be exposed on one port of every node, in addition to 'Cluster IP' type.",
              "enum" : [
                "ClusterIP",
                "ExternalName",
                "LoadBalancer",
                "NodePort"
              ],
              "type" : "string"
            }
          },
          "type" : "object"
        }
      ],
      "default" : {},
      "description" : "Spec defines the behavior of a service. https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status"
    }
  },
  "type" : "object",
  "x-kubernetes-group-version-kind" : [
    {
      "group" : "",
      "kind" : "Service",
      "version" : "v1"
    }
  ]
}