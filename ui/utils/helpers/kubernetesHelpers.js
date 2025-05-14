import { EVENT_TYPES } from '../../lib/event-types';
import { CONNECTION_KINDS, CONNECTION_STATES } from '../Enum';
import {
  useLazyPingKubernetesQuery,
  useGetKubernetesContextsMutation,
  useConfigureKubernetesMutation,
} from '@/rtk-query/system';
import { useUpdateConnectionMutation } from '@/rtk-query/connection';

export const usePingKubernetes = () => {
  const [triggerKubernetes, { isLoading, data, error }] = useLazyPingKubernetesQuery();

  return {
    pingKubernetes: (connectionId, successCallback, errorCallback) => {
      triggerKubernetes(connectionId)
        .unwrap()
        .then((result) => {
          if (successCallback) successCallback(result);
        })
        .catch((err) => {
          if (errorCallback) errorCallback(err);
        });
    },
    isLoading,
    data,
    error,
  };
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

export const useDeleteKubernetesConfig = () => {
  const [updateConnection, { isLoading }] = useUpdateConnectionMutation();

  return {
    deleteKubernetesConfig: (successCallback, errorCallback, connectionId) => {
      updateConnection({
        connectionKind: CONNECTION_KINDS.KUBERNETES,
        connectionPayload: {
          [connectionId]: CONNECTION_STATES.DELETED,
        },
      })
        .unwrap()
        .then((result) => {
          if (successCallback) successCallback(result);
        })
        .catch((err) => {
          if (errorCallback) errorCallback(err);
        });
    },
    isLoading,
  };
};

export const useFetchContexts = () => {
  const [getContexts, { isLoading }] = useGetKubernetesContextsMutation();

  return {
    fetchContexts: (k8sfile, updateProgress) => {
      const formData = new FormData();
      formData.append('k8sfile', k8sfile);

      if (updateProgress) updateProgress({ showProgress: true });

      return getContexts(formData)
        .unwrap()
        .then((result) => {
          if (updateProgress) updateProgress({ showProgress: false });

          if (typeof result !== 'undefined') {
            let ctName = '';
            result.forEach(({ contextName, currentContext }) => {
              if (currentContext) {
                ctName = contextName;
              }
            });

            return { result, currentContextName: ctName };
          }
        })
        .catch((err) => {
          if (updateProgress) updateProgress({ showProgress: false });
          throw err;
        });
    },
    isLoading,
  };
};

export const useSubmitConfig = () => {
  const [configureKubernetes, { isLoading }] = useConfigureKubernetesMutation();

  return {
    submitConfig: (notify, updateProgress, updateK8SConfig, contextName, k8sfile) => {
      const inClusterConfigForm = false;
      const formData = new FormData();
      formData.append('inClusterConfig', inClusterConfigForm ? 'on' : ''); // to simulate form behaviour of a checkbox

      if (!inClusterConfigForm) {
        formData.append('contextName', contextName);
        formData.append('k8sfile', k8sfile);
      }

      if (updateProgress) updateProgress({ showProgress: true });

      return configureKubernetes(formData)
        .unwrap()
        .then((result) => {
          if (updateProgress) updateProgress({ showProgress: false });

          if (typeof result !== 'undefined') {
            if (notify)
              notify({
                message: 'Kubernetes config was validated!',
                event_type: EVENT_TYPES.SUCCESS,
              });

            if (updateK8SConfig) {
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

            return result;
          }
        })
        .catch((err) => {
          if (updateProgress) updateProgress({ showProgress: false });
          alert(err);
        });
    },
    isLoading,
  };
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
