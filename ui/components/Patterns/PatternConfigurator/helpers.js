import React from "react";
import { AddCircle, BuildRounded, DirectionsCar, Filter, SimCard, SupervisedUserCircle, TouchApp } from "@mui/icons-material";
import ExploreIcon from '@mui/icons-material/Explore';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import ListAltIcon from '@mui/icons-material/ListAlt';
import LockIcon from '@mui/icons-material/Lock';

/**
 * nameToIcon returns icons for the rjsf form titles
 *
 * @param {string} name Tooltip name
 * @param {string} color default value is primary
 * @returns CustomIconButton
 */
const NameToIcon = ({ name, color = "#607D8B", ...props }) => {

  const CustomIcon = ({ Icon }) => <Icon style={{ color }}  {...props} />;

  switch (name) {
    //core
    case "Application": return <CustomIcon Icon={TouchApp} />;
    case "Kubernetes Service": return <CustomIcon Icon={BuildRounded} />;
    // Istio
    case "AuthorizationPolicy": return <CustomIcon Icon={LockIcon} />;
    case "DestinationRule": return <CustomIcon Icon={ExploreIcon} />;
    case "EnvoyFilter": return <CustomIcon Icon={Filter} />;
    case "Gateway": return <CustomIcon Icon={ListAltIcon} />;
    case "PeerAuthentication": return <CustomIcon Icon={FileCopyIcon} />;
    case "Sidecar": return <CustomIcon Icon={DirectionsCar} />;
    case "VirtualService": return <CustomIcon Icon={SupervisedUserCircle} />;
    case "WorkloadEntry": return <CustomIcon Icon={SimCard} />;
    // default
    default: return <CustomIcon Icon={AddCircle} />;
  }
}
export default NameToIcon

/**
 * getHumanReadablePatternServiceName takes in the pattern service metadata and returns
 * the readable name of the service
 *
 * @param {*} item pattern service component
 * @returns {string} service name
 */
 export function getHumanReadablePatternServiceName(item) {
    return item?.metadata?.["display.ui.meshery.io/name"];
  }
  
  /**
   * getPatternServiceName takes in the pattern service metadata and returns
   * the name of the service
   *
   * @param {*} item pattern service component
   * @param {boolean} includeDisplayName if set to true, display name is checked first
   * @returns {string} service name
   */
  export function getPatternServiceName(item, includeDisplayName = true) {
    if (includeDisplayName)
      return (
        item?.metadata?.["display.ui.meshery.io/name"] ||
        item?.oam_definition?.metadata?.name ||
        getPatternAttributeName(item) ||
        "NA"
      );
  
    return item?.oam_definition?.metadata?.name || "NA";
  }

  /**
* @param {object} obj check obj is empty or not
* @returns {boolean} true or false based on object is empty
*/

export const isEmptyObj= (obj)=> {
  for (let i in obj) return true;
  return false;
}


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