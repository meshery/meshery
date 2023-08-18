/* eslint-disable react/display-name */
/* eslint-disable no-unused-vars */
import {
  withStyles,
  Typography,
  Grid,
} from "@material-ui/core/";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { updateProgress, updateK8SConfig } from "../../../lib/store";
import AdapterChip from "./AdapterChip"
import { deleteKubernetesConfig, pingKubernetes } from "../helpers/kubernetesHelpers"
import {
  successHandlerGenerator, errorHandlerGenerator, closeButtonForSnackbarAction
} from "../helpers/common"
import { useNotification } from "../../../utils/hooks/useNotification";

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


const KubernetesDataPanel = ({
  clusterInformation, classes, updateK8SConfig
}) => {

  const resetKubernetesConfig = () => updateK8SConfig({ k8sConfig : {
    inClusterConfig : false,
    k8sfile : "",
    contextName : "",
    clusterConfigured : false,
  }, })
  const { notify } = useNotification()

  const handleKubernetesDelete = () => {
    // showProgress()

    const handlerCb = () => resetKubernetesConfig()

    deleteKubernetesConfig(
      successHandlerGenerator(notify, "Kubernetes config removed", handlerCb),
      errorHandlerGenerator(notify, "Not able to remove config")
    )
  }

  const handleKubernetesClick = () => {
    // showProgress()
    pingKubernetes(
      successHandlerGenerator(notify, "Kubernetes pinged"/*, () => hideProgress()*/),
      errorHandlerGenerator(notify, "Kubernetes not pinged"/*, () => hideProgress()*/)
    )

  }

  return (
    <Grid container className={classes.infoContainer} xs={10}>

      <Grid item xs={12}>
        <AdapterChip
          label="Kubernetes"
          handleDelete={handleKubernetesDelete}
          handleClick={handleKubernetesClick}
          isActive={true}
          image="/static/img/kubernetes.svg"
        />
      </Grid>

      <Grid item xs={12} style={{ marginBottom : "1rem" }}>
        <Typography className={classes.infoTitle}>Details</Typography>
      </Grid>

      <Grid item xs={12} container>
        <Grid item xs={12}>
          <Typography className={classes.infoLabel}>Current-Context:  </Typography>
          <Typography className={classes.infoData}>{clusterInformation.contextName}</Typography>
          <Typography className={classes.infoData}>
            {clusterInformation.inClusterConfig
              ? "Using In Cluster Config"
              : "Using out of cluster config"}
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


const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch),
  updateK8SConfig : bindActionCreators(updateK8SConfig, dispatch), });

export default withStyles(styles)(connect(null, mapDispatchToProps)(KubernetesDataPanel))
