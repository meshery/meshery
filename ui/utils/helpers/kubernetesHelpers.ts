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

/**
 * Extracts kubernetes credentials from the data received from the server
 * @param {object} data - k8's config data received from the server
 * @returns {object} - k8's credentials
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
