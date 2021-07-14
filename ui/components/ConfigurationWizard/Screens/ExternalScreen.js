import GrafanaIcon from "../icons/GrafanaIcon.js"
import PrometheusIcon from "../icons/PrometheusIcon.js"
import ServiceSpecificConfig from "../ServiceSpecificConfig.js"
import ExternalsConfig from "../ServiceSpecifComponents/ExternalsConfig.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"


const KubernetesScreen = () => {

  const externalComponentsInfo = [{
    name: "Grafana",
    logoComponent: GrafanaIcon,
    configComp :  <ServiceSpecificConfig components={[ExternalsConfig]} />

  },
  {
    name: "Prometheus",
    logoComponent: PrometheusIcon,
    configComp :  <ServiceSpecificConfig components={[ExternalsConfig]} />

  }
  ]

  return (
    <>
      {externalComponentsInfo && externalComponentsInfo.map(info => (
        <Grid item lg={6} md={12} sm={12} container justify="center" alignItems="center">
          <Grid item>
            <ServiceSwitch serviceInfo={info} /> 
          </Grid>
        </Grid>
      ))}
    </>
  )
}


export default KubernetesScreen
