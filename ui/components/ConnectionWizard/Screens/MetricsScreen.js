/* eslint-disable no-unused-vars */
import GrafanaIcon from "../icons/GrafanaIcon"
import PrometheusIcon from "../icons/PrometheusIcon"
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import MetricsConfig from "../SwitchConfigComponents/Metrics.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"
import VerticalCarousel from "../../VerticalCarousel/VerticalCarousel"
import MetricsDataPanel from "../DataPanels/Metrics"
import {useEffect, useRef, useState} from "react"
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig } from "../../../lib/store.js";
import {fetchPromGrafanaScanData, verifyGrafanaConnection, verifyPrometheusConnection} from "../helpers/metrics"

const MetricsScreen = ({ grafana, prometheus}) => {

  const [isGrafanaConnected, setIsGrafanaConnected] = useState(false)
  const [isPrometheusConnected, setIsPrometheusConnected] = useState(false)
  const [metricsScanUrls, setMetricsScanUrls] = useState({grafana: [], prometheus: []})

  const metricsComponents = [{
    name: "Grafana",
    logoComponent: GrafanaIcon,
    configComp :  <MetricsConfig componentName="Grafana" grafanaScannedUrls={metricsScanUrls.grafana}/>

  },
  {
    name: "Prometheus",
    logoComponent: PrometheusIcon,
    configComp : <MetricsConfig componentName="Prometheus" prometheusScannedUrls={metricsScanUrls.prometheus}/>

  }
  ]

  const getConnectionStatus = (name) => {
    if(name === "Grafana") return isGrafanaConnected
    if(name === "Prometheus") return isPrometheusConnected
  }

  useEffect(() => {
    console.log("Prometheus effect", prometheus)
    verifyPrometheusConnection(prometheus.prometheusURL)
      .then(res => {
        console.log(res)
        if( typeof res !== "undefined") setIsPrometheusConnected(true)
        else setIsPrometheusConnected(false)
      })
      .catch(err => {
        setIsPrometheusConnected(false)
        console.log(err)
      })

  }, [prometheus.ts])

  useEffect(() => {
    console.log("Greafana effect", grafana)
    verifyGrafanaConnection(grafana.grafanaURL)
      .then(res => {
        if( typeof res !== "undefined") setIsGrafanaConnected(true)
        else setIsGrafanaConnected(false)
      })
      .catch(err => setIsGrafanaConnected(false))
  }, [grafana.ts])
   
  useEffect(() => {
    fetchPromGrafanaScanData()
      .then(res => setMetricsScanUrls(res))
      .catch(console.log)

  },[])



  const [activeIndex, setActiveIndex] = useState(0);

  const itemsToBeRendered = metricsComponents.map(comp => {
    return(
      <Grid item lg={4} md={6} sm={12}>
        <ServiceSwitch serviceInfo={comp} isConnected={getConnectionStatus(comp.name)} /> 
      </Grid>
    ) 
  })

  const itemToDisplay = (items, curIndex) => {
    return  items[curIndex] ? items[curIndex] : null
  }



  return (
    <Grid xs={12} container>
      <Grid item xs={4} container justify="center">
        <VerticalCarousel item=
          {itemToDisplay(itemsToBeRendered, activeIndex)} setActiveIndex={setActiveIndex}/>
      </Grid>
      <Grid item xs={8} container justify="center">
        <MetricsDataPanel isConnected={itemToDisplay(metricsComponents, activeIndex).name === "Grafana" ? isGrafanaConnected : isPrometheusConnected} 
          componentName={itemToDisplay(metricsComponents, activeIndex).name}
        />
      </Grid>
    </Grid>
  )
}

const mapDispatchToProps = (dispatch) => ({
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
});

const mapStateToProps = (state) => {
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return {
    grafana,
    prometheus
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MetricsScreen)

