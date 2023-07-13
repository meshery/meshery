export const publish_schema = {
  "type" : "object",
  "properties" : {
    "compatibility" : {
      "type" : "array",
      "items" : {
        "enum" : [
          "Kubernetes",
          "Argo CD",
          "AWS App Mesh",
          "Consul",
          "Fluentd",
          "Istio",
          "Jaeger",
          "Kuma",
          "Linkerd",
          "Network Service Mesh",
          "NGINX Service Mesh",
          "Open Service Mesh",
          "Prometheus",
          "Traefik Mesh"
        ],
        "type" : "string"
      },
      "uniqueItems" : true,
      "description" : "The list of compatible technologies.",
      "x-rjsf-style":6
    },
    "pattern_caveats" : {
      "type" : "string",
      "description" : "Caveats related to the pattern.",
      "x-rjsf-style": 12
    },
    "pattern_info" : {
      "type" : "string",
      "description" : "Additional information about the pattern.",
      "x-rjsf-style" : 12
    },
    "type" : {
      "type" : "string",
      "examples" : [
        "deployment",
        "observability",
        "resiliency",
        "scaling",
        "security",
        "traffic-management",
        "troubleshooting",
        "workloads"
      ],
      "description" : "The category of the pattern.",
      "x-rjsf-style": 6
    }
  },
  "required" : ["compatibility", "pattern_caveats", "pattern_info", "type"]
};
