export const publish_schema = {
  type : "object",
  properties : {
    type : {
      type : "string",
      title : "Type",
      enum : [
        "deployment",
        "observability",
        "resiliency",
        "scaling",
        "security",
        "traffic-management",
        "troubleshooting",
        "workloads",
      ],
      description : "The category of the pattern.",
    },
    compatibility : {
      type : "array",
      title : "Technology",
      items : {
        type : "string",
        enum : [
          "Traefik Mesh",
          "Trickster"
        ],
      },
      uniqueItems : true,
      description : "The list of compatible technologies.",
    },
    pattern_caveats : {
      type : "string",
      title : "Caveats and Consideration",
      format : "textarea",
      description : "Caveats related to the pattern.",
    },
    pattern_info : {
      type : "string",
      title : "Description",
      format : "textarea",
      description : "Additional information about the pattern.",
    },
  },
  required : ["compatibility", "pattern_caveats", "pattern_info", "type"],
};

export const publish_ui_schema = {
  "ui:order" : ["pattern_caveats", "pattern_info", "type", "compatibility"],
};
