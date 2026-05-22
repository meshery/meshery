export const CHART_COLORS = ["#3B697D", "#74A8BE", "#90B9CB", "#CBDEE6"];

export const topicsList = [
    { label: "Deployment", value: "deployment" },
    { label: "Traffic-management", value: "traffic-management" },
    { label: "Security", value: "security" },
    { label: "Workloads", value: "workloads" },
    { label: "Observability", value: "observability" },
    { label: "Troubleshooting", value: "troubleshooting" },
    { label: "Scaling", value: "scaling" },
    { label: "Resiliency", value: "resiliency" },
];

// Remote providers supported by the Meshery Docker Extension.
// To add or change entries, edit install/Makefile.core.mk and run
// `make sync-provider-defaults`. Do not hand-edit the array body.
export const REMOTE_PROVIDERS = [
    // BEGIN AUTO-GENERATED-FROM-MAKEFILE
    {
        name: "Meshery",
        url: "https://cloud.meshery.io",
    },
    {
        name: "Layer5",
        url: "https://cloud.layer5.io",
    },
    // END AUTO-GENERATED-FROM-MAKEFILE
];

export const SELECTED_REMOTE_PROVIDER =
    REMOTE_PROVIDERS.find(
        (provider) =>
            provider.name === process.env.REACT_APP_REMOTE_PROVIDER_NAME,
    ) ?? REMOTE_PROVIDERS[0];

export const SELECTED_PROVIDER_NAME = SELECTED_REMOTE_PROVIDER.name;

export const providerUrl = SELECTED_REMOTE_PROVIDER.url;

export const MESHMAP = "meshmap";
