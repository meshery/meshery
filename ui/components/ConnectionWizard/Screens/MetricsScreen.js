import GrafanaIcon from "../icons/GrafanaIcon.js"
import PrometheusIcon from "../icons/PrometheusIcon.js"
import MetricsConfig from "../SwitchConfigComponents/Metrics"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"


const MetricsScreen = () => {

  const externalComponentsInfo = [{
    name: "Grafana",
    logoComponent: GrafanaIcon,
    configComp :  <MetricsConfig />

  },
  {
    name: "Prometheus",
    logoComponent: PrometheusIcon,
    configComp : <MetricsConfig/>

  }
  ]

  return (
    <>
      {externalComponentsInfo && externalComponentsInfo.map(info => (
        <Grid item lg={6} md={12} sm={12} container justify="center" alignItems="center">
          <Grid item>
            <ServiceSwitch serviceInfo={info} isConnected={false} /> 
          </Grid>
        </Grid>
      ))}
    </>
  )
}


export default MetricsScreen
