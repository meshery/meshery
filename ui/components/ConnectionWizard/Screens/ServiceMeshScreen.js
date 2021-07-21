import ConsulIcon from "../icons/ConsulIcon.js"
import OpenServiceMeshIcon from "../icons/OpenServiceMeshIcon.js"
import LinkerdIcon from "../icons/LinkerdIcon.js"
import ServiceMeshConfig from "../SwitchConfigComponents/ServiceMesh"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"


const KubernetesScreen = () => {

  const serviceMeshComponents = [{
    name: "OpenServiceMesh",
    logoComponent:  OpenServiceMeshIcon,
    configComp :  ServiceMeshConfig
  }, {
    name: "Consul",
    logoComponent: ConsulIcon,
    configComp :  ServiceMeshConfig
  },{
    name: "Linkerd",
    logoComponent: LinkerdIcon,
    configComp :  ServiceMeshConfig
  }]


  return (
    <>
      {serviceMeshComponents.map(comp =>
        <Grid item lg={4} md={6} sm={12}>
          <ServiceSwitch serviceInfo={comp} /> 
        </Grid>
      )}
    </>
  )
}


export default KubernetesScreen

