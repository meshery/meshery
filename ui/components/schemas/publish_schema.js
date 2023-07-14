import _ from "lodash";
import { getMeshModels } from "../../api/meshmodel";

export let publish_schema = null;

getMeshModels()
  .then(({ models }) => {
    const model_names = _.uniq(models?.map((model) => model.displayName));
    publish_schema = _.set(_.cloneDeep(json_schema), "properties.compatibility.items.enum", model_names);
  })
  .catch((err) => {
    console.error(err);
    publish_schema = json_schema;
  });

const json_schema = {
  type: "object",
  properties: {
    compatibility: {
      type: "array",
      title: "Technology",
      items: {
        enum: ["istio", "linkerd"],
        type: "string",
      },
      uniqueItems: true,
      description: "The list of compatible technologies.",
      "x-rjsf-grid-area": 6,
    },
    pattern_caveats: {
      type: "string",
      title: "Caveats and Considerations",
      description: "Caveats related to the design.",
      format : "textarea",
      "x-rjsf-grid-area": 12,
    },
    pattern_info: {
      type: "string",
      title: "Description",
      description: "Additional information about the design.",
      format : "textarea",
      "x-rjsf-grid-area": 12,
    },
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
      "x-rjsf-grid-area": 6,
    },
  },
  required: ["compatibility", "pattern_caveats", "pattern_info", "type"],
};

export const publish_ui_schema = {
  "ui:order": ["type", "compatibility", "pattern_caveats", "pattern_info"],
};

