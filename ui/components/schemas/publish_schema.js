export const publish_schema = {
  type: "object",
  properties: {
    type: {
      type: "string",
      title: "Type",
      examples: [
        "deployment",
        "observability",
        "resiliency",
        "scaling",
        "security",
        "traffic-management",
        "troubleshooting",
        "workloads",
      ],
      description: "The category of the pattern.",
    },
    compatibility: {
      type: "array",
      title: "Technology",
      items: {
        enum: [
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
          "Traefik Mesh",
        ],
        type: "string",
      },
      uniqueItems: false,
      description: "The list of compatible technologies.",
    },
    pattern_caveats: {
      type: "string",
      title: "Caveats and Consideration",
      format: "textarea",
      description: "Caveats related to the pattern.",
    },
    pattern_info: {
      type: "string",
      title: "Description",
      format: "textarea",
      description: "Additional information about the pattern.",
    },
  },
  required: ["compatibility", "pattern_caveats", "pattern_info", "type"],
};

export const publish_ui_schema = {
  "ui:order": ["pattern_caveats", "pattern_info", "type", "compatibility"],
};
