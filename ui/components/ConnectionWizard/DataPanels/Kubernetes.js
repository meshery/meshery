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
import { updateProgress, updateK8SConfig } from "../../../lib/store";
import { pingKubernetesWithNotification, reconfigureKubernetes } from "../helpers/kubernetesHelpers";

const styles = theme => ({

  infoContainer: {
    width: "20rem",
    height: "15rem",
    padding: "1rem 1rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.20)",
    borderRadius: '1rem',
  },
  infoTitle: {
    color: "#647881",
    width: "3rem",
    background: "#F1F3F4",
    padding: ".5rem 5rem .75rem 1.5rem",
    borderRadius: "0.25rem",
    fontSize: ".8rem",
  },
  infoLabel: {
    fontSize: ".9rem",
    color: theme.palette.text.primary,
    marginRight: "1rem",
  },
  infoData: {
    fontSize: ".9rem",
    color: theme.palette.text.secondary,
  },


})

const chipStyles = (theme) => ({
  chipIcon: {
    width: theme.spacing(2.5)
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
})

const KubernetesChip = withStyles(chipStyles)(({classes, handleKubernetesClick, handleKubernetesDelete,  label}) => (

  <Chip
    label={label}
    onDelete={handleKubernetesDelete}
    onClick={handleKubernetesClick}
    icon={<img src="/static/img/kubernetes.svg" className={classes.chipIcon} />}
    className={classes.chip}
    key="k8s-key"
    variant="outlined"
  />
))


const KubernetesPingSnackbarAction = (closeSnackbar) => (key) => (
  <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
    <CloseIcon />
  </IconButton>
)


const KubernetesDataPanel = ({clusterInformation, classes, updateProgress, enqueueSnackbar, closeSnackbar, setClusterInformation, setIsConnected, updateK8SConfig}) => {

  const handleKubernetesDelete = () => {
    updateProgress({showProgress: true})
    reconfigureKubernetes(updateProgress, enqueueSnackbar, KubernetesPingSnackbarAction(closeSnackbar), updateK8SConfig)    
  }
  return (
    <Grid container className={classes.infoContainer} xs={10}>

      <Grid item xs={12}>
        <KubernetesChip 
          handleKubernetesClick={() => pingKubernetesWithNotification(updateProgress, enqueueSnackbar, KubernetesPingSnackbarAction(closeSnackbar))} 
          handleKubernetesDelete = {handleKubernetesDelete}
          label="Kubernetes" 
        />
      </Grid>

      <Grid item xs={12} style={{marginBottom: "1rem"}}>
        <Typography className={classes.infoTitle}>Status</Typography>
      </Grid>

      <Grid item xs={12} container>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Current-Context: {clusterInformation.contextName}</Typography>
          <Typography className={classes.infoData}>
            {clusterInformation.inClusterConfig ? "Using In Cluster Config" : "Using out of cluster config"}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Kubernetes Version</Typography>
          <Typography className={classes.infoData}>
            {clusterInformation.serverVersion && clusterInformation.serverVersion}
          </Typography>
        </Grid>
      </Grid>
    </Grid>

  )
}


const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
});

export default withStyles(styles)(connect(null, mapDispatchToProps)(withSnackbar(KubernetesDataPanel)))
