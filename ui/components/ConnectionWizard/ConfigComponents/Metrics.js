/* eslint-disable no-unused-vars */
/* eslint-disable react/display-name */
import {
  withStyles,
  TextField,
  Grid,
  Button,
  IconButton,
} from "@material-ui/core/";
import { useState } from "react";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig } from "../../../lib/store";
import ReactSelectWrapper from "../../ReactSelectWrapper"
import { handleGrafanaConfigure, handlePrometheusConfigure } from "../helpers/metrics";
import CloseIcon from "@material-ui/icons/Close";
import { useNotification } from "../../../utils/hooks/useNotification";


const styles = () => ({})

const configurationNotificationAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)

const MetricsConfig = ({
  componentName, prometheusScannedUrls, grafanaScannedUrls,updatePrometheusConfig, updateGrafanaConfig, updateProgress
}) => {
  const { notify } = useNotification()
  const  handleConfigurationSubmit = () => {

    if (componentName === "Grafana") handleGrafanaConfigure(notify, url,apiKey, updateProgress, updateGrafanaConfig  )
    if (componentName === "Prometheus") handlePrometheusConfigure(notify, url, updateProgress, updatePrometheusConfig  )
  }
  const getOptions = () => {
    if (componentName === "Grafana")
      return grafanaScannedUrls?.map((graf) => ({ label : graf, value : graf }))
    if (componentName === "Prometheus")
      return prometheusScannedUrls?.map((prom) => ({ label : prom, value : prom }))
  }

  const [url, setUrl] = useState()
  const [apiKey, setApiKey] = useState("")



  return (
    <>
      {componentName === "Prometheus" &&
      <Grid item xs={12} style={{ height : "12.2rem" }}>
        <Grid item xs={12} style={{ marginTop : "2rem", cursor : "pointer" }}>
          <ReactSelectWrapper
            onChange={(select) => setUrl(select
              ? select.value
              : "")}
            options={getOptions()}
            value={{ label : url, value : url }}
            label={`${componentName} Base URL`}
            error={false}
            placeholder={`Address of ${componentName} Server`}
            noOptionsMessage={`No ${componentName} servers discovered`}
          />
        </Grid>
        <Grid item xs={12} style={{ textAlign : "center", marginTop : "2rem" }}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={handleConfigurationSubmit}
          >
          Connect
          </Button>
        </Grid>
      </Grid>
      }
      {componentName === "Grafana" &&
      <Grid item xs={12} style={{ height : "13rem" }}>
        <Grid item xs={12} style={{ marginTop : "1.2rem" }}>
          <ReactSelectWrapper
            onChange={(select) => setUrl(select
              ? select.value
              : "")}
            options={getOptions()}
            value={{ label : url, value : url }}
            label={`${componentName} Base URL`}
            error={false}
            placeholder={`Address of ${componentName} Server`}
            noOptionsMessage={`No ${componentName} servers discovered`}
          />
        </Grid>
        <Grid item xs={12} >
          <TextField
            id="external-config-api-key-input"
            name="externalConfigapiKey"
            label="API key"
            fullWidth
            margin="normal"
            variant="outlined"
            // disabled={inClusterConfigForm === true}
          />
          <Grid item xs={12} style={{ marginTop : "0.8rem", marginBottom : "0.8rem", textAlign : "center" }}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              size="large"
              onClick={handleConfigurationSubmit}
            >
            Connect
            </Button>
          </Grid>
        </Grid>
      </Grid>
      }
    </>
  )
}

const mapDispatchToProps = (dispatch) => ({ updateGrafanaConfig : bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig : bindActionCreators(updatePrometheusConfig, dispatch),
  updateProgress : bindActionCreators(updateProgress, dispatch), });

export default withStyles(styles)(connect(null, mapDispatchToProps)(MetricsConfig))
