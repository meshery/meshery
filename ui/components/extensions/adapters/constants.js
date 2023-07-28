const adaptersDescription = (adapterName) => {
  return `Deploy the Meshery Adapter for ${adapterName} in order to enable deeper lifecycle management of ${adapterName}.`;
}

/*
 * adaptersList.name  -> name of the adapter to display on the card.
 * adaptersList.label -> used as a payload for adapter deployment (like an adapterId).
*/
export const adaptersList = {
  "ISTIO" : {
    name : "Istio",
    label : "meshery-istio",
    imageSrc : "/static/img/istio.svg",
    description : adaptersDescription("Istio"),
    defaultPort : 10000,
    enabled : false,
    url : "",
  },
  "LINKERD" : {
    name : "Linkerd",
    label : "meshery-linkerd",
    imageSrc : "/static/img/linkerd.svg",
    description : adaptersDescription("Linkerd"),
    defaultPort : 10001,
    enabled : false,
    url : "",
  },
  "CONSUL" : {
    name : "Consul",
    label : "meshery-consul",
    imageSrc : "/static/img/consul.svg",
    description : adaptersDescription("Consul"),
    defaultPort : 10002,
    enabled : false,
    url : "",
  },
  "NETWORK_SERVICE_MESH" : {
    name : "Network Service Mesh",
    label : "meshery-nsm",
    imageSrc : "/static/img/networkservicemesh.svg",
    description : adaptersDescription("Network Service Mesh"),
    defaultPort : 10004,
    enabled : false,
    url : "",
  },
  "APP_MESH" : {
    name : "App Mesh",
    label : "meshery-app-mesh",
    imageSrc : "/static/img/app_mesh.svg",
    description : adaptersDescription("App Mesh"),
    defaultPort : 10005,
    enabled : false,
    url : "",
  },
  "TRAEFIK_MESH" : {
    name : "Traefik Mesh",
    label : "meshery-traefik-mesh",
    imageSrc : "/static/img/traefik_mesh.svg",
    description : adaptersDescription("Traefik Mesh"),
    defaultPort : 10006,
    enabled : false,
    url : "",
  },
  "KUMA" : {
    name : "Kuma",
    label : "meshery-kuma",
    imageSrc : "/static/img/kuma.svg",
    description : adaptersDescription("Kuma"),
    defaultPort : 10007,
    enabled : false,
    url : "",
  },
  // TODO: Need to add icon for this.
  // "meshery-cpx": {
  //   name: "Meshery Cpx",
  // label: "meshery-cpx",
  //   imageSrc: "/static/img/",
  //   description: adaptersDescription("Meshery CPX"),
  //   defaultPort: 10008,
  //   enabled: false,
  //   url: "",
  // },
<<<<<<< HEAD

=======
>>>>>>> 02e0c85c1 (Remove Open Service Mesh Adapter Card)
  "NGINX_SERVICE_MESH" : {
    name : "Nginx Service Mesh",
    label : "meshery-nginx-sm",
    imageSrc : "/static/img/nginx.svg",
    description : adaptersDescription("Nginx Serice Mesh"),
    defaultPort : 10010,
    enabled : false,
    url : "",
  },
  "CILIUM_SERVICE_MESH" : {
    name : "Cilium Service Mesh",
    label : "meshery-cilium",
    imageSrc : "/static/img/cilium_service_mesh.svg",
    description : adaptersDescription("Cilium Service Mesh"),
    defaultPort : 10012,
    enabled : false,
    url : "",
  },
};