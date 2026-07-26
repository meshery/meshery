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

// Extension Point: Add more remote providers here as needed.
// This is a list of remote providers that are supported by the Meshery Docker Extension.
export const REMOTE_PROVIDERS = [
    // BEGIN AUTO-GENERATED from install/providers.env - run `make providers-propagate`
    {
        name: "Meshery",
        url: "https://cloud.meshery.io",
    },
    {
        name: "Intel",
        url: "https://perf.smp-spec.io",
    },
    {
        name: "Layer5",
        url: "https://cloud.layer5.io",
    },
    {
        name: "TATA Labs",
        url: "https://platform.tata-consulting.co.uk",
    },
    {
        name: "Cisco ET&I",
        url: "https://collab.eti.cisco.com",
    },
    {
        name: "Metabit Trading",
        url: "https://kickstart.metabit.com",
    },
    {
        name: "OD10 Ventures",
        url: "https://provider.od10.in",
    },
    // END AUTO-GENERATED
];

export const SELECTED_REMOTE_PROVIDER =
    REMOTE_PROVIDERS.find(
        (provider) =>
            provider.name === process.env.REACT_APP_REMOTE_PROVIDER_NAME,
    ) ?? REMOTE_PROVIDERS[0];

export const SELECTED_PROVIDER_NAME = SELECTED_REMOTE_PROVIDER.name;

export const providerUrl = SELECTED_REMOTE_PROVIDER.url;

export const MESHMAP = "meshmap";
