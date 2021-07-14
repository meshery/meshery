import ConsulIcon from "../icons/ConsulIcon.js"
import OpenServiceMeshIcon from "../icons/OpenServiceMeshIcon.js"
import LinkerdIcon from "../icons/LinkerdIcon.js"
import ServiceSpecificConfig from "../ServiceSpecificConfig.js"
import ServiceMeshConfigComp from "../ServiceSpecifComponents/ServiceMeshConfigComp.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"


const KubernetesScreen = () => {

  const serviceMeshComponents = [{
    name: "OpenServiceMesh",
    logoComponent:  OpenServiceMeshIcon,
    configComp :  <ServiceSpecificConfig components={[ServiceMeshConfigComp]} />
  }, {
    name: "Consul",
    logoComponent: ConsulIcon,
    configComp :  <ServiceSpecificConfig components={[ServiceMeshConfigComp]} />
  },{
    name: "Linkerd",
    logoComponent: LinkerdIcon,
    configComp :  <ServiceSpecificConfig components={[ServiceMeshConfigComp]} />
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

