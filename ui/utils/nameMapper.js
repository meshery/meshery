
/**
* @param {string} name name of the mesh
* @returns {object} details of the mesh: {name:"",img:"",color:""}
*/

export const getMeshProperties=(name)=> {
    switch (name) {
      case "istio":
        return { name, image: "/static/img/istio.svg", color: "#466BB0" };
      case "linkerd":
        return { name, image: "/static/img/linkerd.svg", color: "#2beda7" };
      case "kuma":
        return { name, image: "/static/img/kuma.svg", color: "#291953" };
      case "nginx_service_mesh":
        return { name, image: "/static/img/nginx.svg", color: "#009639" };
      case "open_service_mesh":
        return { name, image: "/static/img/openservicemesh.svg", color: "#6ED3B9" };
      case "cilium_service_mesh":
        return {
          name,
          img: "/static/img/cilium_service_mesh.svg",
          color: "#6B91C7"
        };
      case "citrix":
        return {
          name,
          img: "/static/img/citrix_service_mesh.svg",
          color: "#466BB0"
        };
      case "traefik_mesh":
        return { name, image: "/static/img/traefik_mesh.svg", color: "#9D0FB0" };
      case "consul":
        return { name, image: "/static/img/consul.svg", color: "#D62783" };
      case "app_mesh":
        return { name, image: "/static/img/app_mesh.svg", color: "#F49322" };
      case "core":
        return { name, image: "/static/img/meshery-logo.png", color: "#00B39F" };
      case "kubernetes":
        return { name, image: "/static/img/kubernetes.svg", color: "#293CDA" };
      default:
        return {};
    }
  }