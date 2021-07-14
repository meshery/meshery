import KubernetesIcon from "../icons/KubernetesIcon.js"
import ServiceSpecificConfig from "../ServiceSpecificConfig.js"
import KubernetesConfigComp from "../ServiceSpecifComponents/Kubernetes.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import { Grid } from "@material-ui/core"


const KubernetesScreen = () => {

  const kubeserviceInfo = {
    name: "Kubernetes",
    logoComponent: KubernetesIcon,
    configComp :  <ServiceSpecificConfig components={[KubernetesConfigComp]} />

  }

  return (
    <Grid item xs={12} container justify="center" alignItems="center"> 
      <ServiceSwitch serviceInfo={kubeserviceInfo} /> 
    </Grid>
  )
}


export default KubernetesScreen
