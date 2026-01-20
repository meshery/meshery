interface AdapterConfig {
  name: string;
  label: string;
  imageSrc: string;
  description: string;
  defaultPort: number;
  enabled: boolean;
  url: string;
}

export interface AdaptersListType {
  [key: string]: AdapterConfig;
}

const adaptersDescription = (adapterName: string): string => {
  return `Deploy the Meshery Adapter for ${adapterName} in order to enable deeper lifecycle management of ${adapterName}.`;
};

/*
 * adaptersList.name  -> name of the adapter to display on the card.
 * adaptersList.label -> used as a payload for adapter deployment (like an adapterId).
 */
export const adaptersList: AdaptersListType = {
  ISTIO: {
    name: 'Istio',
    label: 'meshery-istio',
    imageSrc: '/static/img/adapters/istio.svg',
    description: adaptersDescription('Istio'),
    defaultPort: 10000,
    enabled: false,
    url: '',
  },
  LINKERD: {
    name: 'Linkerd',
    label: 'meshery-linkerd',
    imageSrc: '/static/img/adapters/linkerd.svg',
    description: adaptersDescription('Linkerd'),
    defaultPort: 10001,
    enabled: false,
    url: '',
  },
  CONSUL: {
    name: 'Consul',
    label: 'meshery-consul',
    imageSrc: '/static/img/adapters/consul.svg',
    description: adaptersDescription('Consul'),
    defaultPort: 10002,
    enabled: false,
    url: '',
  },
  NETWORK_SERVICE_MESH: {
    name: 'Network Service Mesh',
    label: 'meshery-nsm',
    imageSrc: '/static/img/adapters/networkservicemesh.svg',
    description: adaptersDescription('Network Service Mesh'),
    defaultPort: 10004,
    enabled: false,
    url: '',
  },
  TRAEFIK_MESH: {
    name: 'Traefik Mesh',
    label: 'meshery-traefik-mesh',
    imageSrc: '/static/img/adapters/traefik_mesh.svg',
    description: adaptersDescription('Traefik Mesh'),
    defaultPort: 10006,
    enabled: false,
    url: '',
  },
  KUMA: {
    name: 'Kuma',
    label: 'meshery-kuma',
    imageSrc: '/static/img/adapters/kuma.svg',
    description: adaptersDescription('Kuma'),
    defaultPort: 10007,
    enabled: false,
    url: '',
  },
  NGINX_SERVICE_MESH: {
    name: 'NGINX Service Mesh',
    label: 'meshery-nginx-sm',
    imageSrc: '/static/img/adapters/nginx.svg',
    description: adaptersDescription('NGINX Service Mesh'),
    defaultPort: 10010,
    enabled: false,
    url: '',
  },
  CILIUM_SERVICE_MESH: {
    name: 'Cilium Service Mesh',
    label: 'meshery-cilium',
    imageSrc: '/static/img/adapters/cilium_service_mesh.svg',
    description: adaptersDescription('Cilium Service Mesh'),
    defaultPort: 10012,
    enabled: false,
    url: '',
  },
  NIGHTHAWK: {
    name: 'Nighthawk',
    label: 'meshery-nighthawk',
    imageSrc: '/static/img/adapters/nighthawk-logo.svg',
    description: adaptersDescription('Performance Characterization by Meshery Nighthawk'),
    defaultPort: 10013,
    enabled: false,
    url: '',
  },
};

export const ADAPTER_STATUS = {
  ENABLED: 'ENABLED',
  DISABLED: 'DISABLED',
} as const;

export type AdapterStatusType = (typeof ADAPTER_STATUS)[keyof typeof ADAPTER_STATUS];
