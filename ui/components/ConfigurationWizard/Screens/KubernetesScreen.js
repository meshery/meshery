import KubernetesIcon from "../icons/KubernetesIcon.js"
import KubernetesConfig from "../SwitchConfigComponents/Kubernetes.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import { Grid } from "@material-ui/core"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"
import { updateProgress } from "../../../lib/store";
import { withSnackbar } from "notistack";
import { useEffect, useState } from "react"
import {pingKubernetes}  from "../helpers/kubernetesHelpers"
import KubernetesDataPanel from "../DataPanels/Kubernetes"


const KubernetesScreen = ({enqueueSnackbar, k8sconfig, updateProgress, closeSnackbar}) => { 

    const [clusterInformation, setClusterInformation] = useState({
      isClusterConfigured: false,
      inClusterConfig: false,
      kubernetesPingStatus: false
  })
  

useEffect(() => {

    if(k8sconfig.clusterConfigured) setClusterInformation({...clusterInformation, isClusterConfigured: true}) 
    else setClusterInformation({...clusterInformation, isClusterConfigured: false}) 

    if(k8sconfig.inClusterConfig) setClusterInformation({...clusterInformation, inClusterConfig: true}) 
    else setClusterInformation({...clusterInformation, inClusterConfig: false}) 

}, [k8sconfig])

const kubernetesPingUpdate = () => pingKubernetes(
    () => setClusterInformation({...clusterInformation, kubernetesPingStatus: true}),
    () => setClusterInformation({...clusterInformation, kubernetesPingStatus: false})
  )

useEffect(kubernetesPingUpdate, [])


const kubeserviceInfo = {
    name: "Kubernetes",
    logoComponent: KubernetesIcon,
    configComp :  KubernetesConfig,
    clusterInformation 
  }


  return (
    <Grid item xs={12} container justify="center" alignItems="flex-start"> 
        <Grid item xs={6} container justify="center">
          <ServiceSwitch serviceInfo={kubeserviceInfo} /> 
        </Grid>
        <Grid item xs={6} container justify="center">
          <KubernetesDataPanel clusterInformation={kubeserviceInfo.clusterInformation}/>
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

