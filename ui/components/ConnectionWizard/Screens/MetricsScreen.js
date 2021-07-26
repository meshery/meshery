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
import {useEffect, useState} from "react"
import { updateProgress } from "../../../lib/store.js";

const MetricsScreen = ({ updateProgress}) => {

  const metricsComponents = [{
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



  const [activeIndex, setActiveIndex] = useState(0);

  const itemsToBeRendered = metricsComponents.map(comp => {
    return(
      <Grid item lg={4} md={6} sm={12}>
        <ServiceSwitch serviceInfo={comp} isConnected={true} /> 
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
        <MetricsDataPanel  />
      </Grid>
    </Grid>
  )
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  return {
    meshAdapters,
    meshAdaptersts,
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MetricsScreen)

