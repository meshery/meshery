export const publish_schema = {
  "type" : "object",
  "properties" : {
    "compatibility" : {
      "type" : "array",
      "items" : {
        "type" : "string",
        "enum" : [
          "Istio",
          "Linkerd",
          "App Mesh",
          "OSM",
          "Nginx",
          "Kuma",
          "Consul",
          "NSM",
          "Traefik"
        ]
      },
      "uniqueItems" : true,
      "description" : "Select the service mesh that this pattern is compatible with."
    },
    "pattern_caveats" : {
      "type" : "string",
      "description" : "Enter any caveats or limitations of the pattern."
    },
    "pattern_info" : {
      "type" : "string",
      "description" : "Enter any additional information about the pattern."
    },
    "type" : {
      "type" : "string",
      "description" : "Select the category that best describes this pattern.",
      "enum" : [
        "Deployment",
        "Observability",
        "Resiliency",
        "Scaling",
        "Security",
        "Traffic Management",
        "Troubleshooting",
        "Workloads"
      ]
    }
  }
}