/* eslint-disable react/display-name */
/* eslint-disable no-unused-vars */
import CloseIcon from "@material-ui/icons/Close";
import {
  withStyles,
  Typography,
  Grid,
  Chip,
  IconButton,
} from "@material-ui/core/";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import {
  updateProgress, updateK8SConfig, updateGrafanaConfig, updatePrometheusConfig
} from "../../../lib/store";
import { pingAdapterWithNotification } from "../helpers/serviceMeshes"
import { deleteMetricsComponentConfig, pingGrafanaWithNotification, pingPrometheusWithNotification } from "../helpers/metrics";

const styles = theme => ({

  infoContainer : {
    width : "20rem",
    height : "15rem",
    padding : "1rem 1rem",
    boxShadow : "0px 1px 6px 1px rgba(0,0,0,0.20)",
    borderRadius : '1rem',
  },
  infoTitle : {
    color : "#647881",
    width : "3rem",
    background : "#F1F3F4",
    padding : ".5rem 5rem .75rem 1.5rem",
    borderRadius : "0.25rem",
    fontSize : ".8rem",
  },
  infoLabel : {
    fontSize : ".9rem",
    color : theme.palette.text.primary,
    marginRight : "1rem",
    fontWeight : 500
  },
  infoData : { fontSize : ".9rem",
    color : theme.palette.text.secondary, },


})

const chipStyles = (theme) => ({ chipIcon : { width : theme.spacing(2.5) },
  chip : { marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1), }, })


const AdapterChip = withStyles(chipStyles)(({
  classes, handleAdapterClick, isConnected, handleAdapterDelete, componentName
}) => {
  let image = componentName === "Grafana" ? "/static/img/grafana_icon.svg"
    : "/static/img/prometheus_logo_orange_circle.svg";
  let logoIcon = <img src={image} className={classes.chipIcon} />;
  return (
    <Chip
      label={componentName}
      onClick={handleAdapterClick}
      onDelete={isConnected
        ? handleAdapterDelete
        : null}
      icon={logoIcon}
      className={classes.chip}
      key={`adapters-${11}`}
      variant={isConnected
        ? "outlined"
        : "default"}
    />
  )
})


const AdapterPingSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)



const MetricsDataPlane = ({
  classes, updateProgress, enqueueSnackbar,grafana, prometheus, closeSnackbar, isConnected, componentName, updateGrafanaConfig, updatePrometheusConfig
}) => {

  const handleDeleteAdapter = () => {

    updateProgress({ showProgress : true });

    const successCb = (result) => {
      updateProgress({ showProgress : false });
      if (typeof result !== "undefined") {

        if (componentName === "Grafana")
          console.log("Updating grafana config")
        updateGrafanaConfig({ grafana : {
          grafanaURL : "",
          grafanaAPIKey : "",
          grafanaBoardSearch : "",
          grafanaBoards : [],
          selectedBoardsConfigs : [],
        }, })


        if (componentName === "Prometheus")
          console.log("Updating prometheus config")

        updatePrometheusConfig({ prometheus : { prometheusURL : "",
          selectedPrometheusBoardsConfigs : [], }, });

        enqueueSnackbar(`${componentName} was disconnected!` , { variant : "success",
          autoHideDuration : 2000,
          action : AdapterPingSnackbarAction })
      }
    }

    const errorCb = (error) => {
      updateProgress({ showProgress : false });
      enqueueSnackbar(`${componentName} could not be disconnected!: ${error}` , { variant : "error",
        autoHideDuration : 2000,
        action : AdapterPingSnackbarAction })
    }

    deleteMetricsComponentConfig(componentName)(successCb, errorCb)
  }


  const handleAdapterClick = () => {
    if (componentName === "Prometheus") pingPrometheusWithNotification(updateProgress, AdapterPingSnackbarAction(closeSnackbar), enqueueSnackbar)
    if (componentName === "Grafana") pingGrafanaWithNotification(updateProgress, AdapterPingSnackbarAction(closeSnackbar), enqueueSnackbar)
  }

  return (
    <Grid container className={classes.infoContainer} xs={10}>

      <Grid item xs={12}>
        <AdapterChip
          handleAdapterClick={handleAdapterClick}
          isConnected={isConnected}
          handleAdapterDelete= {handleDeleteAdapter}
          componentName={componentName}
        />
      </Grid>

      <Grid item xs={12}>
        <Typography className={classes.infoTitle}>Details</Typography>
      </Grid>

      <Grid item xs={12} container>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Adapter Server Location:</Typography>
          <Typography className={classes.infoData}>{componentName ==="Grafana"
            ? grafana.grafanaURL
            : prometheus.prometheusURL}</Typography>
        </Grid>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Adapter Version</Typography>
          <Typography className={classes.infoData}>

          </Typography>
        </Grid>
      </Grid>
    </Grid>

  )
}


const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig : bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig : bindActionCreators(updatePrometheusConfig, dispatch), });

const mapStateToProps = (state) => {
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return { grafana,
    prometheus };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MetricsDataPlane)))
