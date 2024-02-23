/* eslint-disable no-unused-vars */
import KubernetesIcon from '../icons/KubernetesIcon.js';
import KubernetesConfig from '../ConfigComponents/Kubernetes.js';
import ServiceCard from '../ServiceCard.js';
import { CircularProgress, Grid } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress, updateK8SConfig } from '../../../lib/store';
import { useEffect, useState } from 'react';
import { isKubernetesConnected, pingKubernetes } from '../helpers/kubernetesHelpers';
import KubernetesDataPanel from '../DataPanels/Kubernetes';

const KubernetesScreen = ({ k8sconfig, updateK8SConfig, updateProgress, setStepStatus }) => {
  const [clusterInformation, setClusterInformation] = useState({
    isClusterConfigured: k8sconfig.clusterConfigured,
    inClusterConfig: k8sconfig.inClusterConfig,
    configuredServer: k8sconfig.configuredServer,
    kubernetesPingStatus: false,
    contextName: k8sconfig.contextName,
    serverVersion: k8sconfig.server_version,
    k8sfile: k8sconfig.k8sfile,
  });

  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setStepStatus((prev) => ({ ...prev, kubernetes: isConnected }));
  }, [isConnected]);

  useEffect(() => {
    setIsLoading(true);
    pingKubernetes(
      (res) => {
        setIsLoading(false);
        setClusterInformation({
          ...clusterInformation,
          isClusterConfigured: k8sconfig.clusterConfigured,
          inClusterConfig: k8sconfig.inClusterConfig,
          kubernetesPingStatus: res !== undefined,
          contextName: k8sconfig.contextName,
          serverVersion: k8sconfig.server_version,
          configuredServer: k8sconfig.configuredServer,
        });
      },
      (err) => setIsLoading(false),
    );
  }, [k8sconfig.clusterConfigured, k8sconfig.inClusterConfig, k8sconfig.contextName]);

  useEffect(() => {
    setIsConnected(
      isKubernetesConnected(
        clusterInformation.isClusterConfigured,
        clusterInformation.kubernetesPingStatus,
      ),
    );
  }, [clusterInformation]);

  const kubeserviceInfo = {
    name: 'Kubernetes',
    logoComponent: KubernetesIcon,
    configComp: (
      <KubernetesConfig updateProgress={updateProgress} updateK8SConfig={updateK8SConfig} />
    ),
    clusterInformation,
  };

  const showDataPanel = () =>
    isKubernetesConnected(
      clusterInformation.isClusterConfigured,
      clusterInformation.kubernetesPingStatus,
    );

  return (
    <Grid item xs={12} container justify="center" alignItems="flex-start">
      <Grid
        item
        lg={6}
        sm={12}
        md={12}
        container
        justify="center"
        alignItems="flex-start"
        style={{ paddingLeft: '1rem' }}
      >
        <ServiceCard serviceInfo={kubeserviceInfo} isConnected={isConnected} />
      </Grid>
      <Grid
        item
        lg={6}
        sm={12}
        md={12}
        container
        justify="center"
        alignItems="center"
        style={{ paddingRight: '1rem' }}
      >
        {isLoading ? (
          <CircularProgress />
        ) : (
          showDataPanel() && (
            <KubernetesDataPanel
              clusterInformation={kubeserviceInfo.clusterInformation}
              setIsConnected={setIsConnected}
              setClusterInformation={setClusterInformation}
            />
          )
        )}
      </Grid>
    </Grid>
  );
};

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  return { k8sconfig };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(KubernetesScreen);
