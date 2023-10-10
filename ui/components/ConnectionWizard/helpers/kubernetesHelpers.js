import dataFetch from '../../../lib/data-fetch';
import { EVENT_TYPES } from '../../../lib/event-types';

/**
 * Pings kuberenetes server endpoint
 * @param  {(res) => void} successHandler
 * @param  {(err) => void} errorHandler
 */
export const pingKubernetes = (successHandler, errorHandler, connectionId) => {
  dataFetch(
    '/api/system/kubernetes/ping?connection_id=' + connectionId,
    { credentials: 'include' },
    successHandler,
    errorHandler,
  );
};

/**
 * Figures out if kubernetes connection is established or not
 *
 * @param {true|false} isClusterConfigured - data received from meshery server
 * as to whether or not the server config is found
 * @param {true|false} kubernetesPingStatus - found after pinging the kubernetes
 * server endpoint
 *
 * @return {true|false}
 */
export const isKubernetesConnected = (isClusterConfigured, kubernetesPingStatus) => {
  if (isClusterConfigured) {
    if (kubernetesPingStatus) return true;
  }

  return false;
};

export const deleteKubernetesConfig = (successCb, errorCb, connectionId) =>
  dataFetch(
    '/api/system/kubernetes/contexts/' + connectionId,
    {
      method: 'DELETE',
      credentials: 'include',
    },
    successCb,
    errorCb,
  );

export const fetchContexts = (updateProgress, k8sfile) => {
  const formData = new FormData();
  // if (inClusterConfigForm) {
  //   return;
  // }

  // formData.append('contextName', contextName);
  formData.append('k8sfile', k8sfile);

  updateProgress({ showProgress: true });

  return new Promise((res, rej) => {
    dataFetch(
      '/api/system/kubernetes/contexts',
      {
        method: 'POST',
        credentials: 'include',
        body: formData,
      },
      (result) => {
        updateProgress({ showProgress: false });

        if (typeof result !== 'undefined') {
          let ctName = '';
          result.forEach(({ contextName, currentContext }) => {
            if (currentContext) {
              ctName = contextName;
            }
          });

          res({ result, currentContextName: ctName });
        }
      },
      (err) => rej(err),
    );
  });
};

export const submitConfig = (notify, updateProgress, updateK8SConfig, contextName, k8sfile) => {
  const inClusterConfigForm = false;
  const formData = new FormData();
  formData.append('inClusterConfig', inClusterConfigForm ? 'on' : ''); // to simulate form behaviour of a checkbox
  if (!inClusterConfigForm) {
    formData.append('contextName', contextName);
    formData.append('k8sfile', k8sfile);
  }
  updateProgress({ showProgress: true });
  dataFetch(
    '/api/system/kubernetes',
    {
      method: 'POST',
      credentials: 'include',
      body: formData,
    },
    (result) => {
      updateProgress({ showProgress: false });
      if (typeof result !== 'undefined') {
        notify({ message: 'Kubernetes config was validated!', event_type: EVENT_TYPES.SUCCESS });
        updateK8SConfig({
          k8sConfig: {
            inClusterConfig: inClusterConfigForm,
            k8sfile,
            contextName: result.contextName,
            clusterConfigured: true,
            configuredServer: result.configuredServer,
          },
        });
      }
    },
    (err) => alert(err),
  );
};

/**
 * Extracts kubernetes credentials from the data received from the server
 * @param {object} data - k8's config data received from the server
 * @returns {object} - k8's credentials
 *
 */
export function extractKubernetesCredentials(data) {
  const credentials = {
    credentialName: data.name,
    secret: {
      clusterName: data.cluster.name,
      clusterServerURL: data.cluster.cluster.server,
      auth: {
        clusterUserName: data.auth.name,
        clusterToken: data.auth.user.token,
        clusterClientCertificateData: data.auth.user['client-certificate-data'],
        clusterCertificateAuthorityData: data.cluster.cluster['certificate-authority-data'],
        clusterClientKeyData: data.auth.user['client-key-data'],
      },
    },
  };

  return credentials;
}
