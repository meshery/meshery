import KubernetesIcon from "../icons/KubernetesIcon.js"
import KubernetesConfig from "../SwitchConfigComponents/Kubernetes.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import { Grid } from "@material-ui/core"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"
import { updateProgress } from "../../../lib/store";
import { withSnackbar } from "notistack";
import { useEffect, useState } from "react"
import {isKubernetesConnected, pingKubernetes}  from "../helpers/kubernetesHelpers"
import KubernetesDataPanel from "../DataPanels/Kubernetes"


const KubernetesScreen = ({k8sconfig}) => { 

  const [clusterInformation, setClusterInformation] = useState({
    isClusterConfigured: k8sconfig.clusterConfigured,
    inClusterConfig: k8sconfig.inClusterConfig,
    kubernetesPingStatus: false,
    contextName: k8sconfig.contextName,
    serverVersion: k8sconfig.server_version
  })



  useEffect(() => {

    pingKubernetes(
      (res) => {
        setClusterInformation({
          isClusterConfigured: k8sconfig.clusterConfigured,
          inClusterConfig: k8sconfig.inClusterConfig,
          kubernetesPingStatus: (res !== undefined ? true : false),
          contextName: k8sconfig.contextName,
          serverVersion: k8sconfig.server_version
        })
      },
      (err) => console.log(err)
    )
    
  },[k8sconfig.clusterConfigured, k8sconfig.inClusterConfig, k8sconfig.contextName])




  const kubeserviceInfo = {
    name: "Kubernetes", logoComponent: KubernetesIcon,
    configComp :  KubernetesConfig,
    clusterInformation 
  }

  const showDataPanel = () => !isKubernetesConnected(clusterInformation.isClusterConfigured, clusterInformation.kubernetesPingStatus)

  return (
    <Grid item xs={12} container justify="center" alignItems="flex-start"> 
      <Grid item xs={6} container justify="center">
        <ServiceSwitch serviceInfo={kubeserviceInfo} /> 
      </Grid>
      <Grid item xs={6} container justify="center">
        {
          showDataPanel() &&
        <KubernetesDataPanel clusterInformation={kubeserviceInfo.clusterInformation}/>
        }
      </Grid>
    </Grid>
  )
}


const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig').toJS();
  return {
    k8sconfig,
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});


export default connect(mapStateToProps, mapDispatchToProps)(withSnackbar(KubernetesScreen))

